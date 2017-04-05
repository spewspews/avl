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
	Delete func(*StringInt)
	Lookup func(*StringInt) (*StringInt, bool)
	Value  func(*avl.Node) *StringInt
	*avl.Tree
}

func (m *StringIntMap) SetTree(t *avl.Tree) {
	m.Tree = t
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

func Example() {
	var m StringIntMap
	avl.Make(&m)
	m.Insert(&StringInt{"foo", 10})
	m.Insert(&StringInt{"bar", 11})
	si, ok := m.Lookup(&StringInt{key: "foo"})
	if ok {
		fmt.Println(si.val)
	}
	si, ok = m.Lookup(&StringInt{key: "bar"})
	if ok {
		fmt.Println(si.val)
	}
	si, ok = m.Lookup(&StringInt{key: "baz"})
	if ok {
		fmt.Println(si.val)
	}
	// Output:
	// 10
	// 11
}
