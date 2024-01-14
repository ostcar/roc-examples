#include <stdlib.h>

struct RocStr {
    char* bytes;
    size_t len;
    size_t capacity;
};

extern void roc__mainForHost_1_exposed_generic(void* *ptr, const struct RocStr *arg);
extern void roc__mainForHost_0_caller(const unsigned char *flags, void* *closure_data, unsigned char *output);
