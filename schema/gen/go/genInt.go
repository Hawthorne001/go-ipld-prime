package gengo

import (
	"io"

	"github.com/ipld/go-ipld-prime/schema"
	"github.com/ipld/go-ipld-prime/schema/gen/go/mixins"
)

type intGenerator struct {
	AdjCfg *AdjunctCfg
	mixins.IntTraits
	PkgName string
	Type    schema.TypeInt
}

// --- native content and specializations --->

func (g intGenerator) EmitNativeType(w io.Writer) {
	// Using a struct with a single member is the same size in memory as a typedef,
	//  while also having the advantage of meaning we can block direct casting,
	//   which is desirable because the compiler then ensures our validate methods can't be evaded.
	doTemplate(`
		type _{{ .Type | TypeSymbol }} struct{ x int }
		type {{ .Type | TypeSymbol }} = *_{{ .Type | TypeSymbol }}
	`, w, g.AdjCfg, g)
}

func (g intGenerator) EmitNativeAccessors(w io.Writer) {
	// The node interface's `AsInt` method is almost sufficient... but
	//  this method unboxes without needing to return an error that's statically impossible,
	//   which makes it easier to use in chaining.
	doTemplate(`
		func (n {{ .Type | TypeSymbol }}) Int() int {
			return n.x
		}
	`, w, g.AdjCfg, g)
}

func (g intGenerator) EmitNativeBuilder(w io.Writer) {
	// Scalar types are easy to generate a constructor function for.
	// REVIEW: if this is useful and should be on by default; it also adds a decent amount of noise to a package.
	//  It might make perfect sense to put it on the NodeStyle instead, as `FromInt(int)`.
	// FUTURE: should engage validation flow.
	doTemplate(`
		func New{{ .Type | TypeSymbol }}(v int) {{ .Type | TypeSymbol }} {
			n := _{{ .Type | TypeSymbol }}{v}
			return &n
		}
	`, w, g.AdjCfg, g)
}

func (g intGenerator) EmitNativeMaybe(w io.Writer) {
	// REVIEW: can this be extracted to the mixins package?  it doesn't appear to vary for kind.
	doTemplate(`
		type _{{ .Type | TypeSymbol }}__Maybe struct {
			m schema.Maybe
			v {{if not (MaybeUsesPtr .Type) }}_{{end}}{{ .Type | TypeSymbol }}
		}
		type Maybe{{ .Type | TypeSymbol }} = *_{{ .Type | TypeSymbol }}__Maybe

		func (m Maybe{{ .Type | TypeSymbol }}) IsNull() bool {
			return m.m == schema.Maybe_Null
		}
		func (m Maybe{{ .Type | TypeSymbol }}) IsUndefined() bool {
			return m.m == schema.Maybe_Absent
		}
		func (m Maybe{{ .Type | TypeSymbol }}) Exists() bool {
			return m.m == schema.Maybe_Value
		}
		func (m Maybe{{ .Type | TypeSymbol }}) AsNode() ipld.Node {
			switch m.m {
				case schema.Maybe_Absent:
					return ipld.Undef
				case schema.Maybe_Null:
					return ipld.Null
				case schema.Maybe_Value:
					return {{if not (MaybeUsesPtr .Type) }}&{{end}}m.v
				default:
					panic("unreachable")
			}
		}
		func (m Maybe{{ .Type | TypeSymbol }}) Must() {{ .Type | TypeSymbol }} {
			if !m.Exists() {
				panic("unbox of a maybe rejected")
			}
			return {{if not (MaybeUsesPtr .Type) }}&{{end}}m.v
		}
	`, w, g.AdjCfg, g)
}

// --- type info --->

func (g intGenerator) EmitTypeConst(w io.Writer) {
	doTemplate(`
		// TODO EmitTypeConst
	`, w, g.AdjCfg, g)
}

// --- TypedNode interface satisfaction --->

func (g intGenerator) EmitTypedNodeMethodType(w io.Writer) {
	doTemplate(`
		func ({{ .Type | TypeSymbol }}) Type() schema.Type {
			return nil /*TODO:typelit*/
		}
	`, w, g.AdjCfg, g)
}

func (g intGenerator) EmitTypedNodeMethodRepresentation(w io.Writer) {
	// Perhaps surprisingly, the way to get the representation node pointer
	//  does not actually depend on what the representation strategy is.
	doTemplate(`
		func (n {{ .Type | TypeSymbol }}) Representation() ipld.Node {
			return (*_{{ .Type | TypeSymbol }}__Repr)(n)
		}
	`, w, g.AdjCfg, g)
}

// --- Node interface satisfaction --->

func (g intGenerator) EmitNodeType(w io.Writer) {
	// No additional types needed.  Methods all attach to the native type.
}

func (g intGenerator) EmitNodeTypeAssertions(w io.Writer) {
	doTemplate(`
		var _ ipld.Node = ({{ .Type | TypeSymbol }})(&_{{ .Type | TypeSymbol }}{})
		var _ schema.TypedNode = ({{ .Type | TypeSymbol }})(&_{{ .Type | TypeSymbol }}{})
	`, w, g.AdjCfg, g)
}

func (g intGenerator) EmitNodeMethodAsInt(w io.Writer) {
	doTemplate(`
		func (n {{ .Type | TypeSymbol }}) AsInt() (int, error) {
			return n.x, nil
		}
	`, w, g.AdjCfg, g)
}

