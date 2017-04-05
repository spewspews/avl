package avl_test

import (
	"testing"

	"github.com/emirpasic/gods/trees/avltree"
	"github.com/spewspews/avl"
)

func BenchmarkLookup100(b *testing.B) {
	benchmarkLookup(b, 100)
}

func BenchmarkLookup1000(b *testing.B) {
	benchmarkLookup(b, 1000)
}

func BenchmarkLookup10000(b *testing.B) {
	benchmarkLookup(b, 10000)
}

func BenchmarkLookup100000(b *testing.B) {
	benchmarkLookup(b, 100000)
}

func BenchmarkGoDSGet100(b *testing.B) {
	benchmarkGoDSGet(b, 100)
}

func BenchmarkGoDSGet1000(b *testing.B) {
	benchmarkGoDSGet(b, 1000)
}

func BenchmarkGoDSGet10000(b *testing.B) {
	benchmarkGoDSGet(b, 10000)
}

func BenchmarkGoDSGet100000(b *testing.B) {
	benchmarkGoDSGet(b, 100000)
}

func benchmarkLookup(b *testing.B, size int) {
	b.StopTimer()
	var tree IntTree
	avl.Make(&tree)
	for n := 0; n < size; n++ {
		tree.Insert(n)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			tree.Lookup(n)
		}
	}
}

func benchmarkGoDSGet(b *testing.B, size int) {
	b.StopTimer()
	tree := avltree.NewWithIntComparator()
	for n := 0; n < size; n++ {
		tree.Put(n, nil)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			tree.Get(n)
		}
	}
}
