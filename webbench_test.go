package main

import (
	"cgin/util"
	"testing"
)

func Benchmark_structCopy(b *testing.B) {
	type A struct {
		Name  string
		Age   int
		Attrs []int
	}

	type B struct {
		Name  string
		Age   int
		Attrs []int
	}

	a := &A{
		Name:  "Alice",
		Age:   24,
		Attrs: []int{1, 2, 45},
	}
	bIn := &B{}
	for i := 0; i < b.N; i++ {
		util.StructDeepCopy(a, bIn)
	}
}