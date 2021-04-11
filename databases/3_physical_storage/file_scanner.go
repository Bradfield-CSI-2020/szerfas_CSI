package main

import (
	"encoding/binary"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

/*
Design simplifications:
- Only handling strings for now
- Only bothering with reading and writing for now, not worrying about updating

File format
BLOCK_SIZE = 4096
// [ Free space <2 bytes> | # records on page | list of (offset, len record) <2 bytes each> | records | ...free space... ]


Files of interest:
3_physical_storage/
	tables/
		<new tables go here>
	catalog/
		map[string][]string storing name: []fieldnames
	// no index for now
	// no log for now


func createTable(name string, fields []string) {
	creates a new file in the tables/ directory
	initializes first set of contents in that file to be in the specified file format with zero records
	updates the catalog to contain the new table's fields
}

type FileScanner struct {
	table os.File
	headers []string
	block int
	recordNumber int
}

func NewFileScanner(tableString string) *FileScanner {
	return &FileScanner{
		table: file at `./tables/tableString`
		headers: headers stored at ./catalog under `tableString`
		block: 0
		recordNumber: 0
	}
}

type Record struct {
	// should this just be a slice of strings instead?
}

//note: would want this more flexible but we're sticking to strings for now, so a record is just a slice of strings
func (fs *FileScanner) Write(record []string) error {
	scan along the file a BLOCK_SIZE at a time, looking for a file with sufficient free space
	once found, add the record to that block, updating amount of freespace, length of header/off
	write to disk at that offset - may need to flush
}

func (fs *FileScanner) Next(record []string) error {
	build Record with fs.block and fs.record, zipping field values with field column names from fs.headers
	increment fs.record
	if fs.record exceeds # of records on that block - 1, increment the block and reset the page
	return Record
}

*/

const BLOCK_SIZE = 256

// not the same as 2_sorting_and_hashing_single_table_queries -- limiting to strings here
type Record struct {
	id string
	title string
	genres string
}


func createFile(path string) {
	var absPath, err = filepath.Abs(path)
	if isError(err) { return }
	// detect if file exists
	_, err = os.Stat(absPath)

	// create file if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		if isError(err) { return }
		defer file.Close()
	}

	fmt.Println("==> done creating file", path)
}

//type Relation struct {
//	f *os.File
//}

// func (*Relation) WriteNewBlock(records []string) error {...writes formattedBlock to disk, recurses through remaining records if necessary}
// func (*Relation) Scan() (records []string) {...repeatedly calls ReadBlock and aggregates results...}
// func (*Relation) ReadBlock(blockNum int) (records []string) {}

type Relation struct {
	w *os.File
}

func NewRelation() {

}

// File format: [ Free space <2 bytes> | # records on page | list of (offset, len record) <4 bytes each> | records | ...free space... ]
// Record format [field_length <1 byte> | field string <field_length bytes> | ...repeat...]

