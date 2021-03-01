#include <signal.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <sys/socket.h>
#include <unistd.h>

#define START 2
#define END 10000
#define WORKERS 4

char *TESTS[] = {"brute_force", "brutish", "miller_rabin"};
int num_tests = sizeof(TESTS) / sizeof(char *);

struct job {
  int test; // which test function?
  long n;   // number to test for primality
};

struct job_result {
  int test;
  long n;
  int result;
};

int main(int argc, char *argv[]) {

  int jobfd[2], resultfd[2], i;
  struct job jb;
  struct job_result result;
  pid_t pid;

  // create the two pipes
  pipe(jobfd);
  pipe(resultfd);

  // create all of the worker processes
  for (i = 0; i < WORKERS; i++) {
    pid = fork();

    if (pid == -1) {
      fprintf(stderr, "Failed to fork\n");
      exit(-1);
    }

    if (pid == 0) { // we are the child
      dup2(jobfd[0], STDIN_FILENO);
      dup2(resultfd[1], STDOUT_FILENO);
      close(jobfd[1]);
      close(resultfd[1]);
      execl("primality", "primality", (char *)NULL);
    }

  }
  close(jobfd[0]);
  close(resultfd[1]);

  // start workers
  for (i = 0; i < WORKERS; i++) {
    jb.test = i % 3;
    jb.n = START + i / 3;
    write(jobfd[1], &jb, sizeof(jb));
  }

  for (int completed = 0; completed < num_tests * (END - START); completed++) {
    read(resultfd[0], &result, sizeof(result));
    printf("%20s says %ld %s prime\n", TESTS[result.test], result.n,
           result.result ? "is" : "IS NOT");
    if (i < num_tests * (END - START)) {
      jb.test = i % 3;
      jb.n = START + i / 3;
      write(jobfd[1], &jb, sizeof(jb));
      i++;
    }
  }
}
