#include <math.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

int brute_force(long n);
int brutish(long n);
int miller_rabin(long n);

struct job {
  int test; // which test function?
  long n;   // number to test for primality
};

struct job_result {
  int test;
  long n;
  int result;
};

void exit_with_usage() {
  fprintf(stderr, "Usage: ./primality [brute_force|brutish|miller_rabin]\n");
  exit(1);
}

int main(int argc, char*argv[]) {
  long num;
  struct job j;
  struct job_result result;
  int (*func)(long);

  int (*funcs[])(long) = {brute_force, brutish, miller_rabin};

  if (isatty(STDIN_FILENO)) { // interactive, so expect test name and provide
                              // prompt

    if (argc != 2)
      exit_with_usage();

    if (strcmp(argv[1], "brute_force") == 0)
      func = &brute_force;
    else if (strcmp(argv[1], "brutish") == 0)
      func = &brutish;
    else if (strcmp(argv[1], "miller_rabin") == 0)
      func = &miller_rabin;
    else
      exit_with_usage();

    fprintf(stderr, "Running \"%s\", enter a number:\n> ", argv[1]);
    while (scanf("%ld", &num) == 1) {
      printf("%d\n", (*func)(num));
      fflush(stdout);
      fprintf(stderr, "> ");
    }
    }
  else { // non-interactive, so obtain test function and input from pipe
    for (;;) {
      read(STDIN_FILENO, &j, sizeof(j));
      result.test = j.test;
      result.n = j.n;
      result.result = (funcs[j.test])(j.n);
      write(STDOUT_FILENO, &result, sizeof(result));
    }
  }
}

/*
 * Primality test implementations
 */

// Just test every factor
int brute_force(long n) {
  for (long i = 2; i < n; i++)
    if (n % i == 0)
      return 0;
  return 1;
}

// Test factors, up to sqrt(n)
int brutish(long n) {
  long max = floor(sqrt(n));
  for (long i = 2; i <= max; i++)
    if (n % i == 0)
      return 0;
  return 1;
}

int randint(int a, int b) { return rand() % (++b - a) + a; }

int modpow(int a, int d, int m) {
  int c = a;
  for (int i = 1; i < d; i++)
    c = (c * a) % m;
  return c % m;
}

int witness(int a, int s, int d, int n) {
  int x = modpow(a, d, n);
  if (x == 1)
    return 1;
  for (int i = 0; i < s - 1; i++) {
    if (x == n - 1)
      return 1;
    x = modpow(x, 2, n);
  }
  return (x == n - 1);
}

// TODO we should probably make this a parameter!
int MILLER_RABIN_ITERATIONS = 10;

// An implementation of the probabilistic Miller-Rabin test
int miller_rabin(long n) {
  int a, s = 0, d = n - 1;

  if (n == 2)
    return 1;

  if (!(n & 1) || n <= 1)
    return 0;

  while (!(d & 1)) {
    d >>= 1;
    s += 1;
  }
  for (int i = 0; i < MILLER_RABIN_ITERATIONS; i++) {
    a = randint(2, n - 1);
    if (!witness(a, s, d, n))
      return 0;
  }
  return 1;
}
