package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func isError(e error) {
	if e != nil {
		panic(e)
	}
}

type replica struct {
	Conn net.Conn
	RespBuf *bufio.Reader
	RemoteAddr string
}

var replicas []*replica
var leader *replica
var roundRobinPointer int

func main() {
	args := os.Args[1:]
	// args are addresses of servers at which to connect
	// the first address is the one at which this load balancer server will listen
	ln, err := net.Listen("tcp", args[0])
	fmt.Printf("kvstore load balancer listening at localhost:%s\n", args[0])
	isError(err)
	// all other addresses are of kvstore replicas at which to send requests
	for _, addr := range args[1:] {
		upstreamConn, err := net.Dial("tcp", addr)
		isError(err)
		replicas = append(replicas, &replica{upstreamConn, bufio.NewReader(upstreamConn), addr})
		defer upstreamConn.Close()
	}
	fmt.Printf("replicas: %v\n", replicas)
	// by default, the first replica address is set to leader
	createNewLeader(0)

	for {
		downStreamConn, err := ln.Accept()
		isError(err)
		go handleConnection(downStreamConn)
	}
}

func createNewLeader(index int) {
	// all non-leader addresses are sent to the leader for replication
	// the first address will be targeted for synchronous replication
	// all others will be targeted async
	addressStr := ""
	for i, replica := range replicas {
		if i != index {
			addressStr += " " + replica.RemoteAddr
		}
	}
	fmt.Printf("addressStr is: %s\n", addressStr)
	_, err := replicas[index].Conn.Write([]byte(fmt.Sprintf("promote-to-leader %s\n", addressStr)))
	leader = replicas[index]
	isError(err)
	// this moves the offset to the end of the connection buffer so that future pass throughs do not read the promotion response
	resp, err := leader.RespBuf.ReadString('\n')
	isError(err)
	fmt.Printf(resp)
}

//func confirmLeaderOrFailover() {
//	closed := make(chan bool)
//	go detectClosed(leader.RespBuf, closed)
//	select {
//	case isClosed := <- closed:
//		if isClosed == true {
//			// remove leader from list of replicas
//			leaderAddr := leader.Conn.RemoteAddr()
//			newReplicas := make([]*replica, len(replicas))
//			for _, replica := range replicas {
//				if replica.Conn.RemoteAddr() != leaderAddr {
//					newReplicas = append(newReplicas, replica)
//				}
//			}
//			replicas = newReplicas
//			// unset leader pointer
//			leader = nil
//			// create new leader
//			createNewLeader(0)
//		}
//
//	default:
//		// leader is healthy, continue
//		return
//	}
//}

//func detectClosed(buf *bufio.Reader, closed chan bool) {
//	_, err := leader.RespBuf.Peek(1)
//	if  err == io.EOF {
//		fmt.Printf("%s detected closed LAN connection from leader")
//		leader.Conn.Close()
//		closed <- true
//	} else {
//		closed <- false
//	}
//}

//func confirmLeaderOrFailover() {
//	_, err := leader.RespBuf.Peek(1)
//	if  err == io.EOF {
//		fmt.Printf("%s detected closed LAN connection from leader")
//		// remove leader from list of replicas
//		leaderAddr := leader.Conn.RemoteAddr()
//		newReplicas := make([]*replica, len(replicas))
//		for _, replica := range replicas {
//			if replica.Conn.RemoteAddr() != leaderAddr {
//				newReplicas = append(newReplicas, replica)
//			}
//		}
//		replicas = newReplicas
//
//		// close connection and unset leader pointer
//		leader.Conn.Close()
//		leader = nil
//
//		// create new leader
//		createNewLeader(0)
//	}
//
//}

func handleConnection(downstreamConn net.Conn) {
	inputBuf := bufio.NewReader(downstreamConn)
	for {
		line, err := inputBuf.ReadString('\n')
		isError(err)
		fmt.Printf("input received from client: %s", line)

		args := strings.Fields(line)

		// check to see if leader connection is healthy.
		// If not, execute failover
		//confirmLeaderOrFailover()

		switch {
		case args[0] == "set":
			passThrough(downstreamConn, leader, line)
		case args[0] == "get":
			passThrough(downstreamConn, replicas[roundRobinPointer], line)
			roundRobinPointer++
			if roundRobinPointer > len(replicas) - 1 {roundRobinPointer = 0}
			// guard against any replicas that have been removed from the list
			for(replicas[roundRobinPointer] == nil) {roundRobinPointer++}
			if roundRobinPointer > len(replicas) - 1 {roundRobinPointer = 0}
		}
	}
}

func passThrough(downstreamConn net.Conn, upstreamReplica *replica, msg string) {
	_, err := upstreamReplica.Conn.Write([]byte(msg))
	isError(err)
	resp, err := upstreamReplica.RespBuf.ReadString('\n')
	fmt.Printf("received response from passing through upstream: %s", resp)
	isError(err)
	// pass response through to client
	_, err = downstreamConn.Write([]byte(resp))
	isError(err)
}