func (g intGenerator) EmitNodeMethodStyle(w io.Writer) {
	doTemplate(`
		func ({{ .Type | TypeSymbol }}) Style() ipld.NodeStyle {
			return _{{ .Type | TypeSymbol }}__Style{}
		}
	`, w, g.AdjCfg, g)
}

func (g intGenerator) EmitNodeStyleType(w io.Writer) {
	doTemplate(`
		type _{{ .Type | TypeSymbol }}__Style struct{}

		func (_{{ .Type | TypeSymbol }}__Style) NewBuilder() ipld.NodeBuilder {
			var nb _{{ .Type | TypeSymbol }}__Builder
			nb.Reset()
			return &nb
		}
	`, w, g.AdjCfg, g)
}

// --- NodeBuilder and NodeAssembler --->

func (g intGenerator) GetNodeBuilderGenerator() NodeBuilderGenerator {
	return intBuilderGenerator{
		g.AdjCfg,
		mixins.IntAssemblerTraits{
			g.PkgName,
			g.TypeName,
			"_" + g.AdjCfg.TypeSymbol(g.Type) + "__",
		},
		g.PkgName,
		g.Type,
	}
}

type intBuilderGenerator struct {
	AdjCfg *AdjunctCfg
	mixins.IntAssemblerTraits
	PkgName string
	Type    schema.TypeInt
}

func (g intBuilderGenerator) EmitNodeBuilderType(w io.Writer) {
	doTemplate(`
		type _{{ .Type | TypeSymbol }}__Builder struct {
			_{{ .Type | TypeSymbol }}__Assembler
		}
	`, w, g.AdjCfg, g)
}
func (g intBuilderGenerator) EmitNodeBuilderMethods(w io.Writer) {
	doTemplate(`
		func (nb *_{{ .Type | TypeSymbol }}__Builder) Build() ipld.Node {
			return nb.w
		}
		func (nb *_{{ .Type | TypeSymbol }}__Builder) Reset() {
			var w _{{ .Type | TypeSymbol }}
			*nb = _{{ .Type | TypeSymbol }}__Builder{_{{ .Type | TypeSymbol }}__Assembler{w: &w}}
		}
	`, w, g.AdjCfg, g)
	// Cover up all the assembler methods that are prepared to handle null;
	//  this saves them from having to check for a nil 'fcb'.
	doTemplate(`
		func (nb *_{{ .Type | TypeSymbol }}__Builder) AssignNull() error {
			return mixins.StringAssembler{"{{ .PkgName }}.{{ .TypeName }}"}.AssignNull()
		}
		func (nb *_{{ .Type | TypeSymbol }}__Builder) AssignInt(v int) error {
			*nb.w = _{{ .Type | TypeSymbol }}{v}
			return nil
		}
		func (nb *_{{ .Type | TypeSymbol }}__Builder) AssignNode(v ipld.Node) error {
			if v2, err := v.AsInt(); err != nil {
				return err
			} else {
				return nb.AssignInt(v2)
			}
		}
	`, w, g.AdjCfg, g)
}
func (g intBuilderGenerator) EmitNodeAssemblerType(w io.Writer) {
	doTemplate(`
		type _{{ .Type | TypeSymbol }}__Assembler struct {
			w *_{{ .Type | TypeSymbol }}
			z bool
			fcb func() error
		}
	`, w, g.AdjCfg, g)
}
func (g intBuilderGenerator) EmitNodeAssemblerMethodAssignNull(w io.Writer) {
	// Twist: All generated assemblers quietly accept null... and if used in a context they shouldn't,
	//  either this method gets overriden (this is the case for NodeBuilders),
	//  or the 'na.fcb' (short for "finish callback") returns the rejection.
	//  We don't need a nil check for 'fcb' because all parent assemblers use it, and root builders override all relevant methods.
	//  We don't pass any args to 'fcb' because we assume it comes from something that can already see this whole struct.
	doTemplate(`
		func (na *_{{ .Type | TypeSymbol }}__Assembler) AssignNull() error {
			na.z = true
			return na.fcb()
		}
	`, w, g.AdjCfg, g)
}
func (g intBuilderGenerator) EmitNodeAssemblerMethodAssignInt(w io.Writer) {
	// This method contains a branch to support MaybeUsesPtr because new memory may need to be allocated.
	//  This allocation only happens if the 'w' ptr is nil, which means we're being used on a Maybe;
	//  otherwise, the 'w' ptr should already be set, and we fill that memory location without allocating, as usual.
	doTemplate(`
		func (na *_{{ .Type | TypeSymbol }}__Assembler) AssignInt(v int) error {
			{{- if .Type | MaybeUsesPtr }}
			if na.w == nil {
				na.w = &_{{ .Type | TypeSymbol }}{v}
				return na.fcb()
			}
			{{- end}}
			*na.w = _{{ .Type | TypeSymbol }}{v}
			return na.fcb()
		}
	`, w, g.AdjCfg, g)
}
func (g intBuilderGenerator) EmitNodeAssemblerMethodAssignNode(w io.Writer) {
	doTemplate(`
		func (na *_{{ .Type | TypeSymbol }}__Assembler) AssignNode(v ipld.Node) error {
			if v.IsNull() {
				return na.AssignNull()
			}
			if v2, err := v.AsInt(); err != nil {
				return err
			} else {
				return na.AssignInt(v2)
			}
		}
	`, w, g.AdjCfg, g)
}
func (g intBuilderGenerator) EmitNodeAssemblerOtherBits(w io.Writer) {
	// Nothing needed here for int kinds.
}