func (r *Relation) WriteNewBlock(records [][]string, blockNum int) {
	//fmt.Printf("blockNum is: %d\n", blockNum)
	//fmt.Printf("records to write: %v\n", records)

	formattedBlock := make([]byte, BLOCK_SIZE)
	// debug log
	//fmt.Printf("empty formattedBlock: %v\n", formattedBlock)
	//fmt.Printf("recrods: %v\n", records)

	var sizeRecords int
	var overflow [][]string
	for i := 0; i < len(records); i++ {
		numFields := len(records[i])
		// todo: figure out how to handle edge case in which a record is empty -- may be easier to bound test (could also bound so that testing a constant num of fields for each iteration)
		//if numFields == 0 {continue}
		var sizeRecord int
		for x := 0; x < numFields; x++ {
			if sizeRecords + sizeRecord + len(records[i][x]) + 1 > BLOCK_SIZE - (4 * len(records) + 4) {
				if sizeRecord + len(records[i][x]) + 1 > BLOCK_SIZE - (4 * len(records) + 4) {
					// this means this record itself is too big for our block size -- we're not investing in this edge case for now so throw an error
					panic("single record is too long to store")
				}
				overflow = records[i:]
				records = records[:i]
				sizeRecord = 0
				break
			}
			// +1 for the "varint" byte that goes in front of every field
			sizeRecord += len(records[i][x]) + 1
		}
		sizeRecords += sizeRecord
	}
	numRecords := len(records)
	//fmt.Printf("numRecords: %d\n", numRecords)
	//fmt.Printf("sizeRecords: %d\n", sizeRecords)
	// if length of bytes is greater than BLOCK_SIZE - header adjustment, return error


	// Free space is calculated by taking the block size minus the following:
	// 2 bytes for freeSpace
	// 2 bytes for number of records on page
	// 4 bytes for each tuple<offset, len of record> in record index
	// X bytes for sizeRecords
	indexSize := 4 * numRecords
	// TODO: recalculate to include byte for field length
	freeSpace := uint16(BLOCK_SIZE - 4 - indexSize - sizeRecords)
	binary.BigEndian.PutUint16(formattedBlock[0:2], freeSpace)

	// debug log
	//fmt.Printf("formattedBlock w/ freespace: %v\n", formattedBlock)

	// add numRecords
	binary.BigEndian.PutUint16(formattedBlock[2:4], uint16(numRecords))

	// debug log
	//fmt.Printf("formattedBlock w/ freespace and numRecords: %v\n", formattedBlock)

	//add <offset,len> for each record
	recordIndexOffset := 4
	// recordOffset starts after
	// 2 bytes for freeSpace
	// 2 bytes for number of records on page
	// 4 bytes for each list of tuple<offset, len of record> in record list
	recordOffset := 4 + indexSize
	//fmt.Printf("numRecords: %d\n", numRecords)
	for i := 0; i < numRecords; i++ {
		numFields := len(records[i])
		// add record offset to index
		binary.BigEndian.PutUint16(formattedBlock[recordIndexOffset:recordIndexOffset+2], uint16(recordOffset))
		//fmt.Printf("added recordOffset: %d to slot: %d\n", recordOffset, recordIndexOffset)

		fieldStart := recordOffset
		for x := 0; x < numFields; x++ {
			// add record
			fieldLen := len(records[i][x])
			// add length of field before each field
			formattedBlock[recordOffset] = byte(fieldLen)
			recordOffset++
			for z := 0; z < fieldLen; z++ {
				if recordOffset == BLOCK_SIZE {
					fmt.Printf("inside trouble")
				}
				formattedBlock[recordOffset] = records[i][x][z]
				recordOffset ++
			}
		}

		// add record length to index
		recLen := (recordOffset - fieldStart)
		binary.BigEndian.PutUint16(formattedBlock[recordIndexOffset+2:recordIndexOffset+4], uint16(recLen))
		//fmt.Printf("added recLen: %d to slot: %d\n", recLen, recordIndexOffset+2)

		// increment index offset
		recordIndexOffset += 4

		//debug log
		//fmt.Printf("formattedBlock w/ %d records added: %v\n", i + 1, formattedBlock)
		//fmt.Printf("length of formattedBlock: %d\n", len(formattedBlock))
	}

	// write to the file at offset = blockNum * BLOCK_SIZE
	_, err := r.w.Seek(int64(blockNum * BLOCK_SIZE), 0)
	isError(err)
	//bytes, err := r.w.Write(formattedBlock)
	_, err = r.w.Write(formattedBlock)
	isError(err)
	//fmt.Printf("bytes written: %d\n", bytes)
	//fmt.Printf("written formatted block: %v\n", formattedBlock)
	if len(overflow) > 0 {
		fmt.Printf("overflow: %v\n", overflow)

		r.WriteNewBlock(overflow, blockNum + 1)
	}
}

func (r *Relation) Scan() (records [][]string) {
	_, err := r.w.Seek(int64(0), 0)
	isError(err)
	stat, err := r.w.Stat()
	isError(err)
	size := stat.Size()
	blocks := size / BLOCK_SIZE
	var i int64
	for i = 0; i < blocks; i++ {
		records = append(records, r.ReadBlock(int(i))...)
	}
	return records
}


