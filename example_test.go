package avl_test

import (
	"fmt"
	"math/rand"

	"github.com/spewspews/avl"
)

type StringInt struct {
	key string
	val int
}

// This struct provides functions that will be filled
// in by the call to Make.
type StringIntMap struct {
	Insert func(*StringInt)
	Delete func(*StringInt)
	Lookup func(*StringInt) (*StringInt, bool)
}

// The comparison method is used both to deduce the type
// of values held in the StringIntMap in order to
// create the functions that are provided to
// StringIntMap and also to perform the necessary
// comparisons. It must have a signature like
//     func(a, b T) int
// where T is the type to be held in the tree.
func (StringIntMap) Compare(a, b *StringInt) int {
	switch {
	case a.key < b.key:
		return -1
	default:
		return 0
	case a.key > b.key:
		return 1
	}
}

func Example() {
	var m StringIntMap
	avl.Make(&m)
	// Type safety: the following will not compile
	// m.insert("foo")
	m.Insert(&StringInt{"foo", 10})
	m.Insert(&StringInt{"bar", 11})

	si, ok := m.Lookup(&StringInt{key: "foo"})
	if ok {
		fmt.Println(si.val)
		si.val = 20
	}
	si, ok = m.Lookup(&StringInt{key: "bar"})
	if ok {
		fmt.Println(si.val)
	}

	// The value should have changed.
	si, ok = m.Lookup(&StringInt{key: "foo"})
	if ok {
		fmt.Println(si.val)
	}

	// The folowing does not output anything.
	si, ok = m.Lookup(&StringInt{key: "baz"})
	if ok {
		fmt.Println(si.val)
	}

	// The folowing does not output anything.
	m.Delete(&StringInt{key: "foo"})
	si, ok = m.Lookup(&StringInt{key: "foo"})
	if ok {
		fmt.Println(si.val)
	}

	// Output:
	// 10
	// 11
	// 20
}

func ExampleNext() {
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
