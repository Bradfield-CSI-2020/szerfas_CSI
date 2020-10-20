package main

import (
	"io/ioutil"
	"testing"
)

func TestGet4ByteValue(t *testing.T) {
	// 0000 0001, 0000 0010, 0000 0011, 0000 0100, 0000 0101, 0000 0110, 0000 0111, 0000 1000
	// 0000 0010, 0000 0011, 0000 0100, 0000 0101
	// reverse the byte order: -- unfinished -- easier to just use 1
	//sample := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	sample := []byte{1, 1, 0, 0, 0, 0, 0, 0}
	//fmt.Println(sample)
	output := Get4ByteValue(1, sample)
	expected := 1
	if  output != expected {
		t.Errorf("Get4ByteValue should be %d, not %d", expected, output)
	}
}

func TestCountPackets(t *testing.T) {
	dat, err := ioutil.ReadFile("./net.cap")
	check(err)
	packetCount := CountPackets(dat)
	if packetCount != 99 {
		t.Errorf("packetCount should be 99 but is %d", packetCount)
	}
}