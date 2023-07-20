package parser

import (
	"path"
	"strings"
	"unicode/utf8"
)

type Node interface {
	// position of first charactor of this node
	Pos() Position
	// position of first charactor immediately after this node
	End() Position

	Contains(pos Position) bool

	Children() []Node

	Type() string
}

type Document struct {
	Filename string

	BadHeaders  []*BadHeader
	Includes    []*Include
	CPPIncludes []*CPPInclude
	Namespaces  []*Namespace

	Consts         []*Const
	Typedefs       []*Typedef
	Enums          []*Enum
	Services       []*Service
	Structs        []*Struct
	Unions         []*Union
	Exceptions     []*Exception
	BadDefinitions []*BadDefinition

	Nodes []Node

	Location
}

func NewDocument(headers []Header, defs []Definition, loc Location) *Document {
	doc := &Document{
		Location: loc,
	}

	for _, header := range headers {
		switch header.Type() {
		case "Include":
			doc.Includes = append(doc.Includes, header.(*Include))
		case "CPPInclude":
			doc.CPPIncludes = append(doc.CPPIncludes, header.(*CPPInclude))
		case "Namespace":
			doc.Namespaces = append(doc.Namespaces, header.(*Namespace))
		case "BadHeader":
			doc.BadHeaders = append(doc.BadHeaders, header.(*BadHeader))
		}
		doc.Nodes = append(doc.Nodes, header)
	}

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
		case "BadDefinition":
			doc.BadDefinitions = append(doc.BadDefinitions, def.(*BadDefinition))
		}
		doc.Nodes = append(doc.Nodes, def)
	}
	return doc
}

func (d *Document) Children() []Node {
	return d.Nodes
}

func (d *Document) Type() string {
	return "Document"
}

type Header interface {
	Type() string
	Node
}

type BadHeader struct {
	BadNode bool
	Location
}

func NewBadHeader(loc Location) *BadHeader {
	return &BadHeader{
		BadNode:  true,
		Location: loc,
	}
}

func (h *BadHeader) Type() string {
	return "BadHeader"
}

func (h *BadHeader) Children() []Node {
	return nil
}

type Include struct {
	Path *Literal

	BadNode bool
	Location
}

func NewInclude(path *Literal, loc Location) *Include {
	return &Include{
		Location: loc,
		Path:     path,
	}
}

func NewBadInclude(loc Location) *Include {
	return &Include{
		BadNode:  true,
		Location: loc,
	}
}

func (i *Include) Type() string {
	return "Include"
}

func (i *Include) Name() string {
	_, file := path.Split(i.Path.Value)
	name := strings.TrimRight(file, path.Ext(file))
	return name
}

func (i *Include) Children() []Node {
	return nil
}

type CPPInclude struct {
	Path *Literal

	BadNode bool
	Location
}

func NewCPPInclude(path *Literal, loc Location) *CPPInclude {
	return &CPPInclude{
		Location: loc,
		Path:     path,
	}
}

func NewBadCPPInclude(loc Location) *CPPInclude {
	return &CPPInclude{
		BadNode:  true,
		Location: loc,
	}
}

func (i *CPPInclude) Type() string {
	return "CPPInclude"
}

func (i *CPPInclude) Children() []Node {
	return nil
}

type Namespace struct {
	Language string
	Name     string

	BadNode bool
	Location
}

func NewNamespace(language, name string, loc Location) *Namespace {
	return &Namespace{
		Language: language,
		Name:     name,
		Location: loc,
	}
}

func NewBadNamespace(loc Location) *Namespace {
	return &Namespace{
		BadNode:  true,
		Location: loc,
	}
}

func (n *Namespace) Type() string {
	return "Namespace"
}

func (n *Namespace) Children() []Node {
	return nil
}

type Definition interface {
	Node
	Type() string
	SetComments(comments string)
}

type BadDefinition struct {
	BadNode bool
	Location
}

func NewBadDefinition(loc Location) *BadDefinition {
	return &BadDefinition{
		BadNode:  true,
		Location: loc,
	}
}

func (d *BadDefinition) Type() string {
	return "Definition"
}

func (d *BadDefinition) Children() []Node {
	return nil
}

func (d *BadDefinition) SetComments(string) {
}

type Struct struct {
	Identifier *Identifier
	Fields     []*Field
	Comments   string

	BadNode bool
	Location
}

func NewStruct(identifier *Identifier, fields []*Field, loc Location) *Struct {
	return &Struct{
		Identifier: identifier,
		Fields:     fields,
		Location:   loc,
	}
}

func NewBadStruct(loc Location) *Struct {
	return &Struct{
		BadNode:  true,
		Location: loc,
	}
}

