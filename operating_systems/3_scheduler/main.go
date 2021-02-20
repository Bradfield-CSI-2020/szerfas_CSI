package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"
)

/*
A program whose only purpose is to log any signals it receives (excluding those where it is unable to log anything before being terminated).
Test this methodically with either kill(1) or kill(2)
*/
//func main() {
//	sig_chan := make(chan os.Signal, 1)
//	signal.Notify(sig_chan)
//	i := 0
//	for {
//		sig := <- sig_chan
//		fmt.Printf("Received signal: %s\n", sig)
//		if sig == syscall.SIGINT || sig == syscall.SIGTERM {
//			fmt.Printf("exiting\n")
//			break
//		}
//		i++
//		fmt.Printf("count is: %d\n", i)
//		if i == 7 {
//			fmt.Printf("exiting\n")
//			return
//		}
//	}
//}


/*
A program that runs a child process, and notes how the child terminates (ie normally or due to a signal, and in the later case which signal)
*/
func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	if err := exec.CommandContext(ctx, "sleep", "1").Run(); err != nil {
		// This will fail after 100 milliseconds. The 1 second sleep will be interrupted.
		// Printing here shows the signal (os.Process.Kill) used to terminate the program
		// this fulfills the exercise requirements, but is very high level.
		// More learning is to be had by implementing with use of a lower-level Go library, like syscall
		fmt.Println(err)
	}

	// ////////  BEGIN LOW-LEVEL IMPLEMENTATION  //////// //

	// fork and exec a child process
	// have that child process sleep for 3 seconds and return
	// in the meantime, offer the ueser the ability to send a kill signal to the child process
	// in all cases, record and print to stdoutput how the child process terminated
	procAttr := &syscall.ProcAttr{}
	fmt.Printf("spawning a child process that will sleep for 3s.\nWould you like to send a kill signal (y/n)?\n")
	childPid, err := syscall.ForkExec("/bin/sleep", []string{"5"}, procAttr)
	fmt.Printf("child_pid is: %d\n", childPid)
	if err != nil {
		fmt.Printf("Received error in forking: %s\n", err)
		return
	}

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		panic(fmt.Errorf("input not valid: %s\n", err))
	}
	if input != "y\n" && input != "n\n" {
		fmt.Printf("Invalid input received. Please type only 'y' or 'n'.\n")
		fmt.Printf("Would you like to send a kill signal (y/n)?\n")
	} else if input == "y\n" {
		fmt.Printf("sending kill signal\n")
		err := syscall.Kill(childPid, syscall.SIGKILL)
		if err != nil {
			fmt.Printf("Error in sending kill signal to child process: %s\n", err)
			return
		}
	} else if input == "n\n" {
		fmt.Printf("received input 'n', not sending kill signal\n")
	}
	fmt.Printf("input was: %s\n", input)

	var waitStatus syscall.WaitStatus
	wpid, err := syscall.Wait4(childPid, &waitStatus, 0, nil)
	if err != nil {
		fmt.Printf("received error in calling wait: %s\n", err)
	}
	fmt.Printf("received wpid: %d\n", wpid)
	fmt.Printf("child process exited: %t\n", waitStatus.Exited())
	fmt.Printf("child process signaled: %t\n", waitStatus.Signaled())
	fmt.Printf("signal received: %s\n", waitStatus.Signal())

}





// scratch
//func main() {
	//reader := bufio.NewReader(os.Stdin)
	//input, err := reader.ReadString('\n')
	//if err != nil {
	//	panic(fmt.Errorf("input not valid: %s\n", err))
	//}
	//
	//fmt.Printf("spawning new process with input: %s\n", input)
	//
	//inputCommand := exec.Command("bash", "-c", input)
	//output, err := inputCommand.Output()
	//if err != nil {
	//	panic(fmt.Errorf("command failed: %s\n", err))
	//}
	//fmt.Printf("output of spawned process:\n%s\n", output)
	//}