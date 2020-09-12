package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"testing"
)

//func TestBlur(t *testing.T) {
//	f, err := os.Open("input.jpeg")
//	check(err)
//	in, err := jpeg.Decode(f)
//	check(err)
//
//	got := blur(in)
//	_, ok := got.(image.Image)
//	if !ok {
//		t.Errorf("blur() returned non-image")
//	}
//}

func BenchmarkBlur(b *testing.B) {
	f, err := os.Open("input.jpeg")
	check(err)
	in, err := jpeg.Decode(f)
	check(err)

	var out image.Image
	b.ResetTimer()
	// Perform blur
	//for i := 0; i < b.N; i++ {
	//	out = blur(in)
	//}
	out = blur(in)
	b.StopTimer()
	// Write to output
	fout, err := os.Create("output.jpeg")
	check(err)
	jpeg.Encode(fout, out, &jpeg.Options{Quality: 75})

	fmt.Printf("ok\n")
}
