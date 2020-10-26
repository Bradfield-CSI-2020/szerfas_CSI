package main

import (
	"fmt"
	"log"
	"syscall"
)






type DNSHeader struct {
	ID uint16
 	Flags uint16
	QDCount uint16
	ANCount uint16
	NSCount uint16
	ARCount uint16
}
// 0  1  2  3  4  5  6  7  8  9  10 11 12 13 14 15
//+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//|                      ID                       |
//+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//|QR|   Opcode  |AA|TC|RD|RA|   Z    |   RCODE   |
//+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//|                    QDCOUNT                    |
//+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//|                    ANCOUNT                    |
//+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//|                    NSCOUNT                    |
//+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//|                    ARCOUNT                    |
//+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+


func main() {
	//p :=  make([]byte, 2048)
	//conn, err := net.Dial("udp", "127.0.0.1:1234")
	//if err != nil {
	//	fmt.Printf("Some error %v", err)
	//	return
	//}
	//fmt.Fprintf(conn, "Hi UDP Server, How are you doing?")
	//_, err = bufio.NewReader(conn).Read(p)
	//if err == nil {
	//	fmt.Printf("%s\n", p)
	//} else {
	//	fmt.Printf("Some error %v\n", err)
	//}
	//conn.Close()

	socket, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_UDP)
	if err != nil {
		log.Fatal(err)
	}

	destAddr := [4]byte{8,8,8,8}
	//destAddr := [4]byte{127,0,0,1}
	destPort := 53
	//destPort := 7777
	destSocket := syscall.SockaddrInet4{Port: destPort, Addr: destAddr}
	// can I get away without using bind since connect() and send() may automatically bind?
	//localPort := 0 // get udp port
	err = syscall.Bind(socket, &syscall.SockaddrInet4{Port: 0, Addr: [4]byte{0, 0, 0, 0}})
	if err != nil {
		log.Fatal(err)
	}

	packet := make([]byte, 0)

	// todo: set packet to valid outgoing dns request
	packet = append(packet, 0xd2, 0xf5, 0x01, 0x20, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x0b, 0x62, 0x72, 0x61, 0x64, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x63, 0x73, 0x03, 0x63, 0x6f, 0x6d, 0x00, 0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0x29, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00)
	//packet = append(packet, 0x68, 0x65, 0x6c, 0x6c, 0x6f)
	err = syscall.Sendto(socket, packet, 0, &destSocket)
	if err != nil {
		log.Fatal(err)
	}

	resp_buffer := make([]byte, 1500)
	resp_int, from, err := syscall.Recvfrom(socket, resp_buffer, 0)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("resp_int: %d\n", resp_int)
	fmt.Println(resp_buffer)
	fmt.Printf("from addr in recvfrom: %+v\n", from)
	//header := constructHeader()


}

func ConfirmDNSResponseID(resp []byte) {

}

func parseDNSResponse(resp []byte) {

}



// psudocode
// get socket
// bind socket (hopefully to any port, but can specify and defer cleanup too
// sendto(socket, destaddr, destport)
// listen and get response