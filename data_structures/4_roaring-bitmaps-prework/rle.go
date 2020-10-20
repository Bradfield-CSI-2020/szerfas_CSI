package bitmap

import "math"

func compress(b *uncompressedBitmap) []uint64 {
	var (
		result = make([]uint64, len(b.data))
		compressedValue uint64
		resultIndex uint64 = 0
	)
	// loop through the data
	// if a value is == 0, or if a value is == Max.Int64, then
		//call a helper function(index + 1, data, bool) (counter):
			// loop through data starting at that index
			// increment counter on each loop
			// if bool is false, break as soon as value is not set to 0; if true, break as soon as value is not set to Max.Int64
			// return counter
		//set first bit of uint64 at current index to be 1
		//set second bit to be 0 if value is 0 or 1 if value is Max.Int64
		//bitwise union this first+second-bit-set uint64 with the counter
		// jump i the length of the counter
	// otherwise just copy the value over
	// return result

	for index := 0; index < len(b.data); index++ {
		val := b.data[index]
		if val == 0 {
			rowsToCompress := countRowsToCompress(index, b.data, false )
			compressedValue = 1 << 64
			compressedValue |= rowsToCompress
			result[resultIndex] = compressedValue
			index += int(rowsToCompress) - 1
		} else if val == 1 {
			rowsToCompress := countRowsToCompress(index + 1, b.data, true )
			compressedValue = 1 << 64
			compressedValue |= 1 << 63
			compressedValue |= rowsToCompress
			result[resultIndex] = compressedValue
			index += int(rowsToCompress) - 1
		} else {
			result[resultIndex] = b.data[index]
			resultIndex++
		}
	}
	// todo: resize the underlying array behind this result down to length of the result slice
	return result
}

func countRowsToCompress(index int, data []uint64, seekOnesOrZeroes bool) (counter uint64) {
	counter = 0
	for {
		if seekOnesOrZeroes {
			if data[index] == math.MaxUint64 {
				counter++
				index++
			} else {
				return counter
			}
		} else {
			if data[index] == 0 {
				counter++
				index++
			} else {
				return counter
			}
		}
	}
}

func doubleSliceSize(data []uint64) {
	newSlice := make([]uint64, (len(data) + 1) * 2)
	for i, _ := range data {
		newSlice[i] = data[i]
	}
	data = newSlice
}

func decompress(compressed []uint64) *uncompressedBitmap {
	var data []uint64

	// loop through compressed data
	// if a value is == 0, or if a value is == Max.Int64, then
	//call a helper function(index + 1, data, bool) (counter):
	// loop through data starting at that index
	// increment counter on each loop
	// if bool is false, break as soon as value is not set to 0; if true, break as soon as value is not set to Max.Int64
	// return counter
	//set first bit of uint64 at current index to be 1
	//set second bit to be 0 if value is 0 or 1 if value is Max.Int64
	//bitwise union this first+second-bit-set uint64 with the counter
	// jump i the length of the counter
	// otherwise just copy the value over
	// return result

	return &uncompressedBitmap{
		data: data,
	}
}
