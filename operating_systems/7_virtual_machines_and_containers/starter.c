#define _GNU_SOURCE
#include <sched.h>
#include <sys/wait.h>
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>

#include <errno.h>
#include <fcntl.h>
#include <sched.h>
#include <sys/mount.h>
#include <sys/prctl.h>
#include <linux/prctl.h>
#include <linux/capability.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <sys/wait.h>
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>

#define STACK_SIZE 65536

struct child_config {
  int argc;
  char **argv;
};

/* Entry point for child after `clone` */
int child(void *arg) {
  int num;
  pid_t pid = getpid();
  printf("child pid: %d\n", pid);
  int fd;
//  fd = open("/sys/fs/cgroup/pids/dokker/cgroup.procs", O_RDWR | O_APPEND);
  fd = open("/sys/fs/cgroup/pids/dokker/cgroup.procs", O_RDWR | O_TRUNC);
  if (fd == -1) printf("Errno: %d\n", errno);

//  ssize_t bytes2 = write(1, "hello\n", 6);  // this works fine
//  ssize_t bytes = write(fd, string_pid, 5);
//  if (bytes == -1) printf("Errno: %d\n", errno);
//  printf("%li bytes written\n", bytes);
//  close(fd);
  int bytes = dprintf(fd, "%d\n", getpid());
  printf("%d bytes written\n", bytes);
  if (bytes == -1) printf("Errno: %d\n", errno);
  close(fd);

  // confirm wrote the PID I wanted too -- this was primarily for debugging and took FOREVER to figure out
  // Clearly I need to work through K&R
  fd = open("/sys/fs/cgroup/pids/dokker/cgroup.procs", O_RDWR   );
  if (fd == -1) printf("Errno: %d\n", errno);
  char arr[bytes];
  read(fd, &arr, bytes);
  printf("PID read from /sys/fs/cgroup/pids/dokker/cgroup.procs: %s\n", arr);
  close(fd);

//  // get program to stop because want to check something
//  char name[20];
//  printf("Hello. Pausing so you can inspect /sys/fs/cgroup/pids/dokker/ if you like... hit 'enter' to proceed.\n");
//  fgets(name,20,stdin);

  // limit the number of processes that will work in this cgroup
  fd = open("/sys/fs/cgroup/pids/dokker/pids.max", O_RDWR | O_TRUNC);
  if (fd == -1) printf("Errno: %d\n", errno);
  // limit to just one c group for testing
  bytes = dprintf(fd, "%d\n", 5);
  printf("%d bytes written\n", bytes);
  if (bytes == -1) printf("Errno: %d\n", errno);
  close(fd);


  struct child_config *config = arg;
  if (execvpe(config->argv[0], config->argv, NULL)) {
    fprintf(stderr, "execvpe failed %m.\n");
    return -1;
  }
  return 0;
}


int main(int argc, char**argv) {
  struct child_config config = {0};
//  int flags = CLONE_NEWCGROUP | CLONE_NEWIPC | CLONE_NEWNET | CLONE_NEWNS| CLONE_NEWPID | CLONE_NEWUSER | CLONE_NEWUTS;
  int flags = CLONE_NEWCGROUP | CLONE_NEWNET | CLONE_NEWUTS | CLONE_NEWUSER | CLONE_NEWPID;
  pid_t child_pid = 0;

  // Prepare child configuration
  config.argc = argc - 1;
  config.argv = &argv[1];

  // Allocate stack for child
  char *stack = 0;
  if (!(stack = malloc(STACK_SIZE))) {
    fprintf(stderr, "Malloc failed");
    exit(1);
  }

  // Clone parent, enter child code
  if ((child_pid = clone(child, stack + STACK_SIZE, flags | SIGCHLD, &config)) == -1) {
    fprintf(stderr, "Clone failed");
    exit(2);
  }

  wait(NULL);
  
  return 0;
}
