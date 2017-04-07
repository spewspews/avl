// Package avl implements a type-safe generic AVL balanced binary tree.
package avl

import (
	"errors"
	"fmt"
	"reflect"
)

// A Node of the balanced tree.
type Node struct {
	val reflect.Value
	c   [2]*Node
	p   *Node
	b   int8
}

// Setter provides access to the underlying Tree data structure
// by passing the data structure to this interface's SetTree method.
// This provides access to the general Tree methods Min,
// Max, Root, and Size.
type Setter interface {
	SetTree(*Tree)
}

// Tree is the internal representation of the data structure itself.
type Tree struct {
	root     *Node
	elemType reflect.Type
	size     int
	cmp      func(a, b reflect.Value) int8
}

// DummyTree is for documentation purposes only. It is an example
// of the kind of struct that should be passed as a pointer to avl.Make.
type DummyTree struct {
	// The underlying Tree data structure.
	*Tree

	// Insert inserts a new Dummy element into the Tree. Its argument
	// must match the argument type of DummyTree.Compare.
	Insert func(Dummy)

	// Delete deletes a Dummy element from the tree if found.
	Delete func(Dummy)

	// Lookup returns a Dummy element and true if found.
	Lookup func(Dummy) (Dummy, bool)

	// Value returns the Dummy value from the *avl.Node.
	Value func(*Node) Dummy
}

// Compare is used to determine
// how two elements stored in the avl.Tree compare
// and should return an integer less than, equal to, or greater than
// 0 as the two elements are less
// than, equal to, or greater than each other.
// Any struct passed as a reference to avl.Make must have this method defined
// as it is used to deduce the type of the
// elements stored in the Tree as well as in the implementation
// of the tree operations.
func (DummyTree) Compare(a, b Dummy) int {
	return 0
}

// SetTree is passed a pointer to the underlying avl.Tree structure
// to provide access to the non type specific Tree operations
// such as Tree.Root, Tree.Min, Tree.Max, and Tree.Size.
func (d *DummyTree) SetTree(t *Tree) {
	d.Tree = t
}

// Dummy is an empty value for the purposes of documentation.
type Dummy interface{}

// Make creates and provides implementations of type-safe
// balanced binary tree operations. The argument
// TreeStruct must be a pointer to a struct that has a method
// named Compare with the signature
//     func(α, β T) int
// where T is an arbitrary type. Compare should return an integer
// less than, equal to, or greater than 0 depending on whether
// the value α compares less than, equal to, or greater than β,
// respectively. The TreeStruct itself should contain
// fields for functions of the following types:
//    Insert func(T)
//    Delete func(T)
//    Lookup func(T) (T, bool)
//    Value  func(*Node) T
// Make will provide implementations of these functions that
// allow type-safe access to values in the tree. There is no
// error if any of the above functions are missing. See the
// documentation for DummyTree for more information on these
// functions.
//
// If treeStruct implements the Setter interface, then Make will
// pass the underlying Tree data structure to the SetTree method
// to provide access to the non type-specific methods defined on the
// data structure such as, avl.Min, avl.Max, avl.Root, and avl.Size.
// See the documentation for Node.Next for an example.
func Make(treeStruct interface{}) error {
	tsVal := reflect.ValueOf(treeStruct)

	cmp := tsVal.MethodByName("Compare")
	err := checkCompare(cmp)
	if err != nil {
		return err
	}

	t := &Tree{elemType: cmp.Type().In(0)}
	t.cmp = makeCmp(cmp)
	err = t.makeFnImpls(tsVal)
	if err != nil {
		return err
	}

	if setter, ok := treeStruct.(Setter); ok {
		setter.SetTree(t)
	}

	return nil
}

