package sandbox

import (
	"math/rand"
	"testing"
)

func randomArray(k int) []int {
	a := make([]int, k)
	for i := 0; i < k; i++ {
		a[i] = rand.Intn(k)
	}
	return a
}

func BenchmarkReflectMap(b *testing.B) {
	a := randomArray(1000)
	for i := 0; i < b.N; i++ {
		ReflectMap(Square, a)
	}
}

func BenchmarkNaiveMap(b *testing.B) {
	a := randomArray(1000)
	for i := 0; i < b.N; i++ {
		NaiveMap(Square, a)
	}
}
