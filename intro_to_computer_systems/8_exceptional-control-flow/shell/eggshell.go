package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

func MainLoop() {
	var line string
	var err error

	b := bufio.NewReader(os.Stdin)

	for {
		line, err = b.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("Have an egggxellent day!")
				break
			}
			os.Stderr.WriteString("error! please try another command")
		}
		if strings.Contains(line, "^Q") {
			return
		} else {
			fmt.Printf("output line is: %s", line)
		}
	}
}

// have main loop call multiplexer
// multiplexer does the following:
// declares a channel on which to receive signals, passes that channel to the syscall library
// declares a channel on which to receive user input, passes that channel to a separate go routine called ReceiveInput
// runs switch statement inside infinite loop
// first case is receiving on a channel from signals: prints and exits
// second case is receiving on a channel from input: prints and loops

func main() {
	args := os.Args[1:]
	fmt.Printf("args are: ")
	fmt.Println(args)
	//if len(args) > 0 {
	//	done := make(chan bool)
	//	inputReceived
	//	HandleArgs(done, args)
	//}
	MainLoop2()
	fmt.Println("🍳🍳🍳 Goodbye! 🍳🍳🍳")
	fmt.Println()
}

func MainLoop2() {

	//signals := make(chan os.Signal, 1)  // note: this creates a buffered channel, I wonder if it will work with an unbuffered channel?
	//signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	userInput := make(chan bool)
	done := make(chan bool)

	args := os.Args[1:]
	go ReceiveInput(userInput, done, args)

	for {
		select {
		case <-userInput:
			go ReceiveInput(userInput, done, nil)
		//case <-signals:
		//	fmt.Println()
		//	fmt.Println()
		//	fmt.Println("Have an 🥚zellent day!")
		//	return
		case <-done:
			fmt.Println()
			fmt.Println()
			fmt.Println("Have an 🥚zellent day!")
			return
		}
	}
}

func ReceiveInput(inputReceived chan bool, exit chan bool, args []string) {
	if args != nil && len(args) > 0 {
		HandleArgs(inputReceived, exit, args)
	} else {
		var line string
		var err error
		b := bufio.NewReader(os.Stdin)
		fmt.Printf("🥚> ")
		line, err = b.ReadString('\n')
		if err != nil {
			os.Stderr.WriteString("error! please try another command")
		} else {
			args := strings.Fields(line)
			HandleArgs(inputReceived, exit, args)
		}
	}
}

func HandleArgs(inputReceived chan bool, exit chan bool, args []string) {
	if args[0] == "exit"{
		exit <- true
		return
	}

	switch args[0] {
	case "ls":
		fmt.Println("command is ls")
		path, err := os.Getwd()
		if err != nil {
			log.Println(err)
		}
		files, err := ioutil.ReadDir(path)
		if err != nil {
			log.Println(err)
		}
		for _, f := range files {
			fmt.Println(f.Name())
		}
	case "ls_binary":
		cmd := exec.Command("ls", args[1:]...)   // we can get pid of child via cmd.Process.Pid
		var out bytes.Buffer
		cmd.Stdout = &out
		sigs := make(chan os.Signal, 1)
		// this effectively blocks SIGINT and SIGTERM until we receive an handle, making it safe to start the child without race conditions
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		err := cmd.Start()
		if err != nil {
			log.Fatal(err)
		}
		go func() {
			if sig := <-sigs; sig == syscall.SIGINT {
				err := syscall.Kill(cmd.Process.Pid, syscall.SIGINT)
				if err != nil {
					log.Fatal("unable to interrupt child process")
				}
				fmt.Println()
				fmt.Println("")
			}
			fmt.Println()
			fmt.Println("Ending go routine that would interrupt child process")
		}()
		err = cmd.Wait()
		if err != nil {
			log.Fatal(err)
		}
		// unblock and reset default shell behavior for SIGINT and SIGTERM as we no longer want to pass on to client
		signal.Reset()
		fmt.Println(cmd.Stdout)
	case "sleep":
		cmd := exec.Command("sleep", args[1:]...)   // we can get pid of child via cmd.Process.Pid
		var out bytes.Buffer
		cmd.Stdout = &out
		sigs := make(chan os.Signal, 1)
		// this effectively blocks SIGINT and SIGTERM until we receive a handle, making it safe to start the child without race conditions
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		err := cmd.Start()
		if err != nil {
			log.Fatal(err)
		}
		go func(sigs chan os.Signal) {
			sig := <-sigs
			if sig == syscall.SIGINT {
				err := syscall.Kill(cmd.Process.Pid, syscall.SIGINT)
				if err != nil {
					log.Fatal("unable to interrupt child process")
				}
				fmt.Println()
				fmt.Println("child process interrupted")
			}
			fmt.Println()
		}(sigs)
		err = cmd.Wait()
		if err != nil {
			fmt.Printf("child process exited with error: %s", err)
		}
		// unblock and reset default shell behavior for SIGINT and SIGTERM as we no longer want to pass on to client
		signal.Reset(syscall.SIGINT, syscall.SIGTERM)
		fmt.Println(cmd.Stdout)
	case "pwd":
		fmt.Println("command is pwd")
		path, err := os.Getwd()
		if err != nil {
			log.Println(err)
		}
		fmt.Println(path)
	case "echo":
		fmt.Println(strings.Join(args[1:], " "))
	default:
		fmt.Println("command not recognized, please try another command")
	}
	inputReceived <- true
}