func checkCompare(cmp reflect.Value) error {
	if !cmp.IsValid() {
		return errors.New("Tree interface does not have a Compare method")
	}

	cmpType := cmp.Type()
	if cmpType.Kind() != reflect.Func {
		return errors.New("Compare is not a method")
	}

	if cmpType.NumIn() < 1 {
		return errors.New("Compare method must take two arguments")
	}

	elemType := cmpType.In(0)
	in := []reflect.Type{elemType, elemType}
	out := []reflect.Type{reflect.TypeOf(0)}
	correctType := reflect.FuncOf(in, out, false)
	if cmpType != correctType {
		return fmt.Errorf("Compare method should have signature: %v", correctType)
	}

	return nil
}

func makeCmp(cmp reflect.Value) func(reflect.Value, reflect.Value) int8 {
	args := make([]reflect.Value, 2)
	return func(a, b reflect.Value) int8 {
		args[0] = a
		args[1] = b
		r := cmp.Call(args)[0].Int()
		switch {
		case r < 0:
			return -1
		default:
			return 0
		case r > 0:
			return 1
		}
	}
}

type treeFn struct {
	impl func([]reflect.Value) []reflect.Value
	in   []reflect.Type
	out  []reflect.Type
}

func (t *Tree) makeFnImpls(tsVal reflect.Value) error {
	fns := map[string]treeFn{
		"Insert": {
			t.insert,
			[]reflect.Type{t.elemType},
			[]reflect.Type{},
		},
		"Delete": {
			t.delete,
			[]reflect.Type{t.elemType},
			[]reflect.Type{},
		},
		"Lookup": {
			t.lookup,
			[]reflect.Type{t.elemType},
			[]reflect.Type{t.elemType, reflect.TypeOf(false)},
		},
		"Value": {
			t.value,
			[]reflect.Type{reflect.TypeOf(&Node{})},
			[]reflect.Type{t.elemType},
		},
	}

	for name, tf := range fns {
		fnVal := tsVal.Elem().FieldByName(name)
		if !fnVal.IsValid() {
			continue
		}
		typ := reflect.FuncOf(tf.in, tf.out, false)
		if fnVal.Type() != typ {
			return fmt.Errorf("%s function should have signature: %v", name, typ)
		}
		fnVal.Set(reflect.MakeFunc(typ, tf.impl))
	}

	return nil
}

func (t *Tree) lookup(in []reflect.Value) []reflect.Value {
	val := in[0]
	if val.Type() != t.elemType {
		panic("lookup of wrong type")
	}
	n := t.root
	for n != nil {
		switch t.cmp(val, n.val) {
		case -1:
			n = n.c[0]
		case 0:
			return []reflect.Value{n.val, reflect.ValueOf(true)}
		case 1:
			n = n.c[1]
		}
	}
	return []reflect.Value{reflect.Zero(t.elemType), reflect.ValueOf(false)}
}

func (t *Tree) insert(in []reflect.Value) []reflect.Value {
	val := in[0]
	if val.Type() != t.elemType {
		panic("Inserting wrong type")
	}

	t.insert1(val, nil, &t.root)
	return nil
}

func (t *Tree) insert1(val reflect.Value, p *Node, qp **Node) bool {
	q := *qp
	if q == nil {
		t.size++
		*qp = &Node{val: val, p: p}
		return true
	}

	c := t.cmp(val, q.val)
	if c == 0 {
		q.val = val
		return false
	}

	a := (c + 1) / 2
	fix := t.insert1(val, q, &q.c[a])
	if fix {
		return insertFix(c, qp)
	}
	return false
}

func insertFix(c int8, t **Node) bool {
	s := *t
	if s.b == 0 {
		s.b = c
		return true
	}

	if s.b == -c {
		s.b = 0
		return false
	}

	if s.c[(c+1)/2].b == c {
		s = singlerot(c, s)
	} else {
		s = doublerot(c, s)
	}
	*t = s
	return false
}

func (t *Tree) delete(in []reflect.Value) []reflect.Value {
	val := in[0]
	if val.Type() != t.elemType {
		panic("Deleting wrong type")
	}

	t.delete1(val, &t.root)
	return nil
}

