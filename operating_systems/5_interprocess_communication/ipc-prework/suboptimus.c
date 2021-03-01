#include <signal.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <sys/wait.h>
#include <unistd.h>

int START = 2, END = 20;
char *TESTS[] = {"brute_force", "brutish", "miller_rabin"};
int num_tests = sizeof(TESTS) / sizeof(char *);


// for every test: create a process and a set of pipes to handle that test
// with 3 tests, this does not use all the cores on the machine (8 in my case, often 4)


// alternative design
// create a process for every core on the machine; allow each process to handle any test passed in
// just use two pipes -- all processes read from the two pipes and execute the appropriate job

// Question: why use pipes rather than named pipes or sockets?

int main(int argc, char *argv[]) {
  int testfds[num_tests][2];
  int resultfds[num_tests][2];
  int result, i;
  long n;
  pid_t pid;

  for (i = 0; i < num_tests; i++) {
    pipe(testfds[i]);
    pipe(resultfds[i]);

    pid = fork();

    if (pid == -1) {
      fprintf(stderr, "Failed to fork\n");
      exit(-1);
    }

    if (pid == 0) {
      // we are the child, connect the pipes correctly and exec!
      close(testfds[i][1]);  // test input pipe: close the write because we'll only be reading
      close(resultfds[i][0]);  // result input pipe: close the read because we'll only be writing
      dup2(testfds[i][0], STDIN_FILENO);  // replace STDIN with the test input reader
      dup2(resultfds[i][1], STDOUT_FILENO);  // replace STDOUT with the result input writer
      execl("primality", "primality", TESTS[i], (char *)NULL);
    }

    // we are the parent
    close(testfds[i][0]);  // close the reader b/c we're only writing to test input pipe
    close(resultfds[i][1]);  // close the writer b/c we're only reading from the results input pipe
  }

  // for each number, run each test
  for (n = START; n <= END; n++) {
    for (i = 0; i < num_tests; i++) {

      // we are the parent, so send test case to child and read results
      write(testfds[i][1], &n, sizeof(n));
      read(resultfds[i][0], &result, sizeof(result));
      printf("%15s says %ld %s prime\n", TESTS[i], n, result ? "is" : "IS NOT");
    }
  }
}
