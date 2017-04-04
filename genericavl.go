package genericavl

import (
	"errors"
	"fmt"
	"reflect"
)

type tree struct {
	root                   *Node
	size                   int
	elemType, lookupType   reflect.Type
	insertType, deleteType reflect.Type
	valueType              reflect.Type
	cmp                    func(a, b reflect.Value) int8
}

type Node struct {
	val reflect.Value
	c   [2]*Node
	p   *Node
	b   int8
}

func Make(treeStruct interface{}) error {
	ts := reflect.ValueOf(treeStruct)

	cmp := ts.MethodByName("Compare")
	if err := checkCompare(cmp); err != nil {
		return err
	}

	t := makeTree(cmp)

	tf := treeFn{
		treeFnVal: makeTreeFnVal(ts, "Lookup"),
		impl:      t.lookup,
		typ:       t.lookupType,
	}
	if err := tf.fill(); err != nil {
		return err
	}

	tf = treeFn{
		treeFnVal: makeTreeFnVal(ts, "Delete"),
		impl:      t.delete,
		typ:       t.deleteType,
	}
	if err := tf.fill(); err != nil {
		return err
	}

	tf = treeFn{
		treeFnVal: makeTreeFnVal(ts, "Insert"),
		impl:      t.insert,
		typ:       t.insertType,
	}
	if err := tf.fill(); err != nil {
		return err
	}

	tf = treeFn{
		treeFnVal: makeTreeFnVal(ts, "Value"),
		impl:      t.value,
		typ:       t.valueType,
	}
	if err := tf.fill(); err != nil {
		return err
	}

	min := ts.Elem().FieldByName("Min")
	if min.IsValid() {
		min.Set(reflect.ValueOf(t.min))
	}

	max := ts.Elem().FieldByName("Max")
	if max.IsValid() {
		max.Set(reflect.ValueOf(t.max))
	}

	return nil
}

type treeFn struct {
	treeFnVal
	name string
	impl func([]reflect.Value) []reflect.Value
	typ  reflect.Type
}

type treeFnVal struct {
	name string
	reflect.Value
}

func makeTreeFnVal(ts reflect.Value, name string) treeFnVal {
	return treeFnVal{
		name:  name,
		Value: ts.Elem().FieldByName(name),
	}
}

func (tf *treeFn) fill() error {
	if !tf.IsValid() {
		return nil
	}
	if tf.Type() != tf.typ {
		return fmt.Errorf("%s function has wrong signature: %v", tf.name, tf.Type())
	}
	tf.Set(reflect.MakeFunc(tf.Type(), tf.impl))
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

	if cmpType.NumIn() != 2 || cmpType.NumOut() != 1 {
		return errors.New("Compare method has the wrong signature")
	}

	if cmpType.In(0) != cmpType.In(1) {
		return errors.New("Compare method arguments do not match")
	}

	if kind := cmpType.Out(0).Kind(); kind != reflect.Int {
		return fmt.Errorf("Compare method has wrong return type %v", kind)
	}

	return nil
}

func makeTree(cmp reflect.Value) *tree {
	t := tree{elemType: cmp.Type().In(0)}
	t.insertType = reflect.FuncOf([]reflect.Type{t.elemType}, []reflect.Type{}, false)
	t.deleteType = t.insertType
	t.lookupType = reflect.FuncOf([]reflect.Type{t.elemType}, []reflect.Type{t.elemType, reflect.TypeOf(false)}, false)
	t.valueType = reflect.FuncOf([]reflect.Type{reflect.TypeOf(&Node{})}, []reflect.Type{t.elemType}, false)

	args := make([]reflect.Value, 2)
	t.cmp = func(a, b reflect.Value) int8 {
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

	return &t
}

func (t *tree) lookup(in []reflect.Value) []reflect.Value {
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

func (t *tree) insert(in []reflect.Value) []reflect.Value {
	val := in[0]
	if val.Type() != t.elemType {
		panic("Inserting wrong type")
	}

	t.insert1(val, nil, &t.root)
	return nil
}

func (t *tree) insert1(val reflect.Value, p *Node, qp **Node) bool {
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

func (t *tree) delete(in []reflect.Value) []reflect.Value {
	val := in[0]
	if val.Type() != t.elemType {
		panic("Deleting wrong type")
	}

	t.delete1(val, &t.root)
	return nil
}

func (t *tree) delete1(val reflect.Value, qp **Node) bool {
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

func (t *tree) value(in []reflect.Value) []reflect.Value {
	n := in[0].Interface().(*Node)
	return []reflect.Value{n.val}
}

func (t *tree) min() *Node {
	return t.bottom(0)
}

func (t *tree) max() *Node {
	return t.bottom(1)
}

func (t *tree) bottom(d int) *Node {
	n := t.root
	if n == nil {
		return nil
	}

	for c := n.c[d]; c != nil; c = n.c[d] {
		n = c
	}
	return n
}

func (n *Node) Prev() *Node {
	return n.walk1(0)
}

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
