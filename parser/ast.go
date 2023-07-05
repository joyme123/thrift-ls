package parser

import (
	"unicode/utf8"
)

type Document struct {
	Filename string
	Structs  []*Struct
}

func NewDocument(structs []*Struct) *Document {
	doc := &Document{
		Structs: structs,
	}
	return doc
}

type Struct struct {
	Identifier *Identifier
	Fields     []*Field
}

func NewStruct(identifier *Identifier, fields []*Field) *Struct {
	return &Struct{
		Identifier: identifier,
		Fields:     fields,
	}
}

type Identifier struct {
	Name    string
	Loc     Location
	BadNode bool
}

func NewIdentifier(name string, pos position) *Identifier {
	id := &Identifier{
		Name:    name,
		Loc:     NewLocation(pos, name),
		BadNode: name == "",
	}

	return id
}

func (i *Identifier) ToFieldType() *FieldType {
	t := &FieldType{
		TypeName: &TypeName{
			Name: i.Name,
			Loc:  i.Loc,
		},
	}

	return t
}

func ConvertPosition(pos position) Position {
	return Position{
		Line:   pos.line,
		Col:    pos.col,
		Offset: pos.offset,
	}
}

type Field struct {
	Comments     string
	LineComments string
	Index        int
	Required     *Required
	FieldType    *FieldType
	Identifier   *Identifier
	ConstValue   *ConstValue

	BadNode bool
}

func NewField(comments string, lineComments string, index int, required *Required, fieldType *FieldType, identifier *Identifier, constValue *ConstValue) *Field {
	field := &Field{
		Comments:     comments,
		LineComments: lineComments,
		Index:        index,
		Required:     required,
		FieldType:    fieldType,
		Identifier:   identifier,
		ConstValue:   constValue,
		BadNode:      fieldType == nil,
	}
	return field
}

type Required struct {
	Required bool
	Loc      Location
	BadNode  bool
}

func NewRequired(required bool, pos position) *Required {
	req := &Required{
		Required: required,
	}
	if required {
		req.Loc = NewLocation(pos, "required")
	} else {
		req.Loc = NewLocation(pos, "optional")
	}

	return req
}

func NewBadRequired(text string, pos position) *Required {
	req := &Required{
		Loc:     NewLocation(pos, text),
		BadNode: true,
	}

	return req
}

type FieldType struct {
	TypeName *TypeName
	// only exist when TypeName is map or set or list
	Key *FieldType
	// only exist when TypeName is map
	Value   *FieldType
	BadNode bool
}

func NewFieldType(typeName *TypeName, key *FieldType, value *FieldType) *FieldType {
	return &FieldType{
		TypeName: typeName,
		Key:      key,
		Value:    value,
	}
}

type TypeName struct {
	// TypeName can be:
	// container type: map, set, list
	// base type: bool, byte, i8, i16, i32, i64, double, string, binary
	// struct
	Name string
	Loc  Location
}

func NewTypeName(name string, pos position) *TypeName {
	t := &TypeName{
		Name: name,
		Loc:  NewLocation(pos, name),
	}

	return t
}

type ConstValue struct {
	TypeName string
	Value    any

	// only exist when TypeName is map
	Key any
}

func NewConstValue(typeName string, value any) *ConstValue {
	return &ConstValue{
		TypeName: typeName,
		Value:    value,
	}
}

func NewMapConstValue(key, value *ConstValue) *ConstValue {
	return &ConstValue{
		TypeName: "map",
		Key:      key,
		Value:    value,
	}
}

type Location struct {
	Start Position
	End   Position
}

func NewLocation(startPos position, token string) Location {
	start := ConvertPosition(startPos)
	end := start
	end.Col = start.Col + utf8.RuneCountInString(token)
	end.Offset = start.Offset + utf8.RuneCountInString(token)

	return Location{
		Start: start,
		End:   end,
	}
}

type Position struct {
	Line   int
	Col    int
	Offset int
}
