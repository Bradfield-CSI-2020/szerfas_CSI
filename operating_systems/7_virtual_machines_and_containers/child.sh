# CLONE_NEWUSER
[ `whoami` = "nobody" ] && echo "GOOD. new user" || echo "BAD. user preserved"

# CLONE_NEWNET
ip a | grep -q 2: && echo "BAD, can see network" || echo "GOOD, network isolated"

# CLONE_NEWPID
[ $$ -eq '1' ] && echo "GOOD, new PID namespace $$" || echo "BAD, sharing PID space $$"

# CLONE_NEWCGROUP
cat /proc/$$/cgroup


# CLONE_NEWUTS (Hostname and NIS domain name)
# I don't think the below test for hostname gets at what we want, or does it?
[ `hostname` != ubuntu-sz ] && echo "GOOD, new hostname `hostname`" || echo "BAD, same hostname `hostname`"

# CLONE_NEWNS (mount points)
# how should I test whether I'm on the same file system? Maybe create a file and see if I can see it from outside the process?
# this doesn't work though -- even with sudo I don't seem to have permission to create a file; perhaps I should run *inside* the child func in starter.c but before exec (so before the script)?
# touch test.txt
# sleep 5

# CLONE_NEWIPC
# First thought: try and access something like a named pipe that was created by the parent, but should no longer be visible by the child

# from manpages
#       Namespace   Constant          Isolates
#       Cgroup      CLONE_NEWCGROUP   Cgroup root directory
#       IPC         CLONE_NEWIPC      System V IPC, POSIX message queues
#       Network     CLONE_NEWNET      Network devices, stacks, ports, etc.
#       Mount       CLONE_NEWNS       Mount points
#       PID         CLONE_NEWPID      Process IDs
#       User        CLONE_NEWUSER     User and group IDs
#       UTS         CLONE_NEWUTS      Hostname and NIS domain name