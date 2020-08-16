#include <stdio.h>

typedef unsigned char *byte_pointer;

void show_bytes(int *start, int len) {
    int i;
    for (i = 0; i < len; i ++) {
        printf(" %.2x", start[i]);
    }
    printf("\n");
}

void show_int(int i) {
    show_bytes(&i, sizeof(int));
//    show_bytes((byte_pointer) &i, sizeof(int));
}


void show_bytes_unsigned_char_pointer_start(unsigned char *start, int len) {
    int i;
    for (i = 0; i < len; i ++) {
        printf(" %.2x", start[i]);
    }
    printf("\n");
}


void show_int_casting_pointer_as_unsigned_char(int i) {
    show_bytes_unsigned_char_pointer_start((unsigned char*) &i, sizeof(int));
}

void show_float(float f) {
//    show_bytes((byte_pointer) &f, sizeof(float));
    show_bytes_unsigned_char_pointer_start((unsigned char*) &f, sizeof(float));
}

//void show_pointer(void *p) {
//    show_bytes((byte_pointer) p, sizeof(void *));
//}

// note: this machine uses little endian format
int main() {
    show_int(12345);
    show_int_casting_pointer_as_unsigned_char(12345);
    printf("sizeof int: %lu\n", sizeof(int));
    printf("sizeof int*: %lu\n", sizeof(int*));
    printf("sizeof char*: %lu\n", sizeof(char*));
    printf("sizeof void*: %lu\n", sizeof(void*));
    show_float(12345);


//    float f = 0;
//    unsigned int *p;
//    p = (unsigned int *) &f;
//    (*p) = (*p) | 1U << 31;
//    (*p) = (*p) | 1U << 26;
//    (*p) = (*p) | 1U << 22;
//    (*p) = (*p) | 1U << 20;
//    (*p) = (*p) | 1U << 18;

    float f = 0;
    unsigned int v = *(unsigned int *)&f;
    v |= 1U << 31;
    v |= 1U << 26;
    v |= 1U << 22;
    v |= 1U << 20;
    v |= 1U << 18;

    // expecting this to print equal 42.5 but can't get the bitwise manipulation of floats to work
    printf("%f\t%lu\n", (float) v, sizeof(f));
//    void *p;
//    show_pointer(p);
}


//#include <stdio.h>
//
//typedef unsigned char *byte_pointer;
//
//void show_bytes(byte_pointer start, int len) {
//    int i;
//    for (i = 0; i < len; i++)
//        printf(" %.2x", start[i]);
//    printf("\n");
//}
//
//void show_int(int x) {
//    show_bytes((byte_pointer) &x, sizeof(int));
//}
//
//void show_float(float x) {
//    show_bytes((byte_pointer) &x, sizeof(float));
//}
//
//void show_pointer(void *x) {
//    show_bytes((byte_pointer) &x, sizeof(void *));
//}
//
//int main() {
//    show_int(12345);
//    show_float(12345);
//    void *p;
//    show_pointer(p);
//}