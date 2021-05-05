package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
)

func isError(e error) {
	if e != nil {
		panic(e)
	}
}

// global variables -- may want to pass these down through functions instead
var dataStore map[string]string
var storageFile *os.File
var leadermode bool
var connections = make([]net.Conn, 0)
// used for sync connections with follower when in leadermode
var respBuf *bufio.Reader

func main() {
	num, address := parseNumAddressFromArgs()
	ln, err := net.Listen("tcp", ":" + address)
	isError(err)
	fmt.Printf("kvstore server listening at localhost:%s\n", address)

	fileName := fmt.Sprintf("kv_store%d.json", num)
	storageFile, err = os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0755)
	isError(err)
	defer storageFile.Close()

	dataStore = make(map[string]string)

	for {
		conn, err := ln.Accept()
		isError(err)
		go handleConnection(conn)
	}
}

func parseNumAddressFromArgs() (num int, address string) {
	args := os.Args[1:]
	var err error
	if len(args) == 0 {
		num = 0
	} else {
		num, err = strconv.Atoi(args[0])
		isError(err)
	}
	address = strconv.Itoa(8080 + num)
	return num, address
}

func handleConnection(conn net.Conn) {
	fmt.Printf("received connection with localAddr: %s and remoteAddr: %s\n", conn.LocalAddr(), conn.RemoteAddr())
	inputBuf := bufio.NewReader(conn)


	for {
		line, err := inputBuf.ReadString('\n')
		isError(err)
		fmt.Printf("input received: %s", line)

		args := strings.Fields(line)

		switch {
		case args[0] == "set":
			if len(args) != 2 {
				msg := "Invalid command. Must be of form 'set foo=bar'. Please try again\n"
				os.Stderr.WriteString(msg)
				conn.Write([]byte(msg))
				continue
			}

			kv := strings.Split(args[1], "=")
			key := kv[0]
			val := kv[1]
			err := Set(key, val)
			if err != nil {
				msg := fmt.Sprintf("Failed to set key '%s' with value '%s'. Error: %s\n", key, val, err)
				fmt.Printf(msg)
				conn.Write([]byte(msg))
				continue
			}

			if leadermode {
				// this is statement-based replication which is not ideal, but fine for MVP
				msg := []byte(line)
				_, err = connections[0].Write(msg)
				isError(err)
				// block until receive response from replica
				resp, err := respBuf.ReadString('\n')
				isError(err)
				fmt.Printf("In leader mode and sent single synchronous request to %s\nReceived in response: %s\n", connections[0].RemoteAddr(), resp)
				// send non-blocking async requests to all other replicas
				for _, c := range connections[1:] {
					_, err = c.Write(msg)
					isError(err)
				}
			}
			msg := fmt.Sprintf("Successfully set key '%s' with value '%s'\n", key, val)
			fmt.Printf(msg)
			conn.Write([]byte(msg))
		case args[0] == "get":
			if len(args) != 2 {
				msg := "Invalid command. Must be of form 'get foo'. Please try again\n"
				os.Stderr.WriteString(msg)
				conn.Write([]byte(msg))
			} else {
				key := args[1]
				value, err := Get(key)
				if err != nil {
					msg := fmt.Sprintf("Failed to get key '%s'. Error: %s\n", key, err)
					fmt.Printf(msg)
					conn.Write([]byte(msg))
				} else {
					msg := fmt.Sprintf("%s\n", value)
					fmt.Printf(msg)
					conn.Write([]byte(msg))
				}
			}
		case args[0] == "promote-to-leader":
			leadermode = true
			// promote-to-leader is accompanied by a list of addresses
			// the first is the address of the replica at which to replicate synchronously
			// all others receive updates asynchronously
			for _, addr := range args[1:] {
				// todo: handle closures of connections more gracefully; deferring in a for-loop may be a resource leak
				conn, err := net.Dial("tcp", addr)
				isError(err)
				defer conn.Close()
				connections = append(connections, conn)
			}
			respBuf = bufio.NewReader(connections[0])

			msg := fmt.Sprintf("promoted to leadermode with connections: %v, respBuf: %p\n", connections, respBuf)
			fmt.Printf(msg)
			conn.Write([]byte(msg))
		case args[0] == "demote-from-leader":
			if !leadermode {
				fmt.Printf("received demotion but already not leader")
				continue
			}
			leadermode = false
			for _, c := range connections {
				c.Close()
			}
		}

		//if len(args) != 2 || (args[0] != "set" && args[0] != "get") {
		//	msg := "Invalid command. Must be of form 'set foo=bar' or 'get foo'. Please try again\n"
		//	os.Stderr.WriteString(msg)
		//	conn.Write([]byte(msg))
		//}
		//
		//if len(args) == 2 {
		//	if args[0] == "set" {
		//		kv := strings.Split(args[1], "=")
		//		key := kv[0]
		//		val := kv[1]
		//		err = Set(key, val)
		//		if err != nil {
		//			msg := fmt.Sprintf("Failed to set key '%s' with value '%s'. Error: %s\n", key, val, err)
		//			fmt.Printf(msg)
		//			conn.Write([]byte(msg))
		//		} else {
		//			msg := fmt.Sprintf("Successfully set key '%s' with value '%s'\n", key, val)
		//			fmt.Printf(msg)
		//			conn.Write([]byte(msg))
		//		}
		//		// todo: send to replica
		//	} else if args[0] == "get" {
		//		key := args[1]
		//		value, err := Get(key)
		//		if err != nil {
		//			msg := fmt.Sprintf("Failed to get key '%s'. Error: %s\n", key, err)
		//			fmt.Printf(msg)
		//			conn.Write([]byte(msg))
		//		} else {
		//			msg := fmt.Sprintf("%s\n", value)
		//			fmt.Printf(msg)
		//			conn.Write([]byte(msg))
		//		}
		//	}
		//}
	}
}

func Set(k string, v string) (err error) {
	// todo: handle locking the file so other connections can't create race conditions
	// todo: handle unique ID for each record (might involve another index so that we can handle deletes by unique ID -- scratch that, we can just use key as the ID and overwrite duplicates for this implementation)
	dataStore[k] = v
	return FlushToDisk()
}

func FlushToDisk() error {
	jsonString, err := json.Marshal(dataStore)
	isError(err)
	if err != nil {
		return err
	}
	_, err = storageFile.Seek(0, 0)
	isError(err)
	if err != nil {
		return err
	}
	_, err = storageFile.Write(jsonString)
	isError(err)
	if err != nil {
		return err
	}
	return nil
}

func Get(k string) (v string, err error) {
	err = refreshDataStore()
	if err != nil {
		return "", err
	}
	v, inMap := dataStore[k]
	if !inMap {
		return "", errors.New(fmt.Sprintf("Value for key '%s' not found. Error: %s\n", k, err))
	}
	return v, nil
}

func refreshDataStore() error {
	storageFile.Seek(0, 0)
	byteValue, _ := ioutil.ReadAll(storageFile)
	return json.Unmarshal(byteValue, &dataStore)
}
