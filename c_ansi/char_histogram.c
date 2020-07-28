#include <stdio.h>

main()
{
    int c, i, nwhite, nother;
    int nchar[26];
    int ndigit[10];

    nwhite = nother = 0;

    for (i = 0; i < 26; ++i)
        ndigit[i] = 0;

    for (i = 0; i < 10; ++i)
        ndigit[i] = 0;

    while ((c = getchar()) != EOF) {
        if (c >= 'a' && c <= 'z')
            ++nchar[c - 'a'];
        else if (c == '\n' || c =='\t' || c ==' ')
            ++nwhite;
        else if (c >= '0' && c <= '9')
            ++ndigit[c - '0'];
        else
            ++nother;
    }

    char cc;
    for (cc = 0; cc < 26; cc++)
        printf("%c: %d\n", 'a' + cc, nchar[cc]);
}