package main

import (
	"fmt"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"os"
	"testing"
)

func TestFileScanner(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSize = 1
	parameters.Rng.Seed(1617765204110812000)
	properties := gopter.NewProperties(parameters)
	properties.Property("input to write equals output from read", prop.ForAll(
		func(records [][]string) bool {
			const TABLE_DIR = "/Users/szerfas/src/bradfield/szerfas_CSI/databases/3_physical_storage"
			const TABLE_NAME = "test_table"
			f, err := os.Create(fmt.Sprintf("%s/%s", TABLE_DIR, TABLE_NAME))
			isError(err)
			defer closeFile(f)
			table := Relation{f}
			table.WriteNewBlock(records, 0)
			//answer := table.ReadBlock(0)
			answer := table.Scan()
			//if records == nil && answer == nil {return true}
			//if reflect.DeepEqual(records, make([]string, 0)) && answer == nil {return true}
			//if reflect.DeepEqual(records, make)
			//if records == make([]string, 0) && answer
			return assertBytesEqual(records, answer)
			//if reflect.DeepEqual(records, answer) != true {
			//	printBytes("input: ", records)
			//	printBytes("output: ", answer)
			//	//fmt.Printf("input: %v\nlen(input): %d\noutput: %v\nlen(output): %d\n", records, answer )
			//	//fmt.Printf("input: %v\n output: %v\n", records, answer )
			//	return false
			//} else {
			//	return true
			//}
			//return reflect.DeepEqual(records, answer)
		},
		//TODO: subtract other needed space from BLOCK SIZE and add logic to handle case where string > BLOCK SIZE
		gen.SliceOf(gen.SliceOf(gen.AnyString().SuchThat(func(v string) bool {return len(v) < BLOCK_SIZE}))),
	))
	properties.TestingRun(t)
}

func assertBytesEqual(input [][]string, output[][]string) bool {
	var inputBytes []byte
	var outputBytes []byte
	for i := 0; i < len(input); i++ {
		for x := 0; x < len(input[i]); x++ {
			for z := 0; z < len(input[i][x]); z++ {
				inputBytes = append(inputBytes, input[i][x][z])
				outputBytes = append(outputBytes, output[i][x][z])
				if input[i][x][z] != output[i][x][z] {
					fmt.Printf("input: %v\nlen(input): %d\n", inputBytes, len(inputBytes))
					fmt.Printf("output: %v\nlen(output): %d\n", outputBytes, len(outputBytes))
					return false
				}
			}
		}
	}
	return true
}
//
//func printBytes(varname string, records [][]string) {
//	var bytes []byte
//	for i := 0; i < len(records); i++ {
//		for x := 0; x < len(records[i]); x++ {
//			for z := 0; z < len(records[i][x]); z++ {
//				bytes = append(bytes, records[i][x][z])
//			}
//		}
//	}
//	fmt.Printf("%s: %v\nlen(%s): %d\n", varname, bytes, varname, len(bytes))
//}