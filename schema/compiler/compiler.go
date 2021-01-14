/*
	The schema/compiler package contains concrete implementations of the
	interfaces in the schema package which are used to describe IPLD Schemas,
	and it also provides a Compiler type which is used to construct them.
*/
package compiler

import (
	"fmt"

	"github.com/ipld/go-ipld-prime/schema"
)

// Compiler creates new TypeSystem instances.
// Methods are called on a Compiler instance to add types to the set,
// and when done, the Compile method is called, which can return
// either a list of error values or a new TypeSystem.
//
// Users don't usually use Compiler themselves,
// and this API isn't meant to be especially user-friendly.
// It's better to write IPLD Schemas using the DSL,
// parse and transpile that into the standard DMT format,
// and then read that with `schema/dmt` package and use the `dmt.Compile` feature.
// This lets you spend more time with the human-readable syntax and DMT format,
// which in addition to being better suited for documentation and review,
// is also usable with other IPLD tools and IPLD implementations in other languages.
// (Inside, the `dmt.Compile` feature uses this Compiler for you.)
//
// On error handling:
// Since several sorts of error can't be checked until the whole group of types has been stated
// (for example, referential completeness checks),
// almost none of the methods on Compiler return errors as they go.
// All errors will just be reported altogether at once at the end, when Compile is called.
// Some extremely obvious errors, like trying to use the same TypeName twice, will cause a panic immediately.
// The rule for errors that are raised as panics is that they must have been already avoided if the data were coming from the schemadmt package.
// (E.g., if something could be invalidly sent to the Compiler twice, but was a map key in the schemadmt and so already checked as unique, that's panic-worthy here.
// But if repeats of some identifier are invalid but would a list when expressed in the schemadmt, that's *not* allowed to panic here.)
//
// On immutability:
// The TypeSystem returned by a successful Compile call will be immutable.
// Many methods on the Compiler type are structured to accept data in a way that works towards this immutability.
// In particular, many methods on Compiler take arguments which are "carrier types" for segments of immutable data,
// which must be produced by constructor functions; for one example of this pattern, see the interplay of Compiler.TypeStruct() and MakeStructFieldList().
type Compiler struct {
	// ... and if you're wondering why this type is exported at all?
	//  Well, arguably, it's useful to be able to construct these values without going through the dmt.
	//  At the end of the day, though?  Honestly, import cycle breaking.  This was not the first choice.
	// An implementation which wraps the schemadmt package to make it fit the schema interfaces was the first choice
	//  because it would've saved a *lot* of work (it would've removed the need for this compiler system entirely, among other things);
	//  but that doesn't fly, because the dmt types have to implement schema.Type, and that interface refers to yet more schema.* types.
	//  And that would make an import cycle if we tried to put types wrapping the dmt types into the schema package.  Whoops.
	//  So, here we are.
	//
	// The decision to split out this Compiler type and all its other related Make* functions
	//  from the schema package is largely cosmetic; it technically could've been placed in the schema package.
	//  However, the result of splitting the readable and writable types into packages seemed more readable,
	//   and gives us more elbow room in the godocs to suggest "you probably shouldn't use these directly".
	//  Compared to the already-forced dmt package split, having the creation stuff in one package
	//   and the read-only interfaces in another package just isn't much additional burden.

	// ts gathers all the in-progress types (including anonymous ones),
	// and is eventually the value we return (if Compile is ultimately successful).
	// We insert into this blindly as we go, and check everything for consistency at the end;
	// if those logical checks flunk, we don't allow any reference to it to escape.
	// This is nil'd after any Compile, so when we give a reference to it away,
	// it's immutable from there on out.
	ts *TypeSystem
}

func (c *Compiler) Init() {
	c.ts = &TypeSystem{
		map[schema.TypeReference]schema.Type{},
		nil,
	}
}

func (c *Compiler) Compile() (schema.TypeSystem, error) {
	panic("TODO")
}

func (c *Compiler) addType(t schema.Type) {
	c.mustHaveNameFree(t.Name())
	c.ts.types[schema.TypeReference(t.Name())] = t
	c.ts.list = append(c.ts.list, t)
}
func (c *Compiler) addAnonType(t schema.Type) {
	c.ts.types[schema.TypeReference(t.Name())] = t // FIXME it's... probably a bug that the schema.Type.Name() method doesn't return a TypeReference.  Yeah, it definitely is.  TypeMap and TypeList should have their own name field internally be TypeReference, too, because it's true.  wonder if we should have separate methods on the schema.Type interface for this.  would probably be a usability trap to do so, though (too many user printfs would use the Name function and get blanks and be surprised).
}

func (c *Compiler) mustHaveNameFree(name schema.TypeName) {
	if _, exists := c.ts.types[schema.TypeReference(name)]; exists {
		panic(fmt.Errorf("type name %q already used", name))
	}
}

func (c *Compiler) TypeBool(name schema.TypeName) {
	c.addType(&TypeBool{c.ts, name})
}

func (c *Compiler) TypeString(name schema.TypeName) {
	c.addType(&TypeString{c.ts, name})
}

func (c *Compiler) TypeBytes(name schema.TypeName) {
	c.addType(&TypeBytes{c.ts, name})
}

func (c *Compiler) TypeInt(name schema.TypeName) {
	c.addType(&TypeInt{c.ts, name})
}

