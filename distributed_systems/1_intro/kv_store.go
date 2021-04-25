package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	MainLoop2()
	fmt.Println("Goodbye!")
	fmt.Println()
}

var dataStore map[string]string
var storageFile *os.File

func isError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}
	return (err != nil)
}

func MainLoop2() {
	inputBuf := bufio.NewReader(os.Stdin)
	dataStore = make(map[string]string)

	var err error
	storageFile, err = os.OpenFile("kv_store.json", os.O_RDWR|os.O_CREATE, 0755)
	isError(err)
	defer storageFile.Close()

	err = refreshDataStore()
	isError(err)

	for {
		ReceiveInput(inputBuf)
	}
}

func ReceiveInput(inputBuf *bufio.Reader) {
	var line string
	var err error
	fmt.Printf("DISTRIBUTED KV STORE> ")
	line, err = inputBuf.ReadString('\n')
	if err != nil {
		os.Stderr.WriteString("error! please try another command")
		return
	}

	args := strings.Fields(line)

	if len(args) != 2 || (args[0] != "set" && args[0] != "get") {
		os.Stderr.WriteString("Invalid command. Must be of form 'set foo=bar' or 'get foo'. Please try again\n")
		return
	}

	if len(args) == 2 {
		if args[0] == "set" {
			kv := strings.Split(args[1], "=")
			key := kv[0]
			val := kv[1]
			err = Set(key, val)
			if err != nil {
				fmt.Printf("Failed to set key '%s' with value '%s'. Error: %s\n", key, val, err)
			} else {
				fmt.Printf("Successfully set key '%s' with value '%s'\n", key, val)
			}
		} else if args[0] == "get" {
			key := args[1]
			value, err := Get(key)
			if err != nil {
				fmt.Printf("Failed to get key '%s'. Error: %s\n", key, err)
			} else {
				fmt.Printf("%s\n", value)
			}
		}
	}
}

func Set(k string, v string) (err error) {
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
	if err != nil {
		return err
	}
	_, err = storageFile.Write(jsonString)
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
