#include <signal.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <sys/wait.h>
#include <unistd.h>

int START = 2, END = 10000;
char *TESTS[] = {"brute_force", "brutish", "miller_rabin"};
int num_tests = sizeof(TESTS) / sizeof(char *);

struct jobUnit {  // to be passed in on input pipe
    int method;  // index of set of test methods to use for the job
    long n;  // value under test
    int resultfd;  // output pipe
};

struct resultUnit {  // to be passed into input of result pipe
    int method;  // index of set of test methods used for job
    long n;  // value under test
    bool isPrime;  // result of test
};

// alternative design
// create a process for every core on the machine; allow each process to handle any test passed in
// just use two pipes -- all processes read from the two pipes and execute the appropriate job

// Question: why use pipes rather than named pipes or sockets?

int main(int argc, char *argv[]) {
  int jobfds[2];
  int resultfds[2];
  int i, x;
  struct jobUnit job;
  struct resultUnit result;
  long n;
  pid_t pid;

  int numCPUs = sysconf(_SC_NPROCESSORS_ONLN);

  pipe(jobfds);
  pipe(resultfds);

  // create workers on same number of threads as we have on this machine
  for (i = 0; i < numCPUs; i++) {
    pid = fork();
    if (pid == -1) {
      fprintf(stderr, "Failed to fork\n");
      exit(-1);
    }
    if (pid == 0) {
      // we are the child, connect the pipes correctly and exec!
      // the child process inherited open file descriptors for both ends of the pipe, but we only need one end for each pipe
      close(jobfds[1]);  // close the test input pipe write end because we'll only be reading
      close(resultfds[0]);  // close the result input pipe read end because we'll only be writing
      dup2(jobfds[0], STDIN_FILENO);  // replace STDIN with the test input reader
      dup2(resultfds[1], STDOUT_FILENO);  // replace STDOUT with the result input writer
      execl("primality", "primality", TESTS[i], (char *)NULL);
    }
  }
  // from here on out we're only in the parent because the child processes are now executing primality code

  // parent does not use the test input reader and result input write
  close(jobfds[0]);
  close(resultfds[1]);

  // loop through start and end numbers we want to test
  // create a job for each and send into test input pipe write end
  // then begin pulling from result entry
  // problem with this approach? We're creating a long queue before we start pulling answers
  // it may be better to send just enough jobs to activate each worker, then process and submit the next job with each loop iteration through the end of the test
  // calling this premature optimization, though, and just queueing up all the jobs right away

  for (n = START; n < END; n++) {
    for (i = 0; i < num_tests; i ++) {
      job.method = i % num_tests;
      job.n = n;
      write(resultfds[1], &job, sizeof(job));
    }
  }

  for (n = START; n < END; n++) {
    read(resultfds[0], &result, sizeof(result));
    printf("%15s says %ld %s prime\n", TESTS[result.method], result.n, result.isPrime ? "is" : "IS NOT");
  }
}
