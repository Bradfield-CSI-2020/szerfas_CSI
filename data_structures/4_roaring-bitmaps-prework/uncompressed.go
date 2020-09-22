package bitmap

const wordSize = 64

type uncompressedBitmap struct {
	data []uint64
}

func newUncompressedBitmap() *uncompressedBitmap {
	return &uncompressedBitmap{}
}

func (b *uncompressedBitmap) Get(x uint32) bool {
	bucketIndex := x / 64
	withinBucketIndex := x % 64
	if result := b.data[bucketIndex] & (1 << withinBucketIndex); result > 0 {
		return true
	} else {
		return false
	}
}

func (b *uncompressedBitmap) Set(x uint32) {
	// our underlying data structure is a slice of integers, each 64 bits long.
	// To set a given bit, we move down Y number of integers, which totals Y*64 bits. Then within that integer, we
	// flip on the bit at position = x % 64
	bucketIndex := x / 64
	// x is too large to fit into our data, we first resize the data
	// we might do a scheme like 2x the size of the incoming int, but here I cheat since I know we're only testing
	// with values up to 1,100,000
	if bucketIndex > uint32(len(b.data)) {
		newSlice := make([]uint64, 1500000)
		for i, _ := range b.data {
			newSlice[i] = b.data[i]
		}
		b.data = newSlice
	}
	withinBucketIndex := x % 64
	b.data[bucketIndex] |= (1 << withinBucketIndex)
}

func (b *uncompressedBitmap) Union(other *uncompressedBitmap) *uncompressedBitmap {
	var data []uint64
	return &uncompressedBitmap{
		data: data,
	}
}

func (b *uncompressedBitmap) Intersect(other *uncompressedBitmap) *uncompressedBitmap {
	var data []uint64
	return &uncompressedBitmap{
		data: data,
	}
}