func (s *Struct) Type() string {
	return "Struct"
}

func (s *Struct) SetComments(comments string) {
	s.Comments = comments
}

func (s *Struct) Children() []Node {
	nodes := []Node{s.Identifier}
	for i := range s.Fields {
		nodes = append(nodes, s.Fields[i])
	}

	return nodes
}

type Const struct {
	Name      *Identifier
	ConstType *FieldType
	Value     *ConstValue
	Comments  string

	BadNode bool
	Location
}

func NewConst(name *Identifier, t *FieldType, v *ConstValue, comments string, loc Location) *Const {
	return &Const{
		Name:      name,
		ConstType: t,
		Value:     v,
		Comments:  comments,
		Location:  loc,
	}
}

func NewBadConst(loc Location) *Const {
	return &Const{
		BadNode:  true,
		Location: loc,
	}
}

func (c *Const) Type() string {
	return "Const"
}

func (c *Const) SetComments(comments string) {
	c.Comments = comments
}

func (c *Const) Children() []Node {
	return []Node{c.Name, c.ConstType, c.Value}
}

type Typedef struct {
	T        *FieldType
	Alias    *Identifier
	Comments string
	BadNode  bool

	Location
}

func NewTypedef(t *FieldType, alias *Identifier, loc Location) *Typedef {
	return &Typedef{
		T:        t,
		Alias:    alias,
		Location: loc,
	}
}

func NewBadTypedef(loc Location) *Typedef {
	return &Typedef{
		BadNode:  true,
		Location: loc,
	}
}

func (t *Typedef) Type() string {
	return "Typedef"
}

func (t *Typedef) SetComments(comments string) {
	t.Comments = comments
}

func (t *Typedef) Children() []Node {
	return []Node{t.T, t.Alias}
}

type Enum struct {
	Name     *Identifier
	Values   []*EnumValue
	Comments string

	BadNode bool
	Location
}

func NewEnum(name *Identifier, values []*EnumValue, loc Location) *Enum {
	return &Enum{
		Name:     name,
		Values:   values,
		Location: loc,
	}
}

func NewBadEnum(loc Location) *Enum {
	return &Enum{
		BadNode:  true,
		Location: loc,
	}
}

func (e *Enum) Type() string {
	return "Enum"
}

func (e *Enum) SetComments(comments string) {
	e.Comments = comments
}

func (e *Enum) Children() []Node {
	nodes := []Node{e.Name}
	for i := range e.Values {
		nodes = append(nodes, e.Values[i])
	}
	return nodes
}

type EnumValue struct {
	Name      *Identifier
	ValueNode *ConstValue
	Value     int64 // Value only record enum value. it is not a ast node
	Comments  string

	Location
}

func NewEnumValue(name *Identifier, valueNode *ConstValue, value int64, comments string, loc Location) *EnumValue {
	return &EnumValue{
		Name:      name,
		ValueNode: valueNode,
		Value:     value,
		Comments:  comments,
		Location:  loc,
	}
}

func (e *EnumValue) Children() []Node {
	nodes := []Node{e.Name}
	if e.ValueNode != nil {
		nodes = append(nodes, e.ValueNode)
	}

	return nodes
}

func (e *EnumValue) Type() string {
	return "EnumValue"
}

type Service struct {
	Name      *Identifier
	Extends   *Identifier
	Functions []*Function
	Comments  string

	BadNode bool
	Location
}

func NewService(name *Identifier, extends *Identifier, fns []*Function, loc Location) *Service {
	return &Service{
		Name:      name,
		Extends:   extends,
		Functions: fns,
		Location:  loc,
	}
}

func NewBadService(loc Location) *Service {
	return &Service{
		BadNode:  true,
		Location: loc,
	}
}

func (s *Service) Type() string {
	return "Service"
}

func (s *Service) SetComments(comments string) {
	s.Comments = comments
}

func (s *Service) Children() []Node {
	nodes := []Node{s.Name, s.Extends}
	for i := range s.Functions {
		nodes = append(nodes, s.Functions[i])
	}

	return nodes
}

type Function struct {
	Name         *Identifier
	Oneway       bool
	Void         bool
	FunctionType *FieldType
	Arguments    []*Field
	Throws       []*Field
	Comments     string

	BadNode bool
	Location
}

func NewFunction(name *Identifier, oneway bool, void bool, ft *FieldType, args []*Field, throws []*Field, comments string, loc Location) *Function {
	return &Function{
		Name:         name,
		Oneway:       oneway,
		Void:         void,
		FunctionType: ft,
		Arguments:    args,
		Throws:       throws,
		Comments:     comments,
		Location:     loc,
	}
}

