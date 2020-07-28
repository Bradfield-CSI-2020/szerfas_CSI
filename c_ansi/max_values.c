#include <stdio.h>

//main()
//{
//    int c = 0;
//
//    while (c < 10000000000000) {
//        c = c + 1000;
//        printf("%d\n", c);
//    }
//}


#include <stdio.h>
#include <stdlib.h>
#include <limits.h>
#include <float.h>

int main(int argc, char** argv) {

//    double f = 0b1111111111111111111111111111111111111111111111111111111111111111;
//    float f = 0xFFF;
    double f = 1.7976931348623157e+309;
//    long double f = 0x7FFFFFFFFFFFFFFF;
    printf("float value is %f\n", f);
    printf("trying to get to infinity %f\n", f * f);

//    unsigned char c = 1;
//    signed char c2 = -1;
//
//    printf("adding one and negative one: %d\n", c > c2);
    printf("comparing -1L and 1U: %d\n", -1 < 1U);  /* seems this is false b/c -1 is signed and then promoted to unsigned, making it seem like a larger number */

//    int i;  /* = 260; */
//    char c = 125; /*  = 125 */
//
//    i = c;
//    c = i;
//
//    i = 260;
//    printf("expecting i to be 260 modulo 256: %d\n", i);


    printf("CHAR_BIT    :   %d\n", CHAR_BIT);
    printf("CHAR_MAX    :   %d\n", CHAR_MAX);
    printf("CHAR_MIN    :   %d\n", CHAR_MIN);
    printf("INT_BIT     :   %d\n", sizeof(int) * CHAR_BIT);
    printf("INT_MAX     :   %d\n", INT_MAX);
    printf("INT_MIN     :   %d\n", INT_MIN);
    printf("LONG_MAX    :   %ld\n", (long) LONG_MAX);
    printf("LONG_MIN    :   %ld\n", (long) LONG_MIN);
    printf("SCHAR_MAX   :   %d\n", SCHAR_MAX);
    printf("SCHAR_MIN   :   %d\n", SCHAR_MIN);
    printf("SHRT_MAX    :   %d\n", SHRT_MAX);
    printf("SHRT_MIN    :   %d\n", SHRT_MIN);
    printf("UCHAR_MAX   :   %d\n", UCHAR_MAX);
    printf("UINT_MAX    :   %u\n", (unsigned int) UINT_MAX);
    printf("ULONG_MAX   :   %lu\n", (unsigned long) ULONG_MAX);
    printf("USHRT_MAX   :   %d\n", (unsigned short) USHRT_MAX);

    return 0;
}


//#include <stdio.h>
//#include <stddef.h>
//#include <stdint.h>
//
//#define typename(x) _Generic((x),        /* Get the name of a type */             \
//                                                                                  \
//        _Bool: "_Bool",                  unsigned char: "unsigned char",          \
//         char: "char",                     signed char: "signed char",            \
//    short int: "short int",         unsigned short int: "unsigned short int",     \
//          int: "int",                     unsigned int: "unsigned int",           \
//     long int: "long int",           unsigned long int: "unsigned long int",      \
//long long int: "long long int", unsigned long long int: "unsigned long long int", \
//        float: "float",                         double: "double",                 \
//  long double: "long double",                   char *: "pointer to char",        \
//       void *: "pointer to void",                int *: "pointer to int",         \
//      default: "other")
//
//#define fmt "%20s is '%s'\n"
//int main() {
//
//  size_t s; ptrdiff_t p; intmax_t i; int ai[3] = {0}; return printf( fmt fmt fmt fmt fmt fmt fmt fmt,
//
//     "size_t", typename(s),               "ptrdiff_t", typename(p),
//   "intmax_t", typename(i),      "character constant", typename('0'),
// "0x7FFFFFFF", typename(0x7FFFFFFF),     "0xFFFFFFFF", typename(0xFFFFFFFF),
//"0x7FFFFFFFU", typename(0x7FFFFFFFU),  "array of int", typename(ai));
//}