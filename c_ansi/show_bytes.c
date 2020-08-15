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

//void show_float(float f) {
//    show_bytes((byte_pointer) &f, sizeof(float));
//}
//
//void show_pointer(void *p) {
//    show_bytes((byte_pointer) p, sizeof(void *));
//}

int main() {
    show_int(12345);
    show_int_casting_pointer_as_unsigned_char(12345);
    printf("sizeof int: %lu\n", sizeof(int));
    printf("sizeof int*: %lu\n", sizeof(int*));
    printf("sizeof char*: %lu\n", sizeof(char*));
    printf("sizeof void*: %lu\n", sizeof(void*));
//    show_float(12345);
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