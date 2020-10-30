package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strconv"
	"syscall"
	"time"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type Packet struct {
	// Ignoring source and destination for now
	//SourcePort []byte
	//DestPort []byte
	SequenceNumber int
	// using naive scheme where checksum is equal to the sequence number
	Checksum int
	//headers map[string][]byte
	//data []byte
}

const MAXSEQNUM = 30

func main() {
	var seqNum uint32 = 1

	socket, err  := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_UDP)
	check(err)
	port := 9998
	addr := syscall.SockaddrInet4{Port: port, Addr: [4]byte{}}
	err = syscall.Bind(socket, &addr)
	check(err)

	for {
		err = sendPacket(socket, seqNum)
		check(err)

		timer1 := time.NewTimer(400 * time.Millisecond)
		giveUp := make(chan bool)
		go startTimeoutCoRoutine(seqNum, timer1, socket, 0, giveUp)

		respBuffer := make([]byte, 4096)
		numRec, respAddr, err := syscall.Recvfrom(socket, respBuffer, 0)
		check(err)
		seqNumACK := binary.LittleEndian.Uint32(respBuffer[:4])
		checksum := binary.LittleEndian.Uint32(respBuffer[4:numRec])
		fmt.Printf("received ACK for seqNum: %d\n, checksum: %d\n, from respAddr: %+v\n", seqNumACK, checksum, respAddr)
		if seqNumACK == checksum && seqNumACK == seqNum {
			if !timer1.Stop() {
				// This should always return true, but just in case, drain the channel
				fmt.Printf("drain channel\n")
				go func() {
					<-timer1.C
				}()
			}
			seqNum++
			fmt.Printf("incrementing seqNum!\n")
		}

		if seqNum >= MAXSEQNUM {
			break
		}

	}
	fmt.Printf("done!\n")
}

func sendPacket (socket int, seqNum uint32) error {
	// for now, we'll just send an int; later we'll build on this to send checksums
	p := make([]byte, 8)
	// add sequence number
	binary.LittleEndian.PutUint32(p, seqNum)
	// add checksum
	binary.LittleEndian.PutUint32(p[4:], seqNum)

	destPort, err := strconv.Atoi(os.Args[1])
	check(err)
	destAddr := syscall.SockaddrInet4{Port: destPort, Addr: [4]byte{127, 0, 0, 1}}
	err = syscall.Sendto(socket, p, 0, &destAddr)
	return err
}

func startTimeoutCoRoutine(seqNum uint32, timer *time.Timer, socket int, retryCount int, giveUp chan bool) {
	if retryCount >= 10 && seqNum >= MAXSEQNUM {
		fmt.Printf("tried retrying 10 times without luck, ending\n")
		// todo: still need to implement the use of this channel to end the program
		giveUp <- true
		return
	}
	<-timer.C
	// time expired, so resend
	err := sendPacket(socket, seqNum)
	check(err)
	_ = timer.Reset(100 * time.Millisecond)
	fmt.Println("Time expired, resending and restarting timer")
	retryCount++
	startTimeoutCoRoutine(seqNum, timer, socket, retryCount, giveUp)
}

// todo: implement checksum