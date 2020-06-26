package gendemo

// Code generated by go-ipld-prime gengo.  DO NOT EDIT.

import (
	ipld "github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/node/mixins"
	"github.com/ipld/go-ipld-prime/schema"
)

type _Int struct{ x int }
type Int = *_Int

func (n Int) Int() int {
	return n.x
}
func (_Int__Style) FromInt(v int) (Int, error) {
	n := _Int{v}
	return &n, nil
}

type _Int__Maybe struct {
	m schema.Maybe
	v Int
}
type MaybeInt = *_Int__Maybe

func (m MaybeInt) IsNull() bool {
	return m.m == schema.Maybe_Null
}
func (m MaybeInt) IsUndefined() bool {
	return m.m == schema.Maybe_Absent
}
func (m MaybeInt) Exists() bool {
	return m.m == schema.Maybe_Value
}
func (m MaybeInt) AsNode() ipld.Node {
	switch m.m {
	case schema.Maybe_Absent:
		return ipld.Undef
	case schema.Maybe_Null:
		return ipld.Null
	case schema.Maybe_Value:
		return m.v
	default:
		panic("unreachable")
	}
}
func (m MaybeInt) Must() Int {
	if !m.Exists() {
		panic("unbox of a maybe rejected")
	}
	return m.v
}

var _ ipld.Node = (Int)(&_Int{})
var _ schema.TypedNode = (Int)(&_Int{})

func (Int) ReprKind() ipld.ReprKind {
	return ipld.ReprKind_Int
}
func (Int) LookupString(string) (ipld.Node, error) {
	return mixins.Int{"gendemo.Int"}.LookupString("")
}
func (Int) LookupNode(ipld.Node) (ipld.Node, error) {
	return mixins.Int{"gendemo.Int"}.LookupNode(nil)
}
func (Int) LookupIndex(idx int) (ipld.Node, error) {
	return mixins.Int{"gendemo.Int"}.LookupIndex(0)
}
func (Int) LookupSegment(seg ipld.PathSegment) (ipld.Node, error) {
	return mixins.Int{"gendemo.Int"}.LookupSegment(seg)
}
func (Int) MapIterator() ipld.MapIterator {
	return nil
}
func (Int) ListIterator() ipld.ListIterator {
	return nil
}
func (Int) Length() int {
	return -1
}
func (Int) IsUndefined() bool {
	return false
}
func (Int) IsNull() bool {
	return false
}
func (Int) AsBool() (bool, error) {
	return mixins.Int{"gendemo.Int"}.AsBool()
}
func (n Int) AsInt() (int, error) {
	return n.x, nil
}
func (Int) AsFloat() (float64, error) {
	return mixins.Int{"gendemo.Int"}.AsFloat()
}
func (Int) AsString() (string, error) {
	return mixins.Int{"gendemo.Int"}.AsString()
}
func (Int) AsBytes() ([]byte, error) {
	return mixins.Int{"gendemo.Int"}.AsBytes()
}
func (Int) AsLink() (ipld.Link, error) {
	return mixins.Int{"gendemo.Int"}.AsLink()
}
func (Int) Style() ipld.NodeStyle {
	return _Int__Style{}
}

type _Int__Style struct{}

func (_Int__Style) NewBuilder() ipld.NodeBuilder {
	var nb _Int__Builder
	nb.Reset()
	return &nb
}

type _Int__Builder struct {
	_Int__Assembler
}

func (nb *_Int__Builder) Build() ipld.Node {
	if *nb.m != schema.Maybe_Value {
		panic("invalid state: cannot call Build on an assembler that's not finished")
	}
	return nb.w
}
func (nb *_Int__Builder) Reset() {
	var w _Int
	var m schema.Maybe
	*nb = _Int__Builder{_Int__Assembler{w: &w, m: &m}}
}

type _Int__Assembler struct {
	w *_Int
	m *schema.Maybe
}

func (na *_Int__Assembler) reset() {}
func (_Int__Assembler) BeginMap(sizeHint int) (ipld.MapAssembler, error) {
	return mixins.IntAssembler{"gendemo.Int"}.BeginMap(0)
}
func (_Int__Assembler) BeginList(sizeHint int) (ipld.ListAssembler, error) {
	return mixins.IntAssembler{"gendemo.Int"}.BeginList(0)
}
func (na *_Int__Assembler) AssignNull() error {
	switch *na.m {
	case allowNull:
		*na.m = schema.Maybe_Null
		return nil
	case schema.Maybe_Absent:
		return mixins.IntAssembler{"gendemo.Int"}.AssignNull()
	case schema.Maybe_Value, schema.Maybe_Null:
		panic("invalid state: cannot assign into assembler that's already finished")
	}
	panic("unreachable")
}
func (_Int__Assembler) AssignBool(bool) error {
	return mixins.IntAssembler{"gendemo.Int"}.AssignBool(false)
}
func (na *_Int__Assembler) AssignInt(v int) error {
	switch *na.m {
	case schema.Maybe_Value, schema.Maybe_Null:
		panic("invalid state: cannot assign into assembler that's already finished")
	}
	if na.w == nil {
		na.w = &_Int{}
	}
	na.w.x = v
	*na.m = schema.Maybe_Value
	return nil
}
func (_Int__Assembler) AssignFloat(float64) error {
	return mixins.IntAssembler{"gendemo.Int"}.AssignFloat(0)
}
func (_Int__Assembler) AssignString(string) error {
	return mixins.IntAssembler{"gendemo.Int"}.AssignString("")
}
func (_Int__Assembler) AssignBytes([]byte) error {
	return mixins.IntAssembler{"gendemo.Int"}.AssignBytes(nil)
}
func (_Int__Assembler) AssignLink(ipld.Link) error {
	return mixins.IntAssembler{"gendemo.Int"}.AssignLink(nil)
}
func (na *_Int__Assembler) AssignNode(v ipld.Node) error {
	if v.IsNull() {
		return na.AssignNull()
	}
	if v2, ok := v.(*_Int); ok {
		switch *na.m {
		case schema.Maybe_Value, schema.Maybe_Null:
			panic("invalid state: cannot assign into assembler that's already finished")
		}
		if na.w == nil {
			na.w = v2
			*na.m = schema.Maybe_Value
			return nil
		}
		*na.w = *v2
		*na.m = schema.Maybe_Value
		return nil
	}
	if v2, err := v.AsInt(); err != nil {
		return err
	} else {
		return na.AssignInt(v2)
	}
}
func (_Int__Assembler) Style() ipld.NodeStyle {
	return _Int__Style{}
}
func (Int) Type() schema.Type {
	return nil /*TODO:typelit*/
}
func (n Int) Representation() ipld.Node {
	return (*_Int__Repr)(n)
}

type _Int__Repr = _Int

var _ ipld.Node = &_Int__Repr{}

type _Int__ReprStyle = _Int__Style
type _Int__ReprAssembler = _Int__Assembler
