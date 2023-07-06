package parser

import (
	"unicode/utf8"
)

type Document struct {
	Filename   string
	Consts     []*Const
	Typedefs   []*Typedef
	Enums      []*Enum
	Services   []*Service
	Structs    []*Struct
	Unions     []*Union
	Exceptions []*Exception
}

func NewDocument(defs []Definition) *Document {
	doc := &Document{}

	for _, def := range defs {
		switch def.Type() {
		case "Const":
			doc.Consts = append(doc.Consts, def.(*Const))
		case "Typedef":
			doc.Typedefs = append(doc.Typedefs, def.(*Typedef))
		case "Enum":
			doc.Enums = append(doc.Enums, def.(*Enum))
		case "Service":
			doc.Services = append(doc.Services, def.(*Service))
		case "Struct":
			doc.Structs = append(doc.Structs, def.(*Struct))
		case "Union":
			doc.Unions = append(doc.Unions, def.(*Union))
		case "Exception":
			doc.Exceptions = append(doc.Exceptions, def.(*Exception))
		}
	}
	return doc
}

type Definition interface {
	Type() string
	SetComments(comments string)
}

type Struct struct {
	Identifier *Identifier
	Fields     []*Field
	Comments   string
}

func NewStruct(identifier *Identifier, fields []*Field) *Struct {
	return &Struct{
		Identifier: identifier,
		Fields:     fields,
	}
}

func (s *Struct) Type() string {
	return "Struct"
}

func (s *Struct) SetComments(comments string) {
	s.Comments = comments
}

type Const struct {
	Name      *Identifier
	ConstType *FieldType
	Value     *ConstValue
	Comments  string
}

func NewConst(name *Identifier, t *FieldType, v *ConstValue, comments string) *Const {
	return &Const{
		Name:      name,
		ConstType: t,
		Value:     v,
		Comments:  comments,
	}
}

func (c *Const) Type() string {
	return "Const"
}

func (c *Const) SetComments(comments string) {
	c.Comments = comments
}

type Typedef struct {
	T        *FieldType
	Alias    *Identifier
	Comments string
}

func NewTypedef(t *FieldType, alias *Identifier) *Typedef {
	return &Typedef{
		T:     t,
		Alias: alias,
	}
}

func (t *Typedef) Type() string {
	return "Typedef"
}

func (t *Typedef) SetComments(comments string) {
	t.Comments = comments
}

type Enum struct {
	Name     *Identifier
	Values   []*EnumValue
	Comments string
}

func NewEnum(name *Identifier, values []*EnumValue) *Enum {
	return &Enum{
		Name:   name,
		Values: values,
	}
}

func (e *Enum) Type() string {
	return "Enum"
}

func (e *Enum) SetComments(comments string) {
	e.Comments = comments
}

type EnumValue struct {
	Name     *Identifier
	Value    int64
	Comments string
}

func NewEnumValue(name *Identifier, value int64, comments string) *EnumValue {
	return &EnumValue{
		Name:     name,
		Value:    value,
		Comments: comments,
	}
}

type Service struct {
	Name      *Identifier
	Extends   *Identifier
	Functions []*Function
	Comments  string
}

func NewService(name *Identifier, extends *Identifier, fns []*Function) *Service {
	return &Service{
		Name:      name,
		Extends:   extends,
		Functions: fns,
	}
}

func (s *Service) Type() string {
	return "Service"
}

func (s *Service) SetComments(comments string) {
	s.Comments = comments
}

type Function struct {
	Name         *Identifier
	Oneway       bool
	Void         bool
	FunctionType *FieldType
	Arguments    []*Field
	Throws       []*Field
	Comments     string
}

func NewFunction(name *Identifier, oneway bool, void bool, ft *FieldType, args []*Field, throws []*Field, comments string) *Function {
	return &Function{
		Name:         name,
		Oneway:       oneway,
		Void:         void,
		FunctionType: ft,
		Arguments:    args,
		Throws:       throws,
		Comments:     comments,
	}
}

type Union struct {
	Name     *Identifier
	Fields   []*Field
	Comments string
}

func NewUnion(name *Identifier, fields []*Field) *Union {
	return &Union{
		Name:   name,
		Fields: fields,
	}
}

func (u *Union) Type() string {
	return "Union"
}

func (u *Union) SetComments(comments string) {
	u.Comments = comments
}

type Exception struct {
	Name     *Identifier
	Fields   []*Field
	Comments string
}

func NewException(name *Identifier, fields []*Field) *Exception {
	return &Exception{
		Name:   name,
		Fields: fields,
	}
}

func (e *Exception) Type() string {
	return "Exception"
}

func (e *Exception) SetComments(comments string) {
	e.Comments = comments
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
	KeyType *FieldType
	// only exist when TypeName is map
	ValueType *FieldType
	BadNode   bool
}

func NewFieldType(typeName *TypeName, keyType *FieldType, valueType *FieldType) *FieldType {
	return &FieldType{
		TypeName:  typeName,
		KeyType:   keyType,
		ValueType: valueType,
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
