package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"syscall"
)

// ping google.com
	// create raw socket
	// bind to any local port
	// generate packet
	// put packet into datagram?
	// set TTL of datagram -- can leave this until after can send and receive response from google.com
	// loop through sum number sending each datagram, noting response, and incrementing TTL

// todo:
//Parse responses filter by from address, sequence number
//Keep track of time to record RTT between my machine and every node
//Make hardcoded variables more dynamic - # of pings, target address, sequenceNum, etc.


func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Checksum(data []byte) uint16 {
	// Ones' complement of ones' complement sum of pairs of octets
	total := 0
	for i := 0; i < len(data); i += 2 {
		total += int(binary.BigEndian.Uint16(data[i : i+2]))
	}
	wrapped := uint16((total & 0xffff) + (total >> 16))
	return wrapped ^ 0xffff
}

func main () {
	socket, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_ICMP)
	check(err)
	err = syscall.Bind(socket, &syscall.SockaddrInet4{})

	for i := 0; i < 32; i++ {
		err = syscall.SetsockoptInt(socket, syscall.IPPROTO_IP, syscall.IP_TTL, i)
		check(err)


		// type value of 8 and code value of 0 indicates an Echo (ping) request for IPv4
		var pType uint8 = 8
		var pCode uint8 = 0
		// need any unique number - just using 1 for now
		var identifier uint16 = 1
		var ih, il uint8 = uint8(identifier>>8), uint8(identifier&0xff)
		// sequenceNum used to match responses to requests - starting with one and can increment as send requests
		var sequenceNum uint16 = 1
		var sh, sl uint8 = uint8(sequenceNum>>8), uint8(sequenceNum&0xff)
		// build icmpPacket omitting checksum, setting octets in Big Endian order
		icmpPacket := []byte{pType, pCode, 0, 0, il, ih, sl, sh}
		// calculate checksum and insert into the packet
		checksum := Checksum(icmpPacket)
		icmpPacket[2], icmpPacket[3] = uint8(checksum>>8), uint8(checksum&0xff)

		// send to IP and port for google.com; hardcoding IP address for now
		destPort := 0
		destAddr := [4]byte{172, 217, 5, 110}
		destSockAddr := syscall.SockaddrInet4{Port: destPort, Addr: destAddr}
		err = syscall.Sendto(socket, icmpPacket, 0, &destSockAddr)
		check(err)
		//fmt.Printf("sent: %+v\nto destination:%+v\n", icmpPacket, destSockAddr)
		//fmt.Printf("send error: %s\n", err)

		// recvfrom and print the response
		resp_buffer := make([]byte, 1500)
		//bytesReceived, from, err := syscall.Recvfrom(socket, resp_buffer, 0)
		_, from, err := syscall.Recvfrom(socket, resp_buffer, 0)
		check(err)
		fmt.Printf("received from: %+v\n", from)
		//fmt.Printf("received: %d\n%+v\n", bytesReceived, resp_buffer)
		//fmt.Printf("error: %s\n", err)
	}

}