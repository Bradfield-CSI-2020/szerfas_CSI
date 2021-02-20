# include <stdio.h>
# include <stdlib.h>
# include <unistd.h>
# include <string.h>
# include <sys/wait.h>

char *FORK_FAILED = "fork failed!\n";

int main(int argc, char *argv[]) {
    printf("hello world (pid:%d)\n", (int) getpid());
    int rc = fork();
    if (rc < 0) { // fork failed
        fprintf(stderr, "%s", FORK_FAILED);
        exit(1);
    } else if (rc == 0) { // child (new process)
        printf("hello world, I'm a child (pid:%d)\n", (int) getpid());
        char *myargs[3];
        myargs[0] = strdup("wc");
        myargs[1] = strdup("process_api.c");
        myargs[2] = NULL;  // marks end of array
        execvp(myargs[0], myargs);
        printf("this will not print!");
    } else {
        int rc_wait = wait(NULL);
        printf("hello world, I'm the parent of %d (process: %d)\n", rc, (int) getpid());
    }
}