func (t *Tree) delete1(val reflect.Value, qp **Node) bool {
	q := *qp
	if q == nil {
		return false
	}

	c := t.cmp(val, q.val)
	if c == 0 {
		t.size--
		if q.c[1] == nil {
			if q.c[0] != nil {
				q.c[0].p = q.p
			}
			*qp = q.c[0]
			return true
		}
		fix := deleteMin(&q.c[1], &q.val)
		if fix {
			return deleteFix(-1, qp)
		}
		return false
	}
	a := (c + 1) / 2
	fix := t.delete1(val, &q.c[a])
	if fix {
		return deleteFix(-c, qp)
	}
	return false
}

func deleteMin(qp **Node, min *reflect.Value) bool {
	q := *qp
	if q.c[0] == nil {
		*min = q.val
		if q.c[1] != nil {
			q.c[1].p = q.p
		}
		*qp = q.c[1]
		return true
	}
	fix := deleteMin(&q.c[0], min)
	if fix {
		return deleteFix(1, qp)
	}
	return false
}

func deleteFix(c int8, t **Node) bool {
	s := *t
	if s.b == 0 {
		s.b = c
		return false
	}

	if s.b == -c {
		s.b = 0
		return true
	}

	a := (c + 1) / 2
	if s.c[a].b == 0 {
		s = rotate(c, s)
		s.b = -c
		*t = s
		return false
	}

	if s.c[a].b == c {
		s = singlerot(c, s)
	} else {
		s = doublerot(c, s)
	}
	*t = s
	return true
}

func singlerot(c int8, s *Node) *Node {
	s.b = 0
	s = rotate(c, s)
	s.b = 0
	return s
}

func doublerot(c int8, s *Node) *Node {
	a := (c + 1) / 2
	r := s.c[a]
	s.c[a] = rotate(-c, s.c[a])
	p := rotate(c, s)
	if r.p != p || s.p != p {
		panic("doublerot: bad parents")
	}

	switch {
	default:
		s.b = 0
		r.b = 0
	case p.b == c:
		s.b = -c
		r.b = 0
	case p.b == -c:
		s.b = 0
		r.b = c
	}

	p.b = 0
	return p
}

func rotate(c int8, s *Node) *Node {
	a := (c + 1) / 2
	r := s.c[a]
	s.c[a] = r.c[a^1]
	if s.c[a] != nil {
		s.c[a].p = s
	}
	r.c[a^1] = s
	r.p = s.p
	s.p = r
	return r
}

func (t *Tree) value(in []reflect.Value) []reflect.Value {
	n := in[0].Interface().(*Node)
	return []reflect.Value{n.val}
}

// Size returns the number of elements in the tree.
func (t *Tree) Size() int {
	return t.size
}

// Root returns the root node of the tree.
func (t *Tree) Root() *Node {
	return t.root
}

// Min returns the minimum ordered element of the tree.
func (t *Tree) Min() *Node {
	return t.bottom(0)
}

// Max returns the maximum ordered element of the tree.
func (t *Tree) Max() *Node {
	return t.bottom(1)
}

func (t *Tree) bottom(d int) *Node {
	n := t.root
	if n == nil {
		return nil
	}

	for c := n.c[d]; c != nil; c = n.c[d] {
		n = c
	}
	return n
}

// Prev returns the previous Node in an in-order walk
// of the Tree holding the Node n.
func (n *Node) Prev() *Node {
	return n.walk1(0)
}

// Next returns the next Node in an in-order walk
// of the Tree holding the Node n.
func (n *Node) Next() *Node {
	return n.walk1(1)
}

func (n *Node) walk1(a int) *Node {
	if n == nil {
		return nil
	}

	if n.c[a] != nil {
		n = n.c[a]
		for n.c[a^1] != nil {
			n = n.c[a^1]
		}
		return n
	}

	p := n.p
	for p != nil && p.c[a] == n {
		n = p
		p = p.p
	}
	return p
}
