package genericavl

import (
	"math/rand"
	"testing"
	"time"

	"github.com/emirpasic/gods/trees/avltree"
)

const (
	randMax = 2000
	nNodes  = 1000
	nDels   = 300
)

var rng *rand.Rand

type IntToString struct {
	key int
	val string
}

type IntToStringTree struct {
	Insert func(*IntToString)
	Delete func(*IntToString)
	Lookup func(*IntToString) (*IntToString, bool)
	Value  func(*Node) *IntToString
}

func (i *IntToStringTree) Compare(a, b *IntToString) int {
	switch {
	case a.key < b.key:
		return -1
	default:
		return 0
	case a.key > b.key:
		return 1
	}
}

type IntTree struct {
	*Tree
	Insert func(int)
	Delete func(int)
	Lookup func(int) (int, bool)
	Value  func(*Node) int
}

func (tree *IntTree) Compare(a, b int) int {
	switch {
	case a < b:
		return -1
	default:
		return 0
	case a > b:
		return 1
	}
}

func (tree *IntTree) SetTree(t *Tree) {
	tree.Tree = t
}

func TestMain(m *testing.M) {
	seed := time.Now().UTC().UnixNano()
	rng = rand.New(rand.NewSource(seed))
	m.Run()
}

func TestCreation(t *testing.T) {
	var tree IntToStringTree
	if err := Make(&tree); err != nil {
		t.Error(err)
	}
}

func TestCreationFails(t *testing.T) {
	type Foo int
	var foo Foo
	if err := Make(&foo); err == nil {
		t.Error("Compare came from nowhere")
	}
}

func TestInsertLookup(t *testing.T) {
	var tree IntToStringTree
	if err := Make(&tree); err != nil {
		t.Error(err)
	}
	tree.Insert(&IntToString{key: 1, val: "one"})
	ret, ok := tree.Lookup(&IntToString{key: 1})
	if !ok {
		t.Error("Could not find element")
	}
	if ret.val != "one" {
		t.Errorf("Did not get the right element: %s\n", ret.val)
	}
}

func TestInsertOrdered(t *testing.T) {
	tree := newRandIntTree(nNodes, randMax, t)
	tree.checkOrdered(t)
}

func newRandIntTree(n, randMax int, t *testing.T) *IntTree {
	var tree IntTree
	if err := Make(&tree); err != nil {
		t.Error(err)
	}

	for i := 0; i < n; i++ {
		tree.Insert(rng.Intn(randMax))
	}
	return &tree
}

func (tree *IntTree) checkOrdered(t *testing.T) {
	n := tree.Min()
	for next := n.Next(); next != nil; next = n.Next() {
		t.Logf("Value in node is %d\n", tree.Value(n))
		if tree.Value(next) <= tree.Value(n) {
			t.Errorf("Tree not ordered: %d ≮ %d", tree.Value(next), tree.Value(n))
		}
		n = next
	}
}

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
	Make(&tree)
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
