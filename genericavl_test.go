package genericavl

import (
	"math/rand"
	"testing"
	"time"
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
	Insert func(int)
	Delete func(int)
	Lookup func(int) (int, bool)
	Min    func() *Node
	Max    func() *Node
	Value  func(*Node) int
}

func (i *IntTree) Compare(a, b int) int {
	switch {
	case a < b:
		return -1
	default:
		return 0
	case a > b:
		return 1
	}
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
