package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
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
	MainLoop2()
	fmt.Println("🍳 Goodbye! 🍳")
}

func MainLoop2() {
	signals := make(chan os.Signal, 1)  // note: this creates a buffered channel, I wonder if it will work with an unbuffered channel?
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	userInput := make(chan bool)
	go ReceiveInput(userInput)
	for {
		select {
		case <-userInput:
			go ReceiveInput(userInput)
		case <-signals:
			fmt.Println()
			fmt.Println("🥚🥚🥚 Have an egggzellent day! 🥚🥚🥚")
			return
		}
	}
}

func ReceiveInput(done chan bool) {
	var line string
	var err error
	b := bufio.NewReader(os.Stdin)
	line, err = b.ReadString('\n')
	if err != nil {
		os.Stderr.WriteString("error! please try another command")
	} else {
		fmt.Printf("output line is: %s", line)
	}
	done <-true
}