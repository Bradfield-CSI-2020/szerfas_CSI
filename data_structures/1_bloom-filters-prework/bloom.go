package main

import (
	"encoding/binary"
	"hash/fnv"
	"math/big"
	"github.com/spaolacci/murmur3"
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


//Estimating parameters for our bloom filter:
//Goal: <100kb memory usage, 15% false positive rate, speed < few seconds
//Summary of bloom filters (end of article) at http://www.michaelnielsen.org/ddi/why-bloom-filters-work-the-way-they-do/  suggests the following
//
//Choose a desired max probability of a false positive, p
//Choose a ballpark value for n, number of items inserted into the bloom filter
//Choose a value for m, number of bits used, where m = n/ln(2) * log(1/p)
//Calculate the optimal value of k, number of hash functions used, where k = ln(1/p).
//
//for 15% accuracy this is
//	n			m			k				p
//	235886		280,386		1.897119985		0.15
//
//
//Questions to ask:
//which hash functions to use? 		A: Ideally fast and independent; fnv and murmur given as examples
//how many to use?						A: See estimates above. Started with two and made this work, could do more
//separate bit vector or combined?		A: Literature suggests using a combined bit vector is more memory efficient
//what bit vector sizes work well?		A: See equations above



type myBloomFilter struct {
	data *big.Int // bit vector
}

//const BLOOMFILTER_SIZE = 280386 	// Elapsed time: 85.720142ms	Memory usage: 35048 bytes 	False positive rate: 32.35%
//const BLOOMFILTER_SIZE = 300386 	// Elapsed time: 89.715896ms 	Memory usage: 37548 bytes 	False positive rate: 29.70%
const BLOOMFILTER_SIZE = 600000 	// Elapsed time: 81.917925ms 	Memory usage: 75000 bytes 	False positive rate: 10.56% <-- meets requirements!

func newMyBloomFilter() *myBloomFilter {
	var bitVector big.Int
	bitVector.SetBit(&bitVector, BLOOMFILTER_SIZE + 1, 1)  // ensures the bitvector is initialized to the right size - we'll never check this final bit

	return &myBloomFilter{
		data: &bitVector,
	}
}

func (b *myBloomFilter) add(item string) {
	var i uint32

	// apply first hash
	i = b.fnvHashValue(item)
	if item == "A" {
		fmt.Printf("fnvHashValue for item 'A' is: %d\n", i)
	}
	b.data.SetBit(b.data, int(i % BLOOMFILTER_SIZE), 1)

	// apply second hash
	i = b.murmurHashValue(item)
	if item == "A" {
		fmt.Printf("murmurHashValue for item 'A' is: %d\n", i)
	}
	b.data.SetBit(b.data, int(i % BLOOMFILTER_SIZE), 1)

	// apply third hash
	i = b.blendedHashValue(item)
	if item == "A" {
		fmt.Printf("blendedHashValue for item 'A' is: %d\n", i)
	}
	b.data.SetBit(b.data, int(i % BLOOMFILTER_SIZE), 1)
}

func (b *myBloomFilter) fnvHashValue(item string) uint32 {
	f := fnv.New32a()
	hashResult := []byte(item)
	f.Write(hashResult)
	return f.Sum32()
}

func (b *myBloomFilter) murmurHashValue(item string) uint32 {
	m := murmur3.New32()
	hashResult := []byte(item)
	m.Write(hashResult)
	return m.Sum32()
}

func (b *myBloomFilter) blendedHashValue(item string) uint32 {
	return b.fnvHashValue(item) + b.murmurHashValue(item)
}

func (b *myBloomFilter) maybeContains(item string) bool {
	f := int(b.fnvHashValue(item))
	if item == "A" {
		fmt.Printf("f is %d\n", f)
	}
	isFnvHashSet := b.data.Bit(f % BLOOMFILTER_SIZE) != 0
	if item == "A" {
		fmt.Printf("isFnvHashSet is %t\n", isFnvHashSet)
		fmt.Printf("b.data.Bit(f) is %d\n", b.data.Bit(f))
	}
	m := int(b.murmurHashValue(item))
	if item == "A" {
		fmt.Printf("m is %d\n", m)
	}
	isMurmurHashSet := b.data.Bit(m % BLOOMFILTER_SIZE) != 0
	if item == "A" {
		fmt.Printf("isMurmurHashSet is %t\n", isMurmurHashSet)
		fmt.Printf("b.data.Bit(m) is %d\n", b.data.Bit(m))
	}

	// if using third hash func
	blended := int(b.blendedHashValue(item))
	if item == "A" {
		fmt.Printf("blended is %d\n", m)
	}
	isBlendedHashSet := b.data.Bit(blended % BLOOMFILTER_SIZE) != 0
	if item == "A" {
		fmt.Printf("isMurmurHashSet is %t\n", isBlendedHashSet)
		fmt.Printf("b.data.Bit(m) is %d\n", b.data.Bit(blended))
	}
	//return  isFnvHashSet && isMurmurHashSet
	return  isFnvHashSet && isMurmurHashSet && isBlendedHashSet
}

func (b *myBloomFilter) memoryUsage() int {
	//return binary.Size(b.data)
	return b.data.BitLen() / 8  // will return in bytes
}


