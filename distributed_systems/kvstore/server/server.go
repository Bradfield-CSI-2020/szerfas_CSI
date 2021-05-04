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
	// the first argument is a server number. No number or 0 indicates a leader. Every other number indicates a follower
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
	inputBuf := bufio.NewReader(conn)

	for {
		line, err := inputBuf.ReadString('\n')
		//fmt.Printf("input received: %s", line)

		args := strings.Fields(line)

		if len(args) != 2 || (args[0] != "set" && args[0] != "get") {
			msg := "Invalid command. Must be of form 'set foo=bar' or 'get foo'. Please try again\n"
			os.Stderr.WriteString(msg)
			conn.Write([]byte(msg))
		}

		if len(args) == 2 {
			if args[0] == "set" {
				kv := strings.Split(args[1], "=")
				key := kv[0]
				val := kv[1]
				err = Set(key, val)
				if err != nil {
					msg := fmt.Sprintf("Failed to set key '%s' with value '%s'. Error: %s\n", key, val, err)
					fmt.Printf(msg)
					conn.Write([]byte(msg))
				} else {
					msg := fmt.Sprintf("Successfully set key '%s' with value '%s'\n", key, val)
					fmt.Printf(msg)
					conn.Write([]byte(msg))
				}
				// todo: send to replica
			} else if args[0] == "get" {
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
		}
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