func NewBadFunc(loc Location) *Function {
	return &Function{
		BadNode:  true,
		Location: loc,
	}
}

func (f *Function) Children() []Node {
	nodes := []Node{f.Name, f.FunctionType}
	for i := range f.Arguments {
		nodes = append(nodes, f.Arguments[i])
	}
	for i := range f.Throws {
		nodes = append(nodes, f.Throws[i])
	}

	return nodes
}

func (f *Function) Type() string {
	return "Function"
}

type Union struct {
	Name     *Identifier
	Fields   []*Field
	Comments string

	BadNode bool
	Location
}

func NewUnion(name *Identifier, fields []*Field, loc Location) *Union {
	return &Union{
		Name:     name,
		Fields:   fields,
		Location: loc,
	}
}

func NewBadUnion(loc Location) *Union {
	return &Union{
		BadNode:  true,
		Location: loc,
	}
}

func (u *Union) Type() string {
	return "Union"
}

func (u *Union) SetComments(comments string) {
	u.Comments = comments
}

func (u *Union) Children() []Node {
	nodes := []Node{u.Name}
	for i := range u.Fields {
		nodes = append(nodes, u.Fields[i])
	}
	return nodes
}

type Exception struct {
	Name     *Identifier
	Fields   []*Field
	Comments string

	BadNode bool
	Location
}

func NewException(name *Identifier, fields []*Field, loc Location) *Exception {
	return &Exception{
		Name:     name,
		Fields:   fields,
		Location: loc,
	}
}

func NewBadException(loc Location) *Exception {
	return &Exception{
		BadNode:  true,
		Location: loc,
	}
}

func (e *Exception) Type() string {
	return "Exception"
}

func (e *Exception) SetComments(comments string) {
	e.Comments = comments
}

func (e *Exception) Children() []Node {
	nodes := []Node{e.Name}
	for i := range e.Fields {
		nodes = append(nodes, e.Fields[i])
	}
	return nodes
}

type Identifier struct {
	Name    string
	BadNode bool
	Location
}

func NewIdentifier(name string, pos position) *Identifier {
	id := &Identifier{
		Name:     name,
		Location: NewLocation(pos, name),
		BadNode:  name == "",
	}

	return id
}

func (i *Identifier) ToFieldType() *FieldType {
	t := &FieldType{
		TypeName: &TypeName{
			Name:     i.Name,
			Location: i.Location,
		},
		Location: i.Location,
	}

	return t
}

func (i *Identifier) Children() []Node {
	return nil
}

func (i *Identifier) Type() string {
	return "Identifier"
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
	Index        *FieldIndex
	Required     *Required
	FieldType    *FieldType
	Identifier   *Identifier
	ConstValue   *ConstValue

	BadNode bool
	Location
}

func NewField(comments string, lineComments string, index *FieldIndex, required *Required, fieldType *FieldType, identifier *Identifier, constValue *ConstValue, loc Location) *Field {
	field := &Field{
		Comments:     comments,
		LineComments: lineComments,
		Index:        index,
		Required:     required,
		FieldType:    fieldType,
		Identifier:   identifier,
		ConstValue:   constValue,
		BadNode:      fieldType == nil,
		Location:     loc,
	}
	return field
}

func (f *Field) Children() []Node {
	return []Node{f.Required, f.FieldType, f.Identifier, f.ConstValue}
}

func (f *Field) Type() string {
	return "Field"
}

type FieldIndex struct {
	Value int

	BadNode bool
	Location
}

func NewFieldIndex(v int, loc Location) *FieldIndex {
	return &FieldIndex{
		Value:    v,
		Location: loc,
	}
}

func NewBadFieldIndex(loc Location) *FieldIndex {
	return &FieldIndex{
		BadNode:  true,
		Location: loc,
	}
}

func (f *FieldIndex) Children() []Node {
	return nil
}

func (f *FieldIndex) Type() string {
	return "FieldIndex"
}

type Required struct {
	Required bool
	BadNode  bool
	Location
}

func NewRequired(required bool, pos position) *Required {
	req := &Required{
		Required: required,
	}
	if required {
		req.Location = NewLocation(pos, "required")
	} else {
		req.Location = NewLocation(pos, "optional")
	}

	return req
}

func NewBadRequired(text string, pos position) *Required {
	req := &Required{
		Location: NewLocation(pos, text),
		BadNode:  true,
	}

	return req
}

func (r *Required) Children() []Node {
	return nil
}

func (r *Required) Type() string {
	return "Required"
}

