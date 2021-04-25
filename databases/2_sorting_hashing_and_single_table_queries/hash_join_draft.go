type HashJoinIterator struct {
	// maybe smallTable is already a projection
	smallTable Iterator
	bigTable Iterator // Note: not affiliated with Google

	getKeyFromST func(Tuple) string
	getKeyFromBT func(Tuple) string

	// Takes a tuple from ST, tuple from BT
	// produces a combined tuple for output
	mergeFn func(Tuple, Tuple) Tuple

	hashFn func(string) string
	// by default: hashFn: func(s string) {return s}

	// These are tuples from the smallTable
	hashTable map[string][]Tuple

	// These are merged tuples
	backlog []Tuple
}

func (h *HashJoinIterator) Init() {
	for {
		tuple, ok := h.smallTable.Next()
		if !ok {
			break
		}
		colValue := h.getKeyFromST(tuple)
		h.hashTable[colValue] = append(h.hashTable[colValue], tuple)
	}
}

// Idea: keep a backlog
func (h *HashJoinIterator) Next() (Tuple, bool) {
	if len(h.backlog) > 0 {
		result = h.backlog[0]

		// This is fine in Go (but would be bad
		// in languages where slicing is O(n))
		h.backlog = h.backlog[1:]

		return result, true
	}

	// Small table is in memory
	for {
		// Grab a tuple from bigTable
		tuple, ok := h.bigTable.Next()
		if !ok {
			// No more tuples, we're done
			return Tuple{}, false
		}

		// We have a tuple from bigTable
		// Our goal:
		// match tuples from smallTable with tuples from bigTable
		// that have the same colValue
		colValue := h.getKeyFromBT(tuple)
		key := h.hashFn(colValue)

		matches, ok := h.hashTable[key]
		// Tuple from bigTable didn't match anything,
		// need to go to next tuple from bigTable
		if !ok {
			continue
		}

		// We can assume len(matches) > 0
		// could have multiple matching tuples from small table?
		for _, match := range matches {
			combined := h.mergeFn(match, tuple)
			result = append(result, combined)
		}

		toOutput := result[0]
		h.backlog = append(h.backlog, result[1:]...)
		return toOutput
	}
}