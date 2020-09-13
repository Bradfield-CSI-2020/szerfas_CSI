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
//const BLOOMFILTER_SIZE = 600000 	// Elapsed time: 81.917925ms 	Memory usage: 75000 bytes 	False positive rate: 10.56% <-- meets requirements!
const BLOOMFILTER_SIZE = 800000

func newMyBloomFilter() *myBloomFilter {
	var bitVector big.Int
	bitVector.SetBit(&bitVector, BLOOMFILTER_SIZE + 1, 1)  // ensures the bitvector is initialized to the right size - we'll never check this final bit
	fmt.Printf("bit vector initialized to length: %d\n", bitVector.BitLen())

	return &myBloomFilter{
		data: &bitVector,
	}
}

func (b *myBloomFilter) add(item string) {
	// apply first hash
	fnv := b.fnvHashValue(item)
	if item == "A" {
		fmt.Printf("fnvHashValue for item 'A' is: %d\n", fnv)
	}
	b.data.SetBit(b.data, int(fnv % BLOOMFILTER_SIZE), 1)

	// apply second hash
	m := b.murmurHashValue(item)
	if item == "A" {
		fmt.Printf("murmurHashValue for item 'A' is: %d\n", m)
	}
	b.data.SetBit(b.data, int(m % BLOOMFILTER_SIZE), 1)

	// apply third hash
	//blended := b.blendedHashValue(item)
	blended := fnv + m
	if item == "A" {
		fmt.Printf("blendedHashValue for item 'A' is: %d\n", blended)
	}
	b.data.SetBit(b.data, int(blended % BLOOMFILTER_SIZE), 1)
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
	return murmur.StringSum32(item)  // for some reasons this is significantly faster - shaves off almost 100ms when using blended hash
}

func (b *myBloomFilter) blendedHashValue(item string) uint32 {
	return b.fnvHashValue(item) + b.murmurHashValue(item)
}

func (b *myBloomFilter) maybeContains(item string) bool {
	f := b.fnvHashValue(item)
	if item == "A" {
		fmt.Printf("f is %d\n", f)
	}
	isFnvHashSet := b.data.Bit(int(f % BLOOMFILTER_SIZE)) != 0
	if item == "A" {
		fmt.Printf("isFnvHashSet is %t\n", isFnvHashSet)
		fmt.Printf("b.data.Bit(f) is %d\n", b.data.Bit(int(f)))
	}
	m := b.murmurHashValue(item)
	if item == "A" {
		fmt.Printf("m is %d\n", m)
	}
	isMurmurHashSet := b.data.Bit(int(m % BLOOMFILTER_SIZE)) != 0
	if item == "A" {
		fmt.Printf("isMurmurHashSet is %t\n", isMurmurHashSet)
		fmt.Printf("b.data.Bit(m) is %d\n", b.data.Bit(int(m)))
	}

	// if using third hash func
	//blended := int(b.blendedHashValue(item))
	blended := f + m
	if item == "A" {
		fmt.Printf("blended is %d\n", blended)
	}
	isBlendedHashSet := b.data.Bit(int(blended % BLOOMFILTER_SIZE)) != 0
	if item == "A" {
		fmt.Printf("isMurmurHashSet is %t\n", isBlendedHashSet)
		fmt.Printf("b.data.Bit(m) is %d\n", b.data.Bit(int(blended)))
	}
	//return  isFnvHashSet && isMurmurHashSet
	return  isFnvHashSet && isMurmurHashSet && isBlendedHashSet
}

func (b *myBloomFilter) memoryUsage() int {
	//return binary.Size(b.data)
	return b.data.BitLen() / 8  // will return in bytes
}


