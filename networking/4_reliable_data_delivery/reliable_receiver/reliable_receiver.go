package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strconv"
	"syscall"
)

const MAXSEQNUM = 30 - 1

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	var highestSeqNum uint32 = 0

	socket, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_UDP)
	check(err)
	inboundPort := 9999
	inboundAddr := syscall.SockaddrInet4{Addr: [4]byte{127, 0, 0, 1}, Port: inboundPort}
	err = syscall.Bind(socket, &inboundAddr)
	check(err)
	fmt.Printf("reliable sender listening on port: %d\n", inboundPort)

	for {
		buffer := make([]byte, 4096)
		numRec, fromAddr, err := syscall.Recvfrom(socket, buffer, 0)
		seqNumRec := binary.LittleEndian.Uint32(buffer[:4])
		checksum := binary.LittleEndian.Uint32(buffer[4:numRec])
		fmt.Printf("recieved seqNumber: %d\nchecksum: %d\nfrom: %+v\n", seqNumRec, checksum, fromAddr)
		if seqNumRec == checksum && seqNumRec == (highestSeqNum + 1) {
			highestSeqNum = seqNumRec
		}

		// respond using different unreliable server to mimic real world
		//outboundPort, err := strconv.Atoi(os.Args[1])
		//check(err)
		//outboundAddr := syscall.SockaddrInet4{Addr: [4]byte{127, 0, 0, 1}, Port: outboundPort}
		//resp := make([]byte, 4)
		//binary.LittleEndian.PutUint32(resp, seqNumRec)

		//err = syscall.Sendto(socket, resp, 0, &outboundAddr)
		//check(err)
		err = sendPacket(socket, highestSeqNum)
		check(err)
		fmt.Printf("Ack sent for seqNum: %d\n", highestSeqNum)
		if highestSeqNum >= MAXSEQNUM {
			break
		}
	}
	fmt.Printf("done!")

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