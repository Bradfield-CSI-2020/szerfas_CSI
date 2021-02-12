//package gopter_exp
//
//import "testing/quick"
//import "github.com/leanovate/gopter"
//
//func main() {
//	quick.Generator()
//	gopter.Shrink()
//}


package main

import (
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

func main() {
	parameters := gopter.DefaultTestParametersWithSeed(1234) // Example should generate reproducible results, otherwise DefaultTestParameters() will suffice

	properties := gopter.NewProperties(parameters)

	properties.Property("fail above 100", prop.ForAll(
		func(arg int64) bool {
			return arg <= 100
		},
		gen.Int64(),
	))

	properties.Property("fail above 100 no shrink", prop.ForAllNoShrink(
		func(arg int64) bool {
			return arg <= 100
		},
		gen.Int64(),
	))

	// When using testing.T you might just use: properties.TestingRun(t)
	properties.Run(gopter.ConsoleReporter(false))
}