package avl_test

import (
	"fmt"

	"github.com/spewspews/avl"
)

// StringInt is the key value pair of the mapping
// from strings to integers.
type StringInt struct {
	key string
	val int
}

// StringIntMap implements a map from strings to integers.
// The struct provides functions that will be provided
// by the call to avl.Make.
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
	// This call provides the implementations of
	// StringIntMap.Insert, StringIntMap.Delete, and
	// StringIntMap.Lookup.
	avl.Make(&m)

	// Type safety: the following will not compile
	// m.Insert("foo")
	// StringIntMap.Insert only accepts values of type
	// *StringInt.
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

	// The following does not output anything becuase
	// nothing is found.
	si, ok = m.Lookup(&StringInt{key: "baz"})
	if ok {
		fmt.Println(si.val)
	}

	m.Delete(&StringInt{key: "foo"})
	// The following does not output anything.
	// because the *StringInt with key "foo" has
	// been removed.
	si, ok = m.Lookup(&StringInt{key: "foo"})
	if ok {
		fmt.Println(si.val)
	}
	// Output:
	// 10
	// 11
	// 20
}