type FieldType struct {
	TypeName *TypeName
	// only exist when TypeName is map or set or list
	KeyType *FieldType
	// only exist when TypeName is map
	ValueType *FieldType
	BadNode   bool

	Location
}

func NewFieldType(typeName *TypeName, keyType *FieldType, valueType *FieldType, loc Location) *FieldType {
	return &FieldType{
		TypeName:  typeName,
		KeyType:   keyType,
		ValueType: valueType,
		Location:  loc,
	}
}

func (c *FieldType) Children() []Node {
	nodes := make([]Node, 0, 1)
	nodes = append(nodes, c.TypeName)
	if c.KeyType != nil {
		nodes = append(nodes, c.KeyType)
	}
	if c.ValueType != nil {
		nodes = append(nodes, c.ValueType)
	}

	return nodes
}

func (c *FieldType) Type() string {
	return "FieldType"
}

type TypeName struct {
	// TypeName can be:
	// container type: map, set, list
	// base type: bool, byte, i8, i16, i32, i64, double, string, binary
	// struct, enum, union, exception
	Name string
	Location
}

func NewTypeName(name string, pos position) *TypeName {
	t := &TypeName{
		Name:     name,
		Location: NewLocation(pos, name),
	}

	return t
}

func (t *TypeName) Children() []Node {
	return nil
}

func (t *TypeName) Type() string {
	return "TypeName"
}

type ConstValue struct {
	TypeName string
	Value    any

	// only exist when TypeName is map
	Key any

	Location
}

func NewConstValue(typeName string, value any, loc Location) *ConstValue {
	return &ConstValue{
		TypeName: typeName,
		Value:    value,
		Location: loc,
	}
}

func NewMapConstValue(key, value *ConstValue, loc Location) *ConstValue {
	return &ConstValue{
		TypeName: "map",
		Key:      key,
		Value:    value,
		Location: loc,
	}
}

// TODO(jpf): nodes of key, value
func (c *ConstValue) Children() []Node {
	return nil
}

func (c *ConstValue) Type() string {
	return "ConstValue"
}

type Literal struct {
	Value   string
	BadNode bool

	Location
}

func NewLiteral(v string, loc Location) *Literal {
	return &Literal{
		Value:    v,
		Location: loc,
	}
}

func NewBadLiteral(v string, loc Location) *Literal {
	return &Literal{
		Location: loc,
		BadNode:  true,
	}
}

type Location struct {
	StartPos Position
	EndPos   Position
}

func (l Location) MoveStartInLine(n int) Location {
	newL := l
	newL.StartPos.Col += n
	newL.StartPos.Offset += n

	return newL
}

func (l *Location) Pos() Position {
	return l.StartPos
}

// end col and offset is excluded
func (l *Location) End() Position {
	return l.EndPos
}

func (l *Location) Contains(pos Position) bool {
	if l == nil {
		return false
	}
	return (l.StartPos.Less(pos) || l.StartPos.Equal(pos)) && l.EndPos.Greater(pos)
}

func NewLocationFromPos(start, end Position) Location {
	return Location{StartPos: start, EndPos: end}
}

func NewLocationFromCurrent(c *current) Location {
	return NewLocation(c.pos, string(c.text))
}

func NewLocation(startPos position, text string) Location {
	start := ConvertPosition(startPos)

	nLine := strings.Count(text, "\n")
	lastLineOffset := strings.LastIndexByte(text, '\n')
	if lastLineOffset == -1 {
		lastLineOffset = 0
	}
	lastLine := []byte(text)[lastLineOffset:]
	col := utf8.RuneCount(lastLine) + 1
	if nLine == 0 {
		col += start.Col - 1
	}
	end := Position{
		Line:   start.Line + nLine,
		Col:    col,
		Offset: start.Offset + len(text),
	}

	return Location{
		StartPos: start,
		EndPos:   end,
	}
}

var InvalidPosition = Position{
	Line:   -1,
	Col:    -1,
	Offset: -1,
}

type Position struct {
	Line   int // 1-based line number
	Col    int // 1-based rune count from start of line.
	Offset int // 0-based byte offset
}

func (p *Position) Less(other Position) bool {
	if p.Line < other.Line {
		return true
	} else if p.Line == other.Line {
		return p.Col < other.Col
	}
	return false
}

func (p *Position) Equal(other Position) bool {
	return p.Line == other.Line && p.Col == other.Col
}

func (p *Position) Greater(other Position) bool {
	if p.Line > other.Line {
		return true
	} else if p.Line == other.Line {
		return p.Col > other.Col
	}
	return false
}

func (p *Position) Invalid() bool {
	return p.Line < 1 || p.Col < 1 || p.Offset < 0
}
