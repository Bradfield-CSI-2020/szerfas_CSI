package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
)

var (
	PCAP_GLOBAL_HEADER_BYTE_LENGTH int = 24
	PCAP_PACKET_HEADER_BYTE_LENGTH int = 16
	PACKET_LENGTH_DATA_OFFSET int = 8  // number of bytes from start of packet header to length of captured packet data
	PACKET_NON_TRUNCATED_DATA_OFFSET int = 8  // number of bytes from start of packet header to length of packet data if it was not truncated
	ETHERTYPE_OFFSET int = 12  // number of bytes from start of ethernet frame to ethertype (which may indicate IPv4 or IPv6)
	ETHERFRAME_DESTINATION_OFFSET int = 0
	ETHERFRAME_SOURCE_OFFSET int = 6
	IP_HEADER_TOTAL_LENGTH_OFFSET int = 2
	IP_HEADER_DESTINATION_OFFSET int = 16
	IP_HEADER_SOURCE_OFFSET int = 12
	IP_HEADER_PROTOCOL_OFFSET int = 9  // offset from start of IP header to transport protocol; 6 means TCP
	TCP_HEADER_SOURCE_PORT_OFFSET int = 0
	TCP_HEADER_DESTINATION_PORT_OFFSET int = 2
	TCP_HEADER_SEQUENCE_NUMBER_OFFSET int = 4
	TCP_HEADER_OFFSET_OF_DATA_OFFSET int = 12
	HTTP_SOURCE_PORT_80 int = 80
)


