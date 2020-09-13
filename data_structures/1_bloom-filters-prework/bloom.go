package main

import (
	"encoding/binary"
	"hash/fnv"
	"math/big"
	//"github.com/spaolacci/murmur3"
	murmur "github.com/twmb/murmur3"
	"fmt"
)

type bloomFilter interface {
	add(item string)

	// `false` means the item is definitely not in the set
	// `true` means the item might be in the set
	maybeContains(item string) bool

	// Number of bytes used in any underlying storage
	memoryUsage() int
}

type trivialBloomFilter struct {
	data []uint64
}

func newTrivialBloomFilter() *trivialBloomFilter {
	return &trivialBloomFilter{
		data: make([]uint64, 1000),
	}
}

func (b *trivialBloomFilter) add(item string) {
	// Do nothing
}

func (b *trivialBloomFilter) maybeContains(item string) bool {
	// Technically, any item "might" be in the set
	return true
}

func (b *trivialBloomFilter) memoryUsage() int {
	return binary.Size(b.data)
}







type myBloomFilter struct {
	data *big.Int // bit vector
}

//const BLOOMFILTER_SIZE = 280386 	// Elapsed time: 32.021804ms	Memory usage: 35048 bytes 	False positive rate: 32.35%
//const BLOOMFILTER_SIZE = 300386 	// Elapsed time: 36.116905ms 	Memory usage: 37548 bytes 	False positive rate: 29.70%
//const BLOOMFILTER_SIZE = 600000 	// Elapsed time: 31.764326ms 	Memory usage: 75000 bytes 	False positive rate: 10.56% <-- meets requirements
const BLOOMFILTER_SIZE = 800000 	// Elapsed time: 33.212294ms 	Memory usage: 100000 bytes 	False positive rate: 4.55%  <-- maxes size out at 100kb

func newMyBloomFilter() *myBloomFilter {
	var bitVector big.Int
	bitVector.SetBit(&bitVector, BLOOMFILTER_SIZE + 1, 1)  // ensures the bitvector is initialized to the right size - we'll never check this final bit
	fmt.Printf("bit vector initialized to length: %d\n", bitVector.BitLen())

	return &myBloomFilter{
		data: &bitVector,
	}
}

func (b *myBloomFilter) add(item string) {
	fnv, m, blended := b.getHashValues(item)
	b.data.SetBit(b.data, int(fnv % BLOOMFILTER_SIZE), 1)
	b.data.SetBit(b.data, int(m % BLOOMFILTER_SIZE), 1)
	b.data.SetBit(b.data, int(blended % BLOOMFILTER_SIZE), 1)
	if item == "A" {
		fmt.Printf("fnvHashValue for item 'A' is: %d\n", fnv)
		fmt.Printf("murmurHashValue for item 'A' is: %d\n", m)
		fmt.Printf("blendedHashValue for item 'A' is: %d\n", blended)
	}


}

func (b *myBloomFilter) fnvHashValue(item string) uint32 {
	f := fnv.New32a()
	hashResult := []byte(item)
	f.Write(hashResult)
	return f.Sum32()
}

func (b *myBloomFilter) murmurHashValue(item string) uint32 {
	//m := murmur3.New32()
	//hashResult := []byte(item)
	//m.Write(hashResult)
	//return m.Sum32()
	return murmur.StringSum32(item)  // for some reasons this is significantly faster - shaves off almost 100ms when using blended hash vs above code
}

func (b *myBloomFilter) blendedHashValue(item string) uint32 {
	return b.fnvHashValue(item) + b.murmurHashValue(item)
}

func (b *myBloomFilter) maybeContains(item string) bool {
	f, m, blended := b.getHashValues(item)

	isFnvHashSet := b.data.Bit(int(f % BLOOMFILTER_SIZE)) != 0
	isMurmurHashSet := b.data.Bit(int(m % BLOOMFILTER_SIZE)) != 0
	isBlendedHashSet := b.data.Bit(int(blended % BLOOMFILTER_SIZE)) != 0
	if item == "A" {
		fmt.Printf("isFnvHashSet is %t\n", isFnvHashSet)
		fmt.Printf("b.data.Bit(f) is %d\n", b.data.Bit(int(f)))
		fmt.Printf("isMurmurHashSet is %t\n", isMurmurHashSet)
		fmt.Printf("b.data.Bit(m) is %d\n", b.data.Bit(int(m)))
		fmt.Printf("isBlendedHashSet is %t\n", isBlendedHashSet)
		fmt.Printf("b.data.Bit(m) is %d\n", b.data.Bit(int(blended)))
	}
	return  isFnvHashSet && isMurmurHashSet && isBlendedHashSet
}


func (b *myBloomFilter) getHashValues(item string) (fnv uint32, murmur uint32, blended uint32) {

	fnv = b.fnvHashValue(item)
	murmur = b.murmurHashValue(item)
	blended = fnv + murmur

	if item == "A" {
		fmt.Printf("f is %d\n", fnv)
		fmt.Print("m is %d\n", murmur)
		fmt.Print("blended is %d\n", blended)
	}
	return
}


func (b *myBloomFilter) memoryUsage() int {
	return b.data.BitLen() / 8  // will return in bytes
}


