#include <stdio.h>

main()
{
    char c;
    int i;
    short int si;
    long int li;
    float f;
    double d;
    long double ld;

    long double range[10];
    range[0] = 255;
    range[1] = 65535;
    range[2] = 4294967296;
    range[3] = 255;
    range[4] = 255;
    range[5] = 255;
    range[6] = 255;

    c = 267;
    si = 257;
    i = 65537;

    printf("hello world\n");
    printf("%d\n", 255);
    printf("%c\n", c);
    printf("%hi\n", si);
    printf("%i\n", i);
}