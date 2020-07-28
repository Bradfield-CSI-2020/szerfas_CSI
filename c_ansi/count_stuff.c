#include <stdio.h>

main()
{
    int c;
    int count;

    c = getchar();
    count = 0;
    while( (c = getchar()) != EOF) {
        printf("c is %d\n", c);
        if (c == '\t')
            ++count;
    }
    printf("tab count is %d\n", count);
    printf("done\n");
    printf("Press enter to continue...\n");
    getchar();
}