func (r *Relation) ReadBlock(blockNum int) (records [][]string) {
	//records = make([][]string, 0)
	_, err := r.w.Seek(int64(blockNum * BLOCK_SIZE), 0)
	isError(err)
	formattedBlock := make([]byte, BLOCK_SIZE)
	//bytesRead, err := r.w.Read(formattedBlock)
	_, err = r.w.Read(formattedBlock)
	isError(err)
	//fmt.Printf("bytesRead: %d\n", bytesRead)
	//fmt.Printf("formattedBlock: %v\n", formattedBlock)
	numRecords := int(binary.BigEndian.Uint16(formattedBlock[2:4]))
	var recordOffset uint16
	var recordLen uint16
	// start at offset of index, which is 4
	for i := 4; i < (numRecords * 4) + 4 ; i += 4 {
		recordOffset = binary.BigEndian.Uint16(formattedBlock[i:i+2])
		recordLen = binary.BigEndian.Uint16(formattedBlock[i+2:i+4])
		p := recordOffset
		var record []string
		var field string
		for p < (recordOffset + recordLen) {
			fieldLen := uint16(formattedBlock[p])
			p++
			field = string(formattedBlock[p:p+fieldLen])
			p += fieldLen
			record = append(record, field)
		}
		//recordBytes = formattedBlock[recordOffset:recordOffset + recordLen]
		records = append(records, record)
	}
	//fmt.Printf("records read: %v\n", records)
	return records
}

func closeFile(file *os.File) {
	err := file.Close(); if err != nil {
		log.Fatal(err)
	}
}

func getTableData (path string) (records [][]string, err error) {
	csvfile, err := os.Open(path)
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}
	defer closeFile(csvfile)

	r := csv.NewReader(csvfile)
	//headers, err := r.Read()
	_, err = r.Read()
	if err == io.EOF {
		log.Fatal("empty file at that path")
	}
	if err != nil {
		log.Fatal(err)
	}

	// put headers into separate table
	//fmt.Printf("Headers: %s\n", headers)

	for {
		record, err := r.Read()
		//fmt.Printf("Record: %v\n", record)
		fmt.Printf("err: %s\n", err)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		records = append(records, record)
	}
	return records, nil
}

func main() {
	const TABLE_DIR = "/Users/szerfas/src/bradfield/szerfas_CSI/databases/3_physical_storage"
	const TABLE_NAME = "test_table"
	f, err := os.Create(fmt.Sprintf("%s/%s", TABLE_DIR, TABLE_NAME))
	isError(err)
	defer closeFile(f)
	records, err:= getTableData("/Users/szerfas/src/bradfield/szerfas_CSI/databases/3_physical_storage/ml-20m/10movies.csv")
	//records := [][]string{[]string{"movieId", "title", "genre"}, []string{"some id", "some title", "some genre"}}
	isError(err)
	table := Relation{f}
	table.WriteNewBlock(records, 2)
	answer := table.ReadBlock(2)
	fmt.Println(answer)
}



func createTable(name string, fields []string) {
	// creates table in table directory, initializing first block to empty (but adhering to file format)
	// adds list of fields in catalog directory
	// could store catalog directory in same file format as other tables, but want to move faster for now
}

type FileWriter struct {
	csv *os.File
}

func (fw *FileWriter) createTable(name string, fields []string) {
	// create new file in table directory
	// if file exists, throw error

	// set offset to zero
	// read first line of fw.csv
	// block := DBBLock{records: <insert records here>}
	// write headers to catalog file <-- this can be in same specified format of file if have reusable code

	// loop through remaining records
	// write each new record in specified page format
	// when block size is sufficient, write to disk and start over
}


// implements Iterator interface
// todo: connect files across sub-folders to work together
type FileScanner struct {
	header []string
	table *os.File
	page int
	recordNum int
}

func NewFileScanner(table string) (f *FileScanner, e error) {
	// table is a path to a file on disk, such as `./movies.table`
	path, err := filepath.Abs(table)
	if isError(err) { return nil, err}
	createFile(path)

	// lookup header from header table
	//header

	return &FileScanner{}, nil
}

func (fs *FileScanner) Write(r Record, path string) error {
	// for empty space sufficient for a record by looping through records written so far
	// once find one,
	// create block to write into memory
	return nil
}


func (fs *FileScanner) Next() (r Record, e error) {

	if fs.header == nil {

	}
	// pulls next record from the file

	// if e == EOF, means there is no Next record

	// matches field names from columns with values read off the file
	// returns next record
	return Record{}, nil
}


func isError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}

	return (err != nil)
}