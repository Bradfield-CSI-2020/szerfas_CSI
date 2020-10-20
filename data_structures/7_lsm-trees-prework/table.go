package table

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"
	"unsafe"
)

type Item struct {
	Key, Value string
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type DiskBlock struct {
	Size int
	ItemMap map[string]string
}

func flushItemsToDisk(path string, sortedItems []Item) {
	//var stringsToFlush = []byte
	size := len(sortedItems)
	itemMap := make(map[string]string, size)
	for _, item := range sortedItems {
		itemMap[item.Key] = item.Value
	}

	diskBlock := DiskBlock{size, itemMap}

	jsonBlob, err := json.Marshal(diskBlock)
	check(err)
	//for _, item := range sortedItems {
	//	stringsToFlush = append(stringsToFlush, []byte(item.Key)...)
	//	stringsToFlush = append(stringsToFlush, []byte(item.Value)...)
	//}
	err = ioutil.WriteFile(path, jsonBlob, 0644)
	check(err)
	fmt.Printf("flushed to path: %s,\n JSON: %s\n", path, jsonBlob)
}

func getItemsFromDisk(path string) DiskBlock {
	jsonBlob, err := ioutil.ReadFile(path)
	check(err)
	var result DiskBlock
	err = json.Unmarshal(jsonBlob, &result)
	check(err)
	fmt.Println(result)
	return result
}

func flushIndexToDisk(path string, index []string) {
	jsonBlob, err := json.Marshal(index)
	check(err)
	err = ioutil.WriteFile(path, jsonBlob, 0644)
	check(err)
}

func getIndexFromDisk(path string) []string {
	jsonBlob, err := ioutil.ReadFile(path)
	check(err)
	var result []string
	err = json.Unmarshal(jsonBlob, &result)
	check(err)
	fmt.Println("index from disk")
	fmt.Println(result)
	return result
}

func resizeIndex(index []string) {

}

// Given a sorted list of key/value pairs, write them out according to the format you designed.
func Build(path string, sortedItems []Item) error {
	var (
		n = len(sortedItems)
		index = make([]string, n) // todo: find way to make initial allocation smaller
		byteCount uintptr = 0
		myMacBlock uintptr = 4096
		blockCounter = 0
		lastItemIndexToDisk = 0
		dir = filepath.Dir(path)
	)
	// todo: handle edge case where key + string combo > block size
	fmt.Println("inside Build")
	for i, item := range sortedItems {
		byteCount += unsafe.Sizeof(item)
		if byteCount > myMacBlock {
			fmt.Printf("flushing to disk\n")
			// flush to disc sortedItems[blockCounter:i]
			flushItemsToDisk(filepath.Join(dir, fmt.Sprint(blockCounter)), sortedItems[lastItemIndexToDisk:i - 1])
			// put the last item of that file in the index
			fmt.Printf("last sorted items key: %s\n", sortedItems[i-1].Key)
			index[blockCounter] = sortedItems[i-1].Key
			fmt.Println("index:")
			fmt.Println(index)
			// increment file counter
			blockCounter++
			// set most recent item to disk
			lastItemIndexToDisk = i - 1
			// reset bytecount
			byteCount = 0
		}
		fmt.Printf("ending loop %d\n", i)
	}
	flushItemsToDisk(filepath.Join(dir, fmt.Sprint(blockCounter)), sortedItems[lastItemIndexToDisk:])
	fmt.Printf("len(index): %d\n", len(index))
	index = append([]string(nil), index[:blockCounter]...)
	fmt.Printf("len(index): %d\n", len(index))
	flushIndexToDisk(filepath.Join(dir, "index"), index)
	fmt.Printf("index flushed to disk\n")

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("files in directory:")
	for _, f := range files {
		fmt.Println(f.Name())
	}

	return nil
}

// A Table provides efficient access into sorted key/value data that's organized according
// to the format you designed.
//
// Although a Table shouldn't keep all the key/value data in memory, it should contain
// some metadata to help with efficient access (e.g. size, index, optional Bloom filter).
type Table struct {
	//size int
	Path string
	Index []string
	// todo: bloom filter
}

// Prepares a Table for efficient access. This will likely involve reading some metadata
// in order to populate the fields of the Table struct.
func LoadTable(path string) (*Table, error) {
	dir := filepath.Dir(path)
	index := getIndexFromDisk(filepath.Join(dir, "/index"))
	return &Table{dir, index}, nil
}

func (t *Table) Get(key string) (string, bool, error) {
	// todo: figure out this search issue
	blockNumber := sort.SearchStrings(t.Index, key)
	fmt.Printf("key is: %s\n", key)
	fmt.Printf("block number is: %d\n", blockNumber)
	fmt.Printf("len(t.Index): %d\n", len(t.Index))
	fmt.Println(t.Index)

	blockPath := filepath.Join(t.Path, fmt.Sprint(blockNumber))
	block := getItemsFromDisk(blockPath)
	if result, ok := block.ItemMap[key]; ok {
		return result, true, nil
	} else {
		return "", false, nil
	}
}

func (t *Table) RangeScan(startKey, endKey string) (Iterator, error) {
	return nil, nil
}

type Iterator interface {
	// Advances to the next item in the range. Assumes Valid() == true.
	Next()

	// Indicates whether the iterator is currently pointing to a valid item.
	Valid() bool

	// Returns the Item the iterator is currently pointing to. Assumes Valid() == true.
	Item() Item
}
