package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"syscall"
	"time"
)

// create socket
// bind socket
// listen for incoming connections
// accept incoming connections
// receive from accepted connection

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type CacheEntry struct {
	lastModified time.Time
	object string
}

func main () {
	cache := make(map[string]CacheEntry)

	socket, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
	check(err)

	destAddr := [4]byte{127,0,0,1}
	destPort := 7777
	destSocket := syscall.SockaddrInet4{Port: destPort, Addr: destAddr}

	err = syscall.Bind(socket, &destSocket)
	check(err)

	for {
		err = syscall.Listen(socket, 3)
		fmt.Printf("listening to socket: %d\n", socket)
		if err != nil {
			fmt.Printf("received err in listening: %s", err)
			break
		}
		localSocketFileDescriptor, destSocket, err := syscall.Accept(socket)
		check(err)
		// call go routine to execute listening and sending on the new socket
		//go handleAndEchoConnection(localSocketFileDescriptor, destSocket)
		cacheValue, cacheFlag := handleAndForwardConnection(localSocketFileDescriptor, destSocket, cache)
		if cacheFlag {
			requestLineComponents := strings.Fields(cacheValue.requestLine)
			lmTime, err := time.Parse(time.RFC1123, cacheValue.headerLines["Date"])
			check(err)
			entry := CacheEntry{
				lastModified: lmTime,
				object: cacheValue.body,
			}
			cache[requestLineComponents[1]] = entry
		}

		fmt.Printf("connection accepted\n")
		fmt.Printf("starting new routine at socket: %d\n with destAddr: %+v\n", socket, destSocket)
	}

	fmt.Printf("exiting main\n")
}


func handleAndEchoConnection(localSocketFD int, destSocket syscall.Sockaddr) {
	for {
		// how much
		resp_buffer := make([]byte, 1500)
		bytesReceived, msgSource, err := syscall.Recvfrom(localSocketFD, resp_buffer, 0)
		check(err)
		fmt.Printf("bytes received: %d\n", bytesReceived)
		fmt.Printf("msg source: %+v\n", msgSource)
		packet := resp_buffer[:bytesReceived]
		err = syscall.Sendto(localSocketFD, packet, 0, destSocket)
	}
}

func handleAndForwardConnection(clientSocketFD int, clientSocket syscall.Sockaddr, cache map[string]CacheEntry) (cachEntry HTTPPacket, cacheFlag bool) {
	defer syscall.Close(clientSocketFD)
	// create a new socket
	serverSocketFD, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
	check(err)

	// bind to any available address
	err = syscall.Bind(serverSocketFD, &syscall.SockaddrInet4{Port: 0, Addr: [4]byte{0, 0, 0, 0}})
	check(err)

	serverAddr := [4]byte{127,0,0,1}
	serverPort := 9000
	//serverPort := 7778
	serverSocket := syscall.SockaddrInet4{Port: serverPort, Addr: serverAddr}

	err = syscall.Connect(serverSocketFD, &serverSocket)
	defer syscall.Close(serverSocketFD)
	fmt.Printf("connected to server\n")

	resp_buffer := make([]byte, 1500)
	//for {
	//	bytesReceived, _, err := syscall.Recvfrom(clientSocketFD, resp_buffer, 0)
	//	check(err)
	//	//fmt.Printf("error is: %s\n", err)
	//	packet := resp_buffer[:bytesReceived]
	//	fmt.Printf("# bytes received from client: %d\n", bytesReceived)
	//	fmt.Printf("packet: %s\n", packet)
	//	err = syscall.Sendto(serverSocketFD, packet, 0, &serverSocket)
	//	check(err)
	//	bytesReceived, _, err = syscall.Recvfrom(serverSocketFD, resp_buffer, 0)
	//	check(err)
	//	packet = resp_buffer[:bytesReceived]
	//	fmt.Printf("# bytes received from server: %d\n", bytesReceived)
	//	fmt.Printf("packet: %s\n", packet)
	//	err = syscall.Sendto(clientSocketFD, packet, 0, clientSocket)
	//	check(err)
	//}
	cacheFlag = false
	bytesReceived, _, err := syscall.Recvfrom(clientSocketFD, resp_buffer, 0)
	check(err)
	//fmt.Printf("error is: %s\n", err)
	packet := resp_buffer[:bytesReceived]
	fmt.Printf("# bytes received from client: %d\n", bytesReceived)
	fmt.Printf("packet: %s\n", packet)
	httpPacket :=parseHTTPPacket(string(packet))
	if host, ok := httpPacket.headerLines["Host"]; ok {
		// hardcoding cache for localhost
		if strings.HasPrefix(host, "localhost") {
			cacheFlag = true
			if entry, ok := cache[host]; ok {
				// todo: reach out to server and see if last modified
				// todo: set headers on response
				// send object
				err := syscall.Sendto(clientSocketFD, []byte(entry.object), 0, clientSocket)
				check(err)
				fmt.Printf("sent data from cache: %s\n", entry.object)
				return HTTPPacket{}, !cacheFlag
			}
		}
	}

	err = syscall.Sendto(serverSocketFD, packet, 0, &serverSocket)
	check(err)
	bytesReceived, _, err = syscall.Recvfrom(serverSocketFD, resp_buffer, 0)
	check(err)
	packet = resp_buffer[:bytesReceived]
	fmt.Printf("# bytes received from server: %d\n", bytesReceived)
	fmt.Printf("packet: %s\n", packet)
	httpPacket =parseHTTPPacket(string(packet))
	err = syscall.Sendto(clientSocketFD, packet, 0, clientSocket)
	check(err)
	if cacheFlag {
		return httpPacket, true
	} else {
		return HTTPPacket{}, cacheFlag
	}
}

type HTTPPacket struct {
	requestLine string
	headerLines map[string]string
	body string
}

func parseHTTPPacket(p string) HTTPPacket {

	endOfRequestLine, _ := regexp.Compile("\r\n")
	endOfRequestLineLoc := endOfRequestLine.FindStringIndex(p)
	fmt.Println(endOfRequestLine.FindStringIndex(p))
	endOfHeaderLines, _ := regexp.Compile("\r\n\r\n")
	endOfHeaderLinesLoc := endOfHeaderLines.FindStringIndex(p)
	headers := p[endOfRequestLineLoc[1]:endOfHeaderLinesLoc[0]]
	headersMap := make(map[string]string)
	for _, v := range strings.Split(headers, "\r\n") {
		split, _ := regexp.Compile(":")
		splitLoc := split.FindStringIndex(v)
		key := v[:splitLoc[0]]
		val := v[splitLoc[1]:]
		headersMap[key] = val
	}

	fmt.Println(endOfHeaderLines.FindStringIndex(p))
	return HTTPPacket{
		requestLine: p[:endOfRequestLineLoc[0]],
		headerLines: headersMap,
		body: p[endOfHeaderLinesLoc[1]:],
	}
}