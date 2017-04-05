package avl_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/spewspews/avl"
)

const (
	randMax = 2000
	nNodes  = 1000
	nDels   = 300
)

var rng *rand.Rand

type IntTree struct {
	*avl.Tree
	Insert func(int)
	Delete func(int)
	Lookup func(int) (int, bool)
	Value  func(*avl.Node) int
}

func (IntTree) Compare(a, b int) int {
	switch {
	case a < b:
		return -1
	default:
		return 0
	case a > b:
		return 1
	}
}

func (tree *IntTree) SetTree(t *avl.Tree) {
	tree.Tree = t
}

func TestMain(m *testing.M) {
	seed := time.Now().UTC().UnixNano()
	rng = rand.New(rand.NewSource(seed))
	m.Run()
}

func TestCreationFails(t *testing.T) {
	type Foo int
	var foo Foo
	if err := avl.Make(&foo); err == nil {
		t.Error("Compare came from nowhere")
	}
}

func TestInsertOrdered(t *testing.T) {
	tree := newRandIntTree(nNodes, randMax, t)
	tree.checkOrdered(t)
}

func newRandIntTree(n, randMax int, t *testing.T) *IntTree {
	var tree IntTree
	if err := avl.Make(&tree); err != nil {
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
