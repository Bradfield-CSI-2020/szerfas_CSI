#include <stdio.h>

/* print Fahrenheit-Celsius table for fahr = 0 ,20, ..., 300 */
main()
{
    float celsius, fahr;
    int upper, lower, step;

    step = 10;
    upper = 150;
    lower = 0;

    printf("Celsius\tFahrenheit\n");

    celsius = lower;
    while (celsius <= upper) {
        fahr = 9.0/5.0 * celsius + 32;
        printf("%6.0f\t\t%6.0f\n", celsius, fahr);
        celsius = celsius + step;
    }
}