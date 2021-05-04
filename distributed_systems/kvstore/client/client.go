package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

const (
	tcpLBAddress = "localhost:8080"
)

func isError(e error) {
	if e != nil {
		panic(e)
	}
}

//type MsgReq struct {
//	ReqType string
//	Key string
//	Val string
//}

func main() {
	// start connection to server
	conn, err := net.Dial("tcp", tcpLBAddress)
	isError(err)
	defer conn.Close()

	//var msgReq MsgReq

	// run loop looking for input
	inputBuf := bufio.NewReader(os.Stdin)
	respBuf := bufio.NewReader(conn)
	for {
		line, err := inputBuf.ReadString('\n')
		//fmt.Printf("input is: %s", line)
		isError(err)
		_, err = conn.Write([]byte(line))
		isError(err)
		resp, err := respBuf.ReadString('\n')
		isError(err)
		fmt.Println(resp)

		//args := strings.Fields(line)
		//
		//if len(args) != 2 || (args[0] != "set" && args[0] != "get") {
		//	os.Stderr.WriteString("Invalid command. Must be of form 'set foo=bar' or 'get foo'. Please try again\n")
		//}
		//
		//if len(args) == 2 {
		//	if args[0] == "set" {
		//		kv := strings.Split(args[1], "=")
		//		key := kv[0]
		//		val := kv[1]
		//		err = Set(key, val)
		//		msgReq.ReqType = "set"
		//		msgReq.Key = key
		//		msgReq.Val = val
		//		_, err := conn.Write([]byte(fmt.Sprintf("%v", msgReq)))
		//		if err != nil {
		//			fmt.Printf("Failed to set key '%s' with value '%s'. Error: %s\n", key, val, err)
		//		} else {
		//			fmt.Printf("Successfully set key '%s' with value '%s'\n", key, val)
		//		}
		//	} else if args[0] == "get" {
		//		key := args[1]
		//		value, err := Get(key)
		//		if err != nil {
		//			fmt.Printf("Failed to get key '%s'. Error: %s\n", key, err)
		//		} else {
		//			fmt.Printf("%s\n", value)
		//		}
		//	}
		//}


		//if strings.HasPrefix(line, "get") {
		//	//resp, err := Get(line, conn)
		//	_, err := conn.Write([]byte(line))
		//	isError(err)
		//	//fmt.Printf("%s\n", value)
		//} else if strings.HasPrefix(line, "set") {
		//	//resp, err := Set(line, conn)
		//	_, err := conn.Write([]byte(line))
		//	isError(err)
		//	fmt.Printf("Successfully set key '%s' with value '%s'\n", key, val)
		//} else {
		//	fmt.Printf("command not found: %s\n", line)
		//}
	}
}

//func Get(line string) {
//
//}
//
//func Set(line string) string {
//
//}
