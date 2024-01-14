package main

/*
#cgo LDFLAGS: ./False.o -ldl
#include "./host.h"
*/
import "C"

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"unsafe"
)

func main() {
	if len(os.Args[1:]) == 0 {
		fmt.Println("Please pass a .false file as a command-line argument to the false interpreter!")
		os.Exit(1)
	}

	arg := rocStrFromStr(os.Args[1])

	var callback unsafe.Pointer
	C.roc__mainForHost_1_exposed_generic(&callback, &arg)

	var output C.uchar
	var flags C.uchar
	C.roc__mainForHost_0_caller(&flags, &callback, &output)
	os.Exit(int(output))
}

const is64Bit = uint64(^uintptr(0)) == ^uint64(0)

func rocStrFromStr(str string) C.struct_RocStr {
	var rocStr C.struct_RocStr
	rocStr.len = C.ulong(len(str))
	rocStr.capacity = rocStr.len
	ptr := unsafe.StringData(str)
	rocStr.bytes = (*C.char)(unsafe.Pointer(ptr))
	return rocStr
}

func rocStrRead(rocStr C.struct_RocStr) string {
	if int(rocStr.capacity) < 0 {
		// Small string
		ptr := (*byte)(unsafe.Pointer(&rocStr))

		byteLen := 12
		if is64Bit {
			byteLen = 24
		}

		shortStr := unsafe.String(ptr, byteLen)
		len := shortStr[byteLen-1] ^ 128
		return shortStr[:len]
	}

	// Remove the bit for seamless string
	len := uint64(rocStr.len) & ^uint64(1<<63)
	ptr := (*byte)(unsafe.Pointer(rocStr.bytes))
	return unsafe.String(ptr, len)
}

//export roc_fx_openFile
func roc_fx_openFile(name *C.struct_RocStr) uintptr {
	file, err := os.Open(rocStrRead(*name))
	if err != nil {
		panic(fmt.Sprintf("can not open file: %w", err))
	}
	return uintptr(unsafe.Pointer(file))
}

//export roc_fx_closeFile
func roc_fx_closeFile(filePtr unsafe.Pointer) {
	file := (*os.File)(filePtr)
	file.Close()
}

//export roc_fx_getFileBytes
func roc_fx_getFileBytes(output *C.struct_RocStr, filePtr unsafe.Pointer) {
	file := (*os.File)(filePtr)
	buf := make([]byte, 0x10) // This is intentionally small to ensure correct implementation
	count, err := file.Read(buf)
	if err != nil && err != io.EOF {
		panic(fmt.Sprintf("can not read from file: %v", err))
	}
	str := rocStrFromStr(string(buf[:count]))
	*output = str
}

//export roc_fx_getChar
func roc_fx_getChar() C.char {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return C.char(text[0])
}

//export roc_fx_putLine
func roc_fx_putLine(line *C.struct_RocStr) {
	fmt.Println(rocStrRead(*line))
}

//export roc_fx_putRaw
func roc_fx_putRaw(line *C.struct_RocStr) {
	fmt.Print(rocStrRead(*line))
}

//export roc_alloc
func roc_alloc(size C.ulong, alignment int) unsafe.Pointer {
	return C.malloc(size)
}

//export roc_realloc
func roc_realloc(ptr unsafe.Pointer, newSize, _ C.ulong, alignment int) unsafe.Pointer {
	return C.realloc(ptr, newSize)
}

//export roc_dealloc
func roc_dealloc(ptr unsafe.Pointer, alignment int) {
	C.free(ptr)
}

//export roc_panic
func roc_panic(msg C.struct_RocStr) {
	panic(rocStrRead(msg))
}
