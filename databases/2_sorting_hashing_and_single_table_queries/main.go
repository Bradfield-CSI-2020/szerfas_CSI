package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type QueryComponent struct {
	name string
	args []string
}

func getQuery() []QueryComponent {
	// todo: restructure so that can take various input args, not just a single test case
	return []QueryComponent{
		QueryComponent{"PROJECTION", []string{"title"}},
		QueryComponent{"SELECTION", []string{"id", "EQUALS", "5000"}},
		QueryComponent{"SCAN", []string{"movies"}},
	}
}

type Iterator interface {
	Init()
	Next()
	Close()
}

// todo: how can I make this flexible enough to handle data of different types?
type Record struct {
	id int
	title string
	genres []string
}

//func getField(r *Record, field string) interface{} {
//	v := reflect.ValueOf(r)
//	f := reflect.Indirect(v).FieldByName(field)
//	return int(f.Int())
//}

type Scan struct {
	scanData []Record
	offset int
	len int
	source string
}

func (s *Scan) Init() {
	csvfile, err := os.Open(s.source)
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

	r := csv.NewReader(csvfile)
	headers, err := r.Read()
	if err == io.EOF {
		log.Fatal("empty file at that path")
	}
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Headers: %s", headers)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		i, err := strconv.Atoi(record[0])
		genres := strings.Split(record[2], "|")
		s.scanData = append(s.scanData, Record{i, record[1], genres})
	}
	s.len = len(s.scanData)
	s.offset = 0
	// can I close and re-allocate the file descriptor?
}

func (s *Scan) Next() (Record, error) {
	if s.offset >= s.len {
		return Record{}, io.EOF
	}
	old_offset := s.offset
	s.offset += 1
	return s.scanData[old_offset], nil
}

func (s *Scan) Close() {
	// todo: on a live server I'd give thought to if I need to deallocate this or if GC will handle on its own
}

type Criteria struct {
	column string
	operation string
	value string
}

type Selection struct {
	criteria Criteria
	input *Scan
}

func (s *Selection) Init() {
	s.input.Init()
}

func (s *Selection) Next() (Record, error) {
	input, err := s.input.Next()
	if err != nil {
		log.Fatal(err)
	}
	if s.inputPassesCriteria(input) {
		return input, nil
	} else {
		return s.Next()
	}
}

func (s *Selection) inputPassesCriteria(input Record) bool {
	var val interface{}
	switch s.criteria.column {
	case "id":
		val = input.id
	case "title":
		val = input.title
	case "genres":
		val = input.genres
	}

	switch s.criteria.operation {
	case "EQUAL":
		return val == s.criteria.value
	//case "GREATER THAN":
	//	return val > s.criteria.value
	//case "LESS THAN":
	//	return val > s.criteria.value
	default:
		log.Fatal("should not get here :(")
		return false
	}
}

func (s *Selection) Close() {
	// todo: on a live server I'd give thoguht to if I need to deallocate this or if GC will handle on its own
}

func main() {
	query := getQuery()

	for i := 0; i < len(query); i++ {

	}

}
