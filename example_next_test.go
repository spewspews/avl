package avl_test

import (
	"fmt"
	"math/rand"

	"github.com/spewspews/avl"
)

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

func ExampleNode_Next() {
	var t IntTree
	avl.Make(&t)
	for _, i := range rand.Perm(10) {
		t.Insert(i)
	}
	for n := t.Min(); n != nil; n = n.Next() {
		fmt.Println(t.Value(n))
	}

	// Output:
	// 0
	// 1
	// 2
	// 3
	// 4
	// 5
	// 6
	// 7
	// 8
	// 9
}