func main() {
	// get file
	dat, err := ioutil.ReadFile("./net.cap")
	check(err)
	// pass into CountPackets
	packetCount := CountPackets(dat)
	fmt.Printf("packet count: %d\n", packetCount)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type TCPSegment struct {
	seqNumber int
	content []byte
}

func CountPackets(tcpdump []byte) int {
	//fmt.Printf("length of tcpdump: %d\n", len(tcpdump))
	count := 0
	// jump past the global header
	index := PCAP_GLOBAL_HEADER_BYTE_LENGTH
	// assert the packet completely captured
	sourceTCPPortsSet := make(map[string]bool)
	destinationTCPPortsSet := make(map[string]bool)
	//responseParts := make([]TCPSegment, 0)
	responseParts := make(map[int][]byte)
	for index < len(tcpdump) {
		//fmt.Printf("index: %d\n", index)
		//fmt.Printf("cound: %d\n", count)
		packetLength := Get4ByteValue(index + PACKET_LENGTH_DATA_OFFSET, tcpdump)
		if packetLength != Get4ByteValue(index + PACKET_NON_TRUNCATED_DATA_OFFSET, tcpdump) {
			panic("packet not completely captured")
		} else {
			etherframe := tcpdump[index + PCAP_PACKET_HEADER_BYTE_LENGTH:index + PCAP_PACKET_HEADER_BYTE_LENGTH + packetLength]
			fmt.Printf("ethertype: %d\n", GetEtherTypeFromEthernetFrame(etherframe))
			// todo: ask Oz about organizationally unique ID's here - go into the stretch goal
			fmt.Printf("source MAC: %d\n", GetSourceMACFromEthernetFrame(etherframe))
			fmt.Printf("destination MAC: %d\n", GetDestinationMACFromEthernetFrame(etherframe))
			ipFrame := etherframe[ETHERTYPE_OFFSET + 2:]
			internetHeaderLength := ipFrame[0] % 16 * 4 // we want the least 4 significant binary digits; this should be 20 at a minimum
			fmt.Printf("length of IP header: %d\n", internetHeaderLength)
			fmt.Printf("length of IP packet: %d\n", GetIPTotalLength(ipFrame))
			fmt.Printf("source IP: %s\n", GetSourceIPFromIPFrame(ipFrame))
			//fmt.Printf("source IP: %d\n", GetSourceIPFromIPFrame(ipFrame))
			fmt.Printf("destination IP: %s\n", GetDestinationIPFromIPFrame(ipFrame))
			//fmt.Printf("destination IP: %d\n", GetDestinationIPFromIPFrame(ipFrame))
			fmt.Printf("transport protocol: %d\n", GetTransportProtocolFromIPFrame(ipFrame))  // expect this to be 6, indicating TCP
			tcpSegment := ipFrame[internetHeaderLength:]
			tcpSourcePort := GetTCPSourcePort(tcpSegment)
			fmt.Printf("tcp source port: %d\n", tcpSourcePort)
			sourceTCPPortsSet[strconv.Itoa(tcpSourcePort)] = true
			tcpDestinationPort := GetTCPDestinationPort(tcpSegment)
			fmt.Printf("tcp destination port: %d\n", tcpDestinationPort)
			destinationTCPPortsSet[strconv.Itoa(tcpDestinationPort)] = true
			tcpSequenceNumber := GetTCPSequenceNumber(tcpSegment)
			fmt.Printf("tcp sequence number: %d\n", tcpSequenceNumber)

			tcpHeaderLength := int((tcpSegment[TCP_HEADER_OFFSET_OF_DATA_OFFSET] >> 4) << 2) // this will give is 4 most significant digits of the data offset byte, which is equal to the length of the header; minimum of 20 and max of 60
			fmt.Printf("tcp header length: %d\n", tcpHeaderLength)
			if tcpSourcePort == HTTP_SOURCE_PORT_80 {
				//responseParts = append(responseParts, TCPSegment{tcpSequenceNumber, tcpSegment[tcpHeaderLength:]})
				responseParts[tcpSequenceNumber] = tcpSegment[tcpHeaderLength:]
				//fmt.Printf("html: %s\n", string(tcpSegment[tcpHeaderLength:]))
			}
			count++
			index += packetLength + PCAP_PACKET_HEADER_BYTE_LENGTH
		}
	fmt.Printf("tcp source ports: %+v\n", sourceTCPPortsSet)
	fmt.Printf("tcp destination ports: %+v\n", destinationTCPPortsSet)
	}
	// sort packet order
	//sort.Slice(responseParts, func(i, j int) bool { return responseParts[i].seqNumber < responseParts[j].seqNumber })

	sequenceNumbers := make([]int, 0, len(responseParts))
	httpBytes := make([][]byte, len(responseParts))
	for k, _ := range responseParts {
		sequenceNumbers = append(sequenceNumbers, k)
	}
	sort.Slice(sequenceNumbers, func(i , j int) bool {return sequenceNumbers[i] < sequenceNumbers[j]})
	for i, v := range sequenceNumbers {
		httpBytes[i] = responseParts[v]
	}

	fmt.Printf("ok\n")
	response := bytes.Join(httpBytes, []byte{})
	//sort.Slice(response, func(i, j int) bool { return response[i] < response[j] })
	//fmt.Println(string(response[:1000]))

	// Split into HTTP header and body
	parts := bytes.SplitN(response, []byte{'\r', '\n', '\r', '\n'}, 2)

	//fmt.Println(parts[0])
	fmt.Println(string(parts[0]))
	fmt.Println(parts[1])

	// Write output
	out, _ := os.Create("out.jpeg")
	out.Write(parts[1])
	out.Close()
	fmt.Println("OK, result writen to out.jpeg")

	// jump to next header
	return count
}

func Get4ByteValue(index int, bytes []byte) int {
	//fmt.Printf("length of bytes: %d\n", len(bytes))
	//fmt.Printf("index: %d\n", index)
	//fmt.Printf("%+v\n", bytes[index:index + 4])
	//fmt.Printf("%+v\n", binary.LittleEndian.Uint32(bytes[index:index + 4]))
	// todo: accommodate machines that don't require reverse byte ordering
	return int(binary.LittleEndian.Uint32(bytes[index:index + 4]))
	//return int(binary.BigEndian.Uint32(bytes[index:index + 4]))
}

// todo: ask Oz about differences in byte-ordering
// here's why I think when the host and 	interpreting machines use opposite Endians we want to use littleEndian to decipher, and when on the same we want to use bigEndian:
// when displaying a value to us in a texteditor, the machine has already taken care to display it in bigEndian format for us, so if it's a littleEndian machine, it has flipped it for us. If it's a bigEndian machine, it's showing exactly how the bits were intepreted.
// if the host and interpreting machines use different byte ordering, we want to flip whatever we are seeing - and we are seeing bigEndian format on display.

// here's why I think the littleEndian byte ordering was important for interpreting packets, but once we're clearly in the ethernet protocol, we know that least significant bytes are sent first, so we no longer need to worry about host/interpreting machines.
// Since the tcpdump is displayed in BigEndian visually, we can use BigEndian within Ethernet to get what we need.
func GetEtherTypeFromEthernetFrame(eframe []byte) int {
	// todo: accommodate values below 1500
	return int(binary.BigEndian.Uint16(eframe[ETHERTYPE_OFFSET:ETHERTYPE_OFFSET + 2]))
}

func GetSourceMACFromEthernetFrame(eframe []byte) int {
	return int(binary.BigEndian.Uint16(eframe[ETHERFRAME_SOURCE_OFFSET:ETHERTYPE_OFFSET]))
}

func GetDestinationMACFromEthernetFrame(eframe []byte) int {
	return int(binary.BigEndian.Uint16(eframe[ETHERFRAME_DESTINATION_OFFSET:ETHERFRAME_SOURCE_OFFSET]))
}

func GetIPTotalLength(ipframe []byte) int {
	return int(binary.BigEndian.Uint16(ipframe[IP_HEADER_TOTAL_LENGTH_OFFSET:IP_HEADER_TOTAL_LENGTH_OFFSET + 2]))
}

func GetTransportProtocolFromIPFrame(ipframe []byte) int {
	return int(ipframe[IP_HEADER_PROTOCOL_OFFSET])
}
//func GetSourceIPFromIPFrame(ipframe []byte) int {
//	return int(binary.BigEndian.Uint32(ipframe[IP_HEADER_SOURCE_OFFSET:IP_HEADER_DESTINATION_OFFSET]))
//}

func GetSourceIPFromIPFrame(ipframe []byte) string {
	sourceSlice := ipframe[IP_HEADER_SOURCE_OFFSET:IP_HEADER_DESTINATION_OFFSET]
	result := make([]string, 4)
	for i := 0; i < 4; i++ {
		result[i] = strconv.Itoa(int(sourceSlice[i]))
	}
	return strings.Join(result, ".")
}

//func GetDestinationIPFromIPFrame(ipframe []byte) int {
//	return int(binary.BigEndian.Uint32(ipframe[IP_HEADER_DESTINATION_OFFSET:IP_HEADER_DESTINATION_OFFSET + 4]))
//}

func GetDestinationIPFromIPFrame(ipframe []byte) string {
	destinationSlice := ipframe[IP_HEADER_DESTINATION_OFFSET:IP_HEADER_DESTINATION_OFFSET + 4]
	result := make([]string, 4)
	for i := 0; i < 4; i++ {
		result[i] = strconv.Itoa(int(destinationSlice[i]))
	}
	return strings.Join(result, ".")
}

func GetTCPSourcePort(tcpframe []byte) int {
	return int(binary.BigEndian.Uint16(tcpframe[TCP_HEADER_SOURCE_PORT_OFFSET:TCP_HEADER_DESTINATION_PORT_OFFSET]))
}

func GetTCPDestinationPort(tcpframe []byte) int {
	return int(binary.BigEndian.Uint16(tcpframe[TCP_HEADER_DESTINATION_PORT_OFFSET:TCP_HEADER_DESTINATION_PORT_OFFSET + 2]))
}

func GetTCPSequenceNumber(bytes []byte) int {
	return int(binary.BigEndian.Uint32(bytes[TCP_HEADER_SEQUENCE_NUMBER_OFFSET:TCP_HEADER_SEQUENCE_NUMBER_OFFSET + 4]))
}

//func ConvertEndian(bytes []byte) []byte {
//	result := make([]byte, len(bytes))
//	for i, j := 0, len(bytes)-1; i < j; i, j = i+1, j-1 {
//		result[i], result[j] = bytes[j], bytes[i]
//	}
//	return result
//}