func (c *Compiler) TypeFloat(name schema.TypeName) {
	c.addType(&TypeFloat{c.ts, name})
}

func (c *Compiler) TypeLink(name schema.TypeName, expectedTypeRef schema.TypeName) {
	c.addType(&TypeLink{c.ts, name, expectedTypeRef})
}

func (c *Compiler) TypeStruct(name schema.TypeName, fields structFieldList, rstrat StructRepresentation) {
	t := TypeStruct{
		ts:        c.ts,
		name:      name,
		fields:    fields.x, // it's safe to take this directly because the carrier type means a reference to this slice has never been exported.
		fieldsMap: make(map[StructFieldName]*StructField, len(fields.x)),
		rstrat:    rstrat,
	}
	c.addType(&t)
	for i, f := range fields.x {
		// duplicate names are rejected with a *panic* here because we expect these to already be unique (if this data is coming from the dmt, these were map keys there).
		if _, exists := t.fieldsMap[f.name]; exists {
			panic(fmt.Errorf("type %q already has field named %q", t.name, f.name))
		}
		t.fieldsMap[f.name] = &fields.x[i]
		fields.x[i].parent = &t
	}
}

// structFieldList is a carrier type that just wraps a slice reference.
// It is used so we can let code outside this package hold a value of this type without letting the slice become mutable.
type structFieldList struct {
	x []StructField
}

func MakeStructFieldList(fields ...StructField) structFieldList {
	return structFieldList{fields}
}
func MakeStructField(name StructFieldName, typ schema.TypeReference, optional, nullable bool) StructField {
	return StructField{nil, name, typ, optional, nullable}
}

func MakeStructRepresentation_Map(fieldDetails ...StructRepresentation_Map_FieldDetailsEntry) StructRepresentation {
	rstrat := StructRepresentation_Map{nil, make(map[StructFieldName]StructRepresentation_Map_FieldDetails, len(fieldDetails))}
	for _, fd := range fieldDetails {
		if _, exists := rstrat.fieldDetails[fd.FieldName]; exists {
			panic(fmt.Errorf("field name %q duplicated", fd.FieldName))
		}
		rstrat.fieldDetails[fd.FieldName] = fd.Details
	}
	return rstrat
}

// StructRepresentation_Map_FieldDetailsEntry is a carrier type that associates a field name
// with field detail information that's appropriate to a map representation strategy for a struct.
// It is used to feed data to MakeStructRepresentation_Map in so that that method can build a map
// without exposing a reference to it in a way that would make that map mutable.
type StructRepresentation_Map_FieldDetailsEntry struct {
	FieldName StructFieldName
	Details   StructRepresentation_Map_FieldDetails
}

func (c *Compiler) TypeMap(name schema.TypeName, keyTypeRef schema.TypeName, valueTypeRef schema.TypeReference, valueNullable bool) {
	c.addType(&TypeMap{c.ts, name, keyTypeRef, valueTypeRef, valueNullable})
}

func (c *Compiler) TypeList(name schema.TypeName, valueTypeRef schema.TypeReference, valueNullable bool) {
	c.addType(&TypeList{c.ts, name, valueTypeRef, valueNullable})
}

func (c *Compiler) TypeUnion(name schema.TypeName, members unionMemberList, rstrat UnionRepresentation) {
	t := TypeUnion{
		ts:      c.ts,
		name:    name,
		members: members.x, // it's safe to take this directly because the carrier type means a reference to this slice has never been exported.
		rstrat:  rstrat,
	}
	c.addType(&t)
	// note! duplicate member names *not* rejected at this moment -- that's a job for the validation phase.
	//  this is an interesting contrast to how when buildings struct, dupe field names may be rejected proactively:
	//   the difference is, member names were a list in the dmt form too, so it's important we format a nice error rather than panic if there was invalid data there.
}

// unionMemberList is a carrier type that just wraps a slice reference.
// It is used so we can let code outside this package hold a value of this type without letting the slice become mutable.
type unionMemberList struct {
	x []schema.TypeName
}

func MakeUnionMemberList(members ...schema.TypeName) unionMemberList {
	return unionMemberList{members}
}

func MakeUnionRepresentation_Keyed(discriminantTable unionDiscriminantStringTable) UnionRepresentation {
	return &UnionRepresentation_Keyed{nil, discriminantTable.x}
}

// unionMemberList is a carrier type that just wraps a map reference.
// It is used so we can let code outside this package hold a value of this type without letting the map become mutable.
type unionDiscriminantStringTable struct {
	x map[string]schema.TypeName
}

func MakeUnionDiscriminantStringTable(entries ...UnionDiscriminantStringEntry) unionDiscriminantStringTable {
	x := make(map[string]schema.TypeName, len(entries))
	for _, y := range entries {
		if _, exists := x[y.Discriminant]; exists {
			panic(fmt.Errorf("discriminant string %q duplicated", y.Discriminant))
		}
		x[y.Discriminant] = y.Member
	}
	return unionDiscriminantStringTable{x}
}

// UnionRepresentation_DiscriminantStringEntry is a carrier type that associates a string with a TypeName.
// It is used to feed data to several of the union representation constructors so that those functions
// can build their results without exposing a reference to a map in a way that would make that map mutable.
type UnionDiscriminantStringEntry struct {
	Discriminant string
	Member       schema.TypeName
}
