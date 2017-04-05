package avl_test

import (
	"fmt"
	"math/rand"

	"github.com/spewspews/avl"
)

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
