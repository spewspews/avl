package avl_test

import (
	"fmt"

	"github.com/spewspews/avl"
)

type StringInt struct {
	key string
	val int
}

type StringIntMap struct {
	Insert func(*StringInt)
	Lookup func(*StringInt) (*StringInt, bool)
}

func (m *StringIntMap) Compare(a, b *StringInt) int {
	switch {
	case a.key < b.key:
		return -1
	default:
		return 0
	case a.key > b.key:
		return 1
	}
}

func ExampleMake() {
	var m StringIntMap
	avl.Make(&m)
	// Type safety: the following will not compile
	// m.insert("foo")
	m.Insert(&StringInt{"foo", 10})
	si, ok := m.Lookup(&StringInt{key: "foo"})
	if ok {
		fmt.Println(si.val)
	}
	// Output:
	// 10
}
