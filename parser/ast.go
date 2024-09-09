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

	IsBadNode() bool
	ChildrenBadNode() bool

	// Equals checks ast equals between two node
	Equals(node Node) bool
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

	Comments []*Comment // Comments at end of doc

	Nodes []Node

	Location
}

func NewDocument(headers []Header, defs []Definition, comments []*Comment, loc Location) *Document {
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
	doc.Comments = comments
	for _, comment := range comments {
		doc.Nodes = append(doc.Nodes, comment)
	}
	return doc
}

func (d *Document) Children() []Node {
	return d.Nodes
}

func (d *Document) Type() string {
	return "Document"
}

func (d *Document) IsBadNode() bool {
	return false
}

func (d *Document) ChildrenBadNode() bool {
	children := d.Children()
	for i := range children {
		if children[i].IsBadNode() {
			return true
		}
		if children[i].ChildrenBadNode() {
			return true
		}
	}
	return false
}

func (d *Document) Equals(node Node) bool {
	doc, ok := node.(*Document)
	if !ok {
		return false
	}

	if (d == nil && doc != nil) ||
		(d != nil && doc == nil) {
		return false
	} else if d == nil && doc == nil {
		return true
	}

	if len(doc.Nodes) != len(d.Nodes) {
		return false
	}

	for i := range d.Nodes {
		if !d.Nodes[i].Equals(doc.Nodes[i]) {
			return false
		}
	}

	return true
}

type Header interface {
	Type() string
	SetComments(comments []*Comment, endLineComments []*Comment)
	SetLocation(loc Location)
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

func (h *BadHeader) IsBadNode() bool {
	return true
}

func (h *BadHeader) ChildrenBadNode() bool {
	return false
}

func (h *BadHeader) SetComments([]*Comment, []*Comment) {

}

func (h *BadHeader) Equals(node Node) bool {
	hn, ok := node.(*BadHeader)
	if !ok {
		return false
	}

	if (h != nil && hn == nil) ||
		(h == nil && hn != nil) {
		return false
	}

	return true
}

func (h *BadHeader) SetLocation(loc Location) {
	h.Location = loc
}

type KeywordLiteral struct {
	Text    string
	BadNode bool
	Location
}

func NewKeywordLiteral(c *current) *KeywordLiteral {
	return &KeywordLiteral{
		Text:     string(c.text),
		Location: NewLocationFromCurrent(c),
	}
}

func NewBadKeywordLiteral(c *current) *KeywordLiteral {
	return &KeywordLiteral{
		Text:     string(c.text),
		BadNode:  true,
		Location: NewLocationFromCurrent(c),
	}
}

func (k *KeywordLiteral) Type() string {
	return "KeywordLiteral"
}

func (k *KeywordLiteral) IsBadNode() bool {
	return k.BadNode
}

func (k *KeywordLiteral) Children() []Node {
	return nil
}

func (k *KeywordLiteral) ChildrenBadNode() bool {
	return false
}

func (k *KeywordLiteral) Equals(node Node) bool {
	kl, ok := node.(*KeywordLiteral)
	if !ok {
		return false
	}

	if (k == nil && kl != nil) ||
		(k != nil && kl == nil) {
		return false
	} else if k == nil && kl == nil {
		return true
	}

	if k.BadNode != kl.BadNode {
		return false
	}

	return k.Text == kl.Text
}

type Keyword struct {
	Comments []*Comment
	Literal  *KeywordLiteral

	BadNode bool
	Location
}

func NewKeyword(comments []*Comment, literal *KeywordLiteral, loc Location) Keyword {
	return Keyword{
		Literal:  literal,
		Comments: comments,
		Location: loc,
	}
}

func (i *Keyword) Children() []Node {
	return nil
}

func (i *Keyword) IsBadNode() bool {
	return i.BadNode
}

func (i *Keyword) ChildrenBadNode() bool {
	return false
}

func (i *Keyword) Equals(k *Keyword) bool {
	if (i == nil && k != nil) ||
		(i != nil && k == nil) {
		return false
	} else if i == nil && k == nil {
		return true
	}

	if i.BadNode != k.BadNode {
		return false
	}

	if !i.Literal.Equals(k.Literal) {
		return false
	}

	if len(i.Comments) != len(k.Comments) {
		return false
	}

	for n := range i.Comments {
		if !i.Comments[n].Equals(k.Comments[n]) {
			return false
		}
	}

	return true
}

type IncludeKeyword struct {
	Keyword
}

func (i *IncludeKeyword) Type() string {
	return "IncludeKeyword"
}

func (i *IncludeKeyword) Equals(node Node) bool {
	ik, ok := node.(*IncludeKeyword)
	if !ok {
		return false
	}

	if (i == nil && ik != nil) ||
		(i != nil && ik == nil) {
		return false
	} else if i == nil && ik == nil {
		return true
	}

	return ik.Keyword.Equals(&i.Keyword)
}

type Include struct {
	IncludeKeyword *IncludeKeyword
	Path           *Literal

	Comments        []*Comment
	EndLineComments []*Comment

	BadNode bool
	Location
}

func NewInclude(keyword *IncludeKeyword, path *Literal, loc Location) *Include {
	return &Include{
		IncludeKeyword: keyword,
		Location:       loc,
		Path:           path,
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

func (i *Include) SetComments(comments []*Comment, endLineComments []*Comment) {
	i.Comments = comments
	i.EndLineComments = endLineComments
}

func (i *Include) Name() string {
	_, file := path.Split(i.Path.Value.Text)
	name := strings.TrimRight(file, path.Ext(file))
	return name
}

func (i *Include) Children() []Node {
	nodes := []Node{i.IncludeKeyword, i.Path}

	for _, com := range i.Comments {
		nodes = append(nodes, com)
	}
	for _, com := range i.EndLineComments {
		nodes = append(nodes, com)
	}

	return nodes
}

func (i *Include) IsBadNode() bool {
	return i.BadNode
}

func (i *Include) ChildrenBadNode() bool {
	children := i.Children()
	for i := range children {
		if children[i].IsBadNode() {
			return true
		}
		if children[i].ChildrenBadNode() {
			return true
		}
	}
	return false
}

func (i *Include) SetLocation(loc Location) {
	i.Location = loc
}

func (i *Include) Equals(node Node) bool {
	in, ok := node.(*Include)
	if !ok {
		return false
	}

	if (i == nil && in != nil) ||
		(i != nil && in == nil) {
		return false
	} else if i == nil && in == nil {
		return true
	}

	if i.BadNode != in.BadNode {
		return false
	}

	if !i.IncludeKeyword.Equals(in.IncludeKeyword) {
		return false
	}

	if !i.Path.Equals(in.Path) {
		return false
	}

	if len(i.Comments) != len(in.Comments) {
		return false
	}

	for n := range i.Comments {
		if !i.Comments[n].Equals(in.Comments[n]) {
			return false
		}
	}

	if len(i.EndLineComments) != len(in.EndLineComments) {
		return false
	}

	for n := range i.EndLineComments {
		if !i.EndLineComments[n].Equals(in.EndLineComments[n]) {
			return false
		}
	}

	return true
}

type CPPIncludeKeyword struct {
	Keyword
}

func (c *CPPIncludeKeyword) Type() string {
	return "CPPIncludeKeyword"
}

func (c *CPPIncludeKeyword) Equals(node Node) bool {
	cn, ok := node.(*CPPIncludeKeyword)
	if !ok {
		return false
	}

	if (c == nil && cn != nil) ||
		(c != nil && cn == nil) {
		return false
	} else if c == nil && cn == nil {
		return true
	}

	return c.Keyword.Equals(&cn.Keyword)
}

type CPPInclude struct {
	CPPIncludeKeyword *CPPIncludeKeyword
	Path              *Literal

	Comments        []*Comment
	EndLineComments []*Comment

	BadNode bool
	Location
}

func NewCPPInclude(keyword *CPPIncludeKeyword, path *Literal, loc Location) *CPPInclude {
	return &CPPInclude{
		CPPIncludeKeyword: keyword,
		Location:          loc,
		Path:              path,
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

func (i *CPPInclude) SetComments(comments []*Comment, endLineComments []*Comment) {
	i.Comments = comments
	i.EndLineComments = endLineComments
}

func (i *CPPInclude) Children() []Node {
	res := []Node{i.CPPIncludeKeyword, i.Path}
	for _, com := range i.Comments {
		res = append(res, com)
	}
	for _, com := range i.EndLineComments {
		res = append(res, com)
	}
	return res
}

func (i *CPPInclude) IsBadNode() bool {
	return i.BadNode
}

func (i *CPPInclude) ChildrenBadNode() bool {
	children := i.Children()
	for i := range children {
		if children[i].IsBadNode() {
			return true
		}
		if children[i].ChildrenBadNode() {
			return true
		}
	}

	return false
}

func (i *CPPInclude) SetLocation(loc Location) {
	i.Location = loc
}

func (i *CPPInclude) Equals(node Node) bool {
	cn, ok := node.(*CPPInclude)
	if !ok {
		return false
	}

	if (i == nil && cn != nil) ||
		(i != nil && cn == nil) {
		return false
	} else if i == nil && cn == nil {
		return true
	}

	if i.BadNode != cn.BadNode {
		return false
	}

	if !i.CPPIncludeKeyword.Equals(cn.CPPIncludeKeyword) {
		return false
	}

	if !i.Path.Equals(cn.Path) {
		return false
	}

	if len(i.Comments) != len(cn.Comments) {
		return false
	}

	for n := range i.Comments {
		if !i.Comments[n].Equals(cn.Comments[n]) {
			return false
		}
	}

	if len(i.EndLineComments) != len(cn.EndLineComments) {
		return false
	}

	for n := range i.EndLineComments {
		if !i.EndLineComments[n].Equals(cn.EndLineComments[n]) {
			return false
		}
	}

	return true
}

type NamespaceKeyword struct {
	Keyword
}

func (n *NamespaceKeyword) Type() string {
	return "NamespaceKeyword"
}

func (n *NamespaceKeyword) Equals(node Node) bool {
	nn, ok := node.(*NamespaceKeyword)
	if !ok {
		return false
	}

	if (n == nil && nn != nil) ||
		(n != nil && nn == nil) {
		return false
	} else if n == nil && nn == nil {
		return true
	}

	return n.Keyword.Equals(&nn.Keyword)
}

type NamespaceScope struct {
	Identifier
}

func (ns *NamespaceScope) Equals(node Node) bool {
	nn, ok := node.(*NamespaceScope)
	if !ok {
		return false
	}

	if (ns == nil && nn != nil) ||
		(ns != nil && nn == nil) {
		return false
	} else if ns == nil && nn == nil {
		return true
	}

	return ns.Identifier.Equals(&nn.Identifier)
}

type Namespace struct {
	NamespaceKeyword *NamespaceKeyword
	Language         *NamespaceScope
	Name             *Identifier

	Annotations     *Annotations
	Comments        []*Comment
	EndLineComments []*Comment

	BadNode bool
	Location
}

func NewNamespace(keyword *NamespaceKeyword, language *NamespaceScope, name *Identifier, annotations *Annotations, loc Location) *Namespace {
	return &Namespace{
		NamespaceKeyword: keyword,
		Language:         language,
		Name:             name,
		Annotations:      annotations,

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

func (n *Namespace) SetComments(comments []*Comment, endLineComments []*Comment) {
	n.Comments = comments
	n.EndLineComments = endLineComments
}

func (n *Namespace) Children() []Node {
	ret := []Node{n.NamespaceKeyword, n.Language, n.Name}

	for i := range n.Comments {
		ret = append(ret, n.Comments[i])
	}
	for i := range n.EndLineComments {
		ret = append(ret, n.EndLineComments[i])
	}

	if n.Annotations != nil {
		ret = append(ret, n.Annotations)
	}

	return ret
}

func (n *Namespace) IsBadNode() bool {
	return n.BadNode
}

func (n *Namespace) ChildrenBadNode() bool {
	children := n.Children()
	for i := range children {
		if children[i].IsBadNode() {
			return true
		}
		if children[i].ChildrenBadNode() {
			return true
		}
	}
	return false
}

func (n *Namespace) SetLocation(loc Location) {
	n.Location = loc
}

func (n *Namespace) Equals(node Node) bool {
	nn, ok := node.(*Namespace)
	if !ok {
		return false
	}

	if (n == nil && nn != nil) ||
		(n != nil && nn == nil) {
		return false
	} else if n == nil && nn == nil {
		return true
	}

	if n.BadNode != nn.BadNode {
		return false
	}

	if !n.NamespaceKeyword.Equals(nn.NamespaceKeyword) {
		return false
	}

	if !n.Language.Equals(nn.Language) {
		return false
	}

	if !n.Name.Equals(nn.Name) {
		return false
	}

	if !n.Annotations.Equals(nn.Annotations) {
		return false
	}

	if len(n.Comments) != len(nn.Comments) {
		return false
	}

	for i := range n.Comments {
		if !n.Comments[i].Equals(nn.Comments[i]) {
			return false
		}
	}

	if len(n.EndLineComments) != len(nn.EndLineComments) {
		return false
	}

	for i := range n.EndLineComments {
		if !n.EndLineComments[i].Equals(nn.EndLineComments[i]) {
			return false
		}
	}

	return true
}

type Definition interface {
	Node
	Type() string
	SetComments(comments []*Comment, endLineComments []*Comment)
	SetAnnotations(annotations *Annotations)
	SetLocation(loc Location)
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

func (d *BadDefinition) SetComments([]*Comment, []*Comment) {
}

func (d *BadDefinition) SetAnnotations(annos *Annotations) {

}

func (d *BadDefinition) SetLocation(loc Location) {
	d.Location = loc
}

func (d *BadDefinition) IsBadNode() bool {
	return true
}

func (d *BadDefinition) ChildrenBadNode() bool {
	return false
}

func (d *BadDefinition) Equals(node Node) bool {
	dn, ok := node.(*BadDefinition)
	if !ok {
		return false
	}

	if (d == nil && dn != nil) ||
		(d != nil && dn == nil) {
		return false
	} else if d == nil && dn == nil {
		return true
	}

	if d.BadNode != dn.BadNode {
		return false
	}

	return true
}

type StructKeyword struct {
	Keyword
}

func (s *StructKeyword) Type() string {
	return "StructKeyword"
}

func (s *StructKeyword) Equals(node Node) bool {
	sn, ok := node.(*StructKeyword)
	if !ok {
		return false
	}

	if (s == nil && sn != nil) ||
		(s != nil && sn == nil) {
		return false
	} else if s == nil && sn == nil {
		return true
	}

	return s.Keyword.Equals(&sn.Keyword)
}

type LCurKeyword struct {
	Keyword
}

func (s *LCurKeyword) Type() string {
	return "LCurKeyword"
}

func (s *LCurKeyword) Equals(node Node) bool {
	sn, ok := node.(*LCurKeyword)
	if !ok {
		return false
	}

	if (s == nil && sn != nil) ||
		(s != nil && sn == nil) {
		return false
	} else if s == nil && sn == nil {
		return true
	}

	return s.Keyword.Equals(&sn.Keyword)
}

type RCurKeyword struct {
	Keyword
}

func (s *RCurKeyword) Type() string {
	return "RCurKeyword"
}

func (s *RCurKeyword) Equals(node Node) bool {
	sn, ok := node.(*RCurKeyword)
	if !ok {
		return false
	}

	if (s == nil && sn != nil) ||
		(s != nil && sn == nil) {
		return false
	} else if s == nil && sn == nil {
		return true
	}

	return s.Keyword.Equals(&sn.Keyword)
}

type Struct struct {
	StructKeyword *StructKeyword
	LCurKeyword   *LCurKeyword
	RCurKeyword   *RCurKeyword
	Identifier    *Identifier
	Fields        []*Field

	Comments        []*Comment
	EndLineComments []*Comment
	Annotations     *Annotations

	BadNode bool
	Location
}

func NewStruct(structKeyword *StructKeyword, lCurKeyword *LCurKeyword, rCurKeyword *RCurKeyword, identifier *Identifier, fields []*Field, loc Location) *Struct {
	return &Struct{
		StructKeyword: structKeyword,
		LCurKeyword:   lCurKeyword,
		RCurKeyword:   rCurKeyword,
		Identifier:    identifier,
		Fields:        fields,
		Location:      loc,
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

func (s *Struct) SetComments(comments []*Comment, endLineComments []*Comment) {
	s.Comments = comments
	s.EndLineComments = endLineComments
}

func (s *Struct) SetAnnotations(annos *Annotations) {
	s.Annotations = annos
}

func (s *Struct) Children() []Node {
	nodes := []Node{s.StructKeyword, s.LCurKeyword, s.RCurKeyword, s.Identifier}
	for i := range s.Fields {
		nodes = append(nodes, s.Fields[i])
	}

	for i := range s.Comments {
		nodes = append(nodes, s.Comments[i])
	}
	for i := range s.EndLineComments {
		nodes = append(nodes, s.EndLineComments[i])
	}
	if s.Annotations != nil {
		nodes = append(nodes, s.Annotations)
	}

	return nodes
}

func (s *Struct) IsBadNode() bool {
	return s.BadNode
}

func (s *Struct) ChildrenBadNode() bool {
	children := s.Children()
	for i := range children {
		if children[i].IsBadNode() {
			return true
		}
		if children[i].ChildrenBadNode() {
			return true
		}
	}
	return false
}

func (s *Struct) SetLocation(loc Location) {
	s.Location = loc
}

func (s *Struct) Equals(node Node) bool {
	sn, ok := node.(*Struct)
	if !ok {
		return false
	}

	if (s == nil && sn != nil) ||
		(s != nil && sn == nil) {
		return false
	} else if s == nil && sn == nil {
		return true
	}

	if s.BadNode != sn.BadNode {
		return false
	}

	if !s.StructKeyword.Equals(sn.StructKeyword) {
		return false
	}

	if !s.LCurKeyword.Equals(sn.LCurKeyword) {
		return false
	}

	if !s.RCurKeyword.Equals(sn.RCurKeyword) {
		return false
	}

	if !s.Identifier.Equals(sn.Identifier) {
		return false
	}

	if len(s.Fields) != len(sn.Fields) {
		return false
	}

	for i := range s.Fields {
		if !s.Fields[i].Equals(sn.Fields[i]) {
			return false
		}
	}

	if len(s.Comments) != len(sn.Comments) {
		return false
	}

	for i := range s.Comments {
		if !s.Comments[i].Equals(sn.Comments[i]) {
			return false
		}
	}

	if len(s.EndLineComments) != len(sn.EndLineComments) {
		return false
	}

	for i := range s.EndLineComments {
		if !s.EndLineComments[i].Equals(sn.EndLineComments[i]) {
			return false
		}
	}

	if !s.Annotations.Equals(sn.Annotations) {
		return false
	}

	return true
}

type ConstKeyword struct {
	Keyword
}

func (c *ConstKeyword) Type() string {
	return "ConstKeyword"
}

func (c *ConstKeyword) Equals(node Node) bool {
	cn, ok := node.(*ConstKeyword)
	if !ok {
		return false
	}

	if (c == nil && cn != nil) ||
		(c != nil && cn == nil) {
		return false
	} else if c == nil && cn == nil {
		return true
	}

	return c.Keyword.Equals(&cn.Keyword)
}

type EqualKeyword struct {
	Keyword
}

func NewBadEqualKeyword() *EqualKeyword {
	return &EqualKeyword{
		Keyword: Keyword{
			BadNode: true,
		},
	}
}

func (e *EqualKeyword) Type() string {
	return "EqualKeyword"
}

func (e *EqualKeyword) Equals(node Node) bool {
	en, ok := node.(*EqualKeyword)
	if !ok {
		return false
	}

	if (e == nil && en != nil) ||
		(e != nil && en == nil) {
		return false
	} else if e == nil && en == nil {
		return true
	}

	return e.Keyword.Equals(&en.Keyword)
}

type ListSeparatorKeyword struct {
	Keyword
	Text string // , or ;
}

func (e *ListSeparatorKeyword) Type() string {
	return "ListSeparator"
}

func (e *ListSeparatorKeyword) Equals(node Node) bool {
	en, ok := node.(*ListSeparatorKeyword)
	if !ok {
		return false
	}

	if (e == nil && en != nil) ||
		(e != nil && en == nil) {
		return false
	} else if e == nil && en == nil {
		return true
	}

	return e.Keyword.Equals(&en.Keyword)
}

type Const struct {
	ConstKeyword         *ConstKeyword
	EqualKeyword         *EqualKeyword
	ListSeparatorKeyword *ListSeparatorKeyword // can be nil
	Name                 *Identifier
	ConstType            *FieldType
	Value                *ConstValue

	Comments        []*Comment
	EndLineComments []*Comment
	Annotations     *Annotations

	BadNode bool
	Location
}

func NewConst(constKeyword *ConstKeyword, equalKeyword *EqualKeyword, listSeparatorKeyword *ListSeparatorKeyword, name *Identifier, t *FieldType, v *ConstValue, loc Location) *Const {
	return &Const{
		ConstKeyword:         constKeyword,
		EqualKeyword:         equalKeyword,
		ListSeparatorKeyword: listSeparatorKeyword,
		Name:                 name,
		ConstType:            t,
		Value:                v,
		Location:             loc,
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

func (c *Const) SetComments(comments []*Comment, endLineComments []*Comment) {
	c.Comments = comments
	c.EndLineComments = endLineComments
}

func (c *Const) SetAnnotations(annos *Annotations) {
	c.Annotations = annos
}

func (c *Const) Children() []Node {
	res := []Node{c.ConstKeyword, c.EqualKeyword, c.Name, c.ConstType, c.Value}
	if c.ListSeparatorKeyword != nil {
		res = append(res, c.ListSeparatorKeyword)
	}

	for i := range c.Comments {
		res = append(res, c.Comments[i])
	}
	for i := range c.EndLineComments {
		res = append(res, c.EndLineComments[i])
	}

	if c.Annotations != nil {
		res = append(res, c.Annotations)
	}

	return res
}

func (c *Const) IsBadNode() bool {
	return c.BadNode
}

func (c *Const) ChildrenBadNode() bool {
	children := c.Children()
	for i := range children {
		if children[i].IsBadNode() {
			return true
		}
		if children[i].ChildrenBadNode() {
			return true
		}
	}
	return false
}

func (c *Const) SetLocation(loc Location) {
	c.Location = loc
}

func (c *Const) Equals(node Node) bool {
	cn, ok := node.(*Const)
	if !ok {
		return false
	}

	if (c == nil && cn != nil) ||
		(c != nil && cn == nil) {
		return false
	} else if c == nil && cn == nil {
		return true
	}

	if c.BadNode != cn.BadNode {
		return false
	}

	if !c.ConstKeyword.Equals(cn.ConstKeyword) {
		return false
	}

	if !c.EqualKeyword.Equals(cn.EqualKeyword) {
		return false
	}

	if !c.ListSeparatorKeyword.Equals(cn.ListSeparatorKeyword) {
		return false
	}

	if !c.Name.Equals(cn.Name) {
		return false
	}

	if !c.ConstType.Equals(cn.ConstType) {
		return false
	}
	if !c.Value.Equals(cn.Value) {
		return false
	}

	if len(c.Comments) != len(cn.Comments) {
		return false
	}

	for i := range c.Comments {
		if !c.Comments[i].Equals(cn.Comments[i]) {
			return false
		}
	}

	if len(c.EndLineComments) != len(cn.EndLineComments) {
		return false
	}

	for i := range c.EndLineComments {
		if !c.EndLineComments[i].Equals(cn.EndLineComments[i]) {
			return false
		}
	}

	if !c.Annotations.Equals(cn.Annotations) {
		return false
	}

	return true
}

type TypedefKeyword struct {
	Keyword
}

func (t *TypedefKeyword) Type() string {
	return "TypedefKeyword"
}

func (t *TypedefKeyword) Equals(node Node) bool {
	tn, ok := node.(*TypedefKeyword)
	if !ok {
		return false
	}

	if (t == nil && tn != nil) ||
		(t != nil && tn == nil) {
		return false
	} else if t == nil && tn == nil {
		return true
	}

	return t.Keyword.Equals(&tn.Keyword)
}

type Typedef struct {
	TypedefKeyword *TypedefKeyword
	T              *FieldType
	Alias          *Identifier

	Comments        []*Comment
	EndLineComments []*Comment
	Annotations     *Annotations
	BadNode         bool

	Location
}

func NewTypedef(keyword *TypedefKeyword, t *FieldType, alias *Identifier, loc Location) *Typedef {
	return &Typedef{
		TypedefKeyword: keyword,
		T:              t,
		Alias:          alias,
		Location:       loc,
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

func (t *Typedef) SetComments(comments []*Comment, endLineComments []*Comment) {
	t.Comments = comments
	t.EndLineComments = endLineComments
}

func (t *Typedef) SetAnnotations(annos *Annotations) {
	t.Annotations = annos
}

func (t *Typedef) Children() []Node {
	nodes := []Node{t.TypedefKeyword, t.T, t.Alias}

	for i := range t.Comments {
		nodes = append(nodes, t.Comments[i])
	}
	for i := range t.EndLineComments {
		nodes = append(nodes, t.EndLineComments[i])
	}
	if t.Annotations != nil {
		nodes = append(nodes, t.Annotations)
	}

	return nodes
}

func (t *Typedef) IsBadNode() bool {
	return t.BadNode
}

func (t *Typedef) ChildrenBadNode() bool {
	children := t.Children()
	for i := range children {
		if children[i].IsBadNode() {
			return true
		}
		if children[i].ChildrenBadNode() {
			return true
		}
	}
	return false
}

func (t *Typedef) SetLocation(loc Location) {
	t.Location = loc
}

func (t *Typedef) Equals(node Node) bool {
	tn, ok := node.(*Typedef)
	if !ok {
		return false
	}

	if (t == nil && tn != nil) ||
		(t != nil && tn == nil) {
		return false
	} else if t == nil && tn == nil {
		return true
	}

	if t.BadNode != tn.BadNode {
		return false
	}

	if !t.TypedefKeyword.Equals(tn.TypedefKeyword) {
		return false
	}

	if !t.T.Equals(tn.T) {
		return false
	}

	if !t.Alias.Equals(tn.Alias) {
		return false
	}

	if len(t.Comments) != len(tn.Comments) {
		return false
	}

	for i := range t.Comments {
		if !t.Comments[i].Equals(tn.Comments[i]) {
			return false
		}
	}

	if len(t.EndLineComments) != len(tn.EndLineComments) {
		return false
	}

	for i := range t.EndLineComments {
		if !t.EndLineComments[i].Equals(tn.EndLineComments[i]) {
			return false
		}
	}

	if !t.Annotations.Equals(tn.Annotations) {
		return false
	}

	return true
}

type EnumKeyword struct {
	Keyword
}

func (e *EnumKeyword) Type() string {
	return "EnumKeyword"
}

func (e *EnumKeyword) Equals(node Node) bool {
	en, ok := node.(*EnumKeyword)
	if !ok {
		return false
	}

	if (e == nil && en != nil) ||
		(e != nil && en == nil) {
		return false
	} else if e == nil && en == nil {
		return true
	}

	return e.Keyword.Equals(&en.Keyword)
}

type Enum struct {
	EnumKeyword *EnumKeyword
	LCurKeyword *LCurKeyword
	RCurKeyword *RCurKeyword
	Name        *Identifier
	Values      []*EnumValue

	Comments        []*Comment
	EndLineComments []*Comment
	Annotations     *Annotations

	BadNode bool
	Location
}

func NewEnum(enumKeyword *EnumKeyword, lCurKeyword *LCurKeyword, rCurKeyword *RCurKeyword, name *Identifier, values []*EnumValue, loc Location) *Enum {
	return &Enum{
		EnumKeyword: enumKeyword,
		LCurKeyword: lCurKeyword,
		RCurKeyword: rCurKeyword,
		Name:        name,
		Values:      values,
		Location:    loc,
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

func (e *Enum) SetComments(comments []*Comment, endlineComments []*Comment) {
	e.Comments = comments
	e.EndLineComments = endlineComments
}

func (e *Enum) SetAnnotations(annos *Annotations) {
	e.Annotations = annos
}

func (e *Enum) Children() []Node {
	nodes := []Node{e.Name}
	for i := range e.Values {
		nodes = append(nodes, e.Values[i])
	}

	for i := range e.Comments {
		nodes = append(nodes, e.Comments[i])
	}
	for i := range e.EndLineComments {
		nodes = append(nodes, e.EndLineComments[i])
	}
	if e.Annotations != nil {
		nodes = append(nodes, e.Annotations)
	}

	return nodes
}

func (e *Enum) IsBadNode() bool {
	return e.BadNode
}

func (e *Enum) ChildrenBadNode() bool {
	children := e.Children()
	for i := range children {
		if children[i].IsBadNode() {
			return true
		}
		if children[i].ChildrenBadNode() {
			return true
		}
	}
	return false
}

func (e *Enum) SetLocation(loc Location) {
	e.Location = loc
}

func (e *Enum) Equals(node Node) bool {
	en, ok := node.(*Enum)
	if !ok {
		return false
	}

	if (e == nil && en != nil) ||
		(e != nil && en == nil) {
		return false
	} else if e == nil && en == nil {
		return true
	}

	if e.BadNode != en.BadNode {
		return false
	}

	if !e.EnumKeyword.Equals(en.EnumKeyword) {
		return false
	}

	if !e.LCurKeyword.Equals(en.LCurKeyword) {
		return false
	}

	if !e.RCurKeyword.Equals(en.RCurKeyword) {
		return false
	}

	if !e.Name.Equals(en.Name) {
		return false
	}
	if len(e.Values) != len(en.Values) {
		return false
	}
	for i := range e.Values {
		if !e.Values[i].Equals(en.Values[i]) {
			return false
		}
	}

	if len(e.Comments) != len(en.Comments) {
		return false
	}

	for i := range e.Comments {
		if !e.Comments[i].Equals(en.Comments[i]) {
			return false
		}
	}

	if len(e.EndLineComments) != len(en.EndLineComments) {
		return false
	}

	for i := range e.EndLineComments {
		if !e.EndLineComments[i].Equals(en.EndLineComments[i]) {
			return false
		}
	}

	if !e.Annotations.Equals(en.Annotations) {
		return false
	}

	return true
}

type EnumValue struct {
	ListSeparatorKeyword *ListSeparatorKeyword // can be nil
	EqualKeyword         *EqualKeyword         // can be nil
	Name                 *Identifier
	ValueNode            *ConstValue
	Value                int64 // Value only record enum value. it is not a ast node
	Annotations          *Annotations
	Comments             []*Comment
	EndLineComments      []*Comment

	BadNode bool
	Location
}

func NewBadEnumValue(loc Location) *EnumValue {
	return &EnumValue{
		BadNode:  true,
		Location: loc,
	}
}

func NewEnumValue(listSeparatorKeyword *ListSeparatorKeyword, equalKeyword *EqualKeyword, name *Identifier, valueNode *ConstValue, value int64, annotations *Annotations, loc Location) *EnumValue {
	return &EnumValue{
		ListSeparatorKeyword: listSeparatorKeyword,
		EqualKeyword:         equalKeyword,
		Name:                 name,
		ValueNode:            valueNode,
		Value:                value,
		Annotations:          annotations,
		Location:             loc,
	}
}

func (e *EnumValue) Children() []Node {
	nodes := []Node{e.Name}
	if e.ValueNode != nil {
		nodes = append(nodes, e.ValueNode)
	}
	if e.ListSeparatorKeyword != nil {
		nodes = append(nodes, e.ListSeparatorKeyword)
	}
	if e.EqualKeyword != nil {
		nodes = append(nodes, e.EqualKeyword)
	}
	for i := range e.Comments {
		nodes = append(nodes, e.Comments[i])
	}
	for i := range e.EndLineComments {
		nodes = append(nodes, e.EndLineComments[i])
	}
	if e.Annotations != nil {
		nodes = append(nodes, e.Annotations)
	}

	return nodes
}

func (e *EnumValue) Type() string {
	return "EnumValue"
}

func (e *EnumValue) SetComments(comments []*Comment, endLineComments []*Comment) {
	e.Comments = comments
	e.EndLineComments = endLineComments
}

func (e *EnumValue) IsBadNode() bool {
	return e.BadNode
}

func (e *EnumValue) ChildrenBadNode() bool {
	children := e.Children()
	for i := range children {
		if children[i].IsBadNode() {
			return true
		}
		if children[i].ChildrenBadNode() {
			return true
		}
	}
	return false
}

func (e *EnumValue) Equals(node Node) bool {
	en, ok := node.(*EnumValue)
	if !ok {
		return false
	}

	if (e == nil && en != nil) ||
		(e != nil && en == nil) {
		return false
	} else if e == nil && en == nil {
		return true
	}

	if e.BadNode != en.BadNode {
		return true
	}

	// comma 不影响语义
	// if !e.ListSeparatorKeyword.Equals(en.ListSeparatorKeyword) {
	// 	return false
	// }

	if !e.EqualKeyword.Equals(en.EqualKeyword) {
		return false
	}

	if !e.Name.Equals(en.Name) {
		return false
	}

	if !e.ValueNode.Equals(en.ValueNode) {
		return false
	}

	if e.Value != en.Value {
		return false
	}

	if len(e.Comments) != len(en.Comments) {
		return false
	}

	for i := range e.Comments {
		if !e.Comments[i].Equals(en.Comments[i]) {
			return false
		}
	}

	if len(e.EndLineComments) != len(en.EndLineComments) {
		return false
	}

	for i := range e.EndLineComments {
		if !e.EndLineComments[i].Equals(en.EndLineComments[i]) {
			return false
		}
	}

	if !e.Annotations.Equals(en.Annotations) {
		return false
	}

	return true
}

type ServiceKeyword struct {
	Keyword
}

func (s *ServiceKeyword) Type() string {
	return "ServiceKeyword"
}

func (s *ServiceKeyword) Equals(node Node) bool {
	sn, ok := node.(*ServiceKeyword)
	if !ok {
		return false
	}

	if (s == nil && sn != nil) ||
		(s != nil && sn == nil) {
		return false
	} else if s == nil && sn == nil {
		return true
	}

	return s.Keyword.Equals(&sn.Keyword)
}

type ExtendsKeyword struct {
	Keyword
}

func (s *ExtendsKeyword) Type() string {
	return "ExtendsKeyword"
}

func (s *ExtendsKeyword) Equals(node Node) bool {
	sn, ok := node.(*ExtendsKeyword)
	if !ok {
		return false
	}

	if (s == nil && sn != nil) ||
		(s != nil && sn == nil) {
		return false
	} else if s == nil && sn == nil {
		return true
	}

	return s.Keyword.Equals(&sn.Keyword)
}

type Service struct {
	ServiceKeyword *ServiceKeyword
	ExtendsKeyword *ExtendsKeyword // can be nil
	LCurKeyword    *LCurKeyword
	RCurKeyword    *RCurKeyword
	Name           *Identifier
	Extends        *Identifier
	Functions      []*Function

	Comments        []*Comment
	EndLineComments []*Comment
	Annotations     *Annotations

	BadNode bool
	Location
}

func NewService(serviceKeyword *ServiceKeyword, extendsKeyword *ExtendsKeyword, lCurKeyword *LCurKeyword, rCurKeyword *RCurKeyword, name *Identifier, extends *Identifier, fns []*Function, loc Location) *Service {
	return &Service{
		ServiceKeyword: serviceKeyword,
		ExtendsKeyword: extendsKeyword,
		LCurKeyword:    lCurKeyword,
		RCurKeyword:    rCurKeyword,
		Name:           name,
		Extends:        extends,
		Functions:      fns,
		Location:       loc,
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

func (s *Service) SetComments(comments []*Comment, endLineComments []*Comment) {
	s.Comments = comments
	s.EndLineComments = endLineComments
}

func (s *Service) SetAnnotations(annos *Annotations) {
	s.Annotations = annos
}

func (s *Service) Children() []Node {
	nodes := []Node{s.ServiceKeyword, s.LCurKeyword, s.RCurKeyword}
	if s.ExtendsKeyword != nil {
		nodes = append(nodes, s.ExtendsKeyword)
	}
	if s.Name != nil {
		nodes = append(nodes, s.Name)
	}
	if s.Extends != nil {
		nodes = append(nodes, s.Extends)
	}
	for i := range s.Functions {
		nodes = append(nodes, s.Functions[i])
	}

	for i := range s.Comments {
		nodes = append(nodes, s.Comments[i])
	}
	for i := range s.EndLineComments {
		nodes = append(nodes, s.EndLineComments[i])
	}
	if s.Annotations != nil {
		nodes = append(nodes, s.Annotations)
	}

	return nodes
}

func (s *Service) IsBadNode() bool {
	return s.BadNode
}

func (s *Service) ChildrenBadNode() bool {
	children := s.Children()
	for i := range children {
		if children[i].IsBadNode() {
			return true
		}
		if children[i].ChildrenBadNode() {
			return true
		}
	}
	return false
}

func (s *Service) SetLocation(loc Location) {
	s.Location = loc
}

func (s *Service) Equals(node Node) bool {
	sn, ok := node.(*Service)
	if !ok {
		return false
	}

	if (s == nil && sn != nil) ||
		(s != nil && sn == nil) {
		return false
	} else if s == nil && sn == nil {
		return true
	}

	if s.BadNode != sn.BadNode {
		return false
	}

	if !s.ServiceKeyword.Equals(sn.ServiceKeyword) {
		return false
	}

	if !s.ExtendsKeyword.Equals(sn.ExtendsKeyword) {
		return false
	}

	if !s.LCurKeyword.Equals(sn.LCurKeyword) {
		return false
	}

	if !s.RCurKeyword.Equals(sn.RCurKeyword) {
		return false
	}

	if !s.Name.Equals(sn.Name) {
		return false
	}

	if !s.Extends.Equals(sn.Extends) {
		return false
	}

	if len(s.Functions) != len(sn.Functions) {
		return false
	}

	for i := range s.Functions {
		if !s.Functions[i].Equals(sn.Functions[i]) {
			return false
		}
	}

	if len(s.Comments) != len(sn.Comments) {
		return false
	}

	for i := range s.Comments {
		if !s.Comments[i].Equals(sn.Comments[i]) {
			return false
		}
	}

	if len(s.EndLineComments) != len(sn.EndLineComments) {
		return false
	}

	for i := range s.EndLineComments {
		if !s.EndLineComments[i].Equals(sn.EndLineComments[i]) {
			return false
		}
	}

	if !s.Annotations.Equals(sn.Annotations) {
		return false
	}

	return true
}

type OnewayKeyword struct {
	Keyword
}

func (o *OnewayKeyword) Type() string {
	return "OnewayKeyword"
}

func (o *OnewayKeyword) Equals(node Node) bool {
	on, ok := node.(*OnewayKeyword)
	if !ok {
		return false
	}

	if (o == nil && on != nil) ||
		(o != nil && on == nil) {
		return false
	} else if o == nil && on == nil {
		return true
	}

	return o.Keyword.Equals(&on.Keyword)
}

type LParKeyword struct {
	Keyword
}

func (l *LParKeyword) Type() string {
	return "LParKeyword"
}

func (l *LParKeyword) Equals(node Node) bool {
	ln, ok := node.(*LParKeyword)
	if !ok {
		return false
	}

	if (l == nil && ln != nil) ||
		(l != nil && ln == nil) {
		return false
	} else if l == nil && ln == nil {
		return true
	}

	return l.Keyword.Equals(&ln.Keyword)
}

type RParKeyword struct {
	Keyword
}

func (r *RParKeyword) Type() string {
	return "RParKeyword"
}

func (r *RParKeyword) Equals(node Node) bool {
	rn, ok := node.(*RParKeyword)
	if !ok {
		return false
	}

	if (r == nil && rn != nil) ||
		(r != nil && rn == nil) {
		return false
	} else if r == nil && rn == nil {
		return true
	}

	return r.Keyword.Equals(&rn.Keyword)
}

type VoidKeyword struct {
	Keyword
}

func (v *VoidKeyword) Type() string {
	return "VoidKeyword"
}

func (v *VoidKeyword) Equals(node Node) bool {
	vn, ok := node.(*VoidKeyword)
	if !ok {
		return false
	}

	if (v == nil && vn != nil) ||
		(v != nil && vn == nil) {
		return false
	} else if v == nil && vn == nil {
		return true
	}

	return v.Keyword.Equals(&vn.Keyword)
}

type ThrowsKeyword struct {
	Keyword
}

func (t *ThrowsKeyword) Type() string {
	return "ThrowsKeyword"
}

func (t *ThrowsKeyword) Equals(node Node) bool {
	tn, ok := node.(*ThrowsKeyword)
	if !ok {
		return false
	}

	if (t == nil && tn != nil) ||
		(t != nil && tn == nil) {
		return false
	} else if t == nil && tn == nil {
		return true
	}

	return t.Keyword.Equals(&tn.Keyword)
}

type Throws struct {
	ThrowsKeyword *ThrowsKeyword
	LParKeyword   *LParKeyword
	RParKeyword   *RParKeyword

	Fields []*Field

	BadNode bool
	Location
}

func NewThrows(throwsKeyword *ThrowsKeyword, lparKeyword *LParKeyword, rparKeyword *RParKeyword, fields []*Field, loc Location) *Throws {
	return &Throws{
		ThrowsKeyword: throwsKeyword,
		LParKeyword:   lparKeyword,
		RParKeyword:   rparKeyword,
		Fields:        fields,
		Location:      loc,
	}
}

func (t *Throws) Type() string {
	return "Throws"
}

func (t *Throws) IsBadNode() bool {
	return t.BadNode
}

func (t *Throws) ChildrenBadNode() bool {
	children := t.Children()
	for i := range children {
		if children[i].IsBadNode() {
			return true
		}
		if children[i].ChildrenBadNode() {
			return true
		}
	}
	return false
}

func (t *Throws) Children() []Node {
	nodes := []Node{t.ThrowsKeyword, t.LParKeyword, t.RParKeyword}
	for i := range t.Fields {
		nodes = append(nodes, t.Fields[i])
	}
	return nodes
}

func (t *Throws) Equals(node Node) bool {
	tn, ok := node.(*Throws)
	if !ok {
		return false
	}

	if (t == nil && tn != nil) ||
		(t != nil && tn == nil) {
		return false
	} else if t == nil && tn == nil {
		return true
	}

	if t.BadNode != tn.BadNode {
		return false
	}

	if !t.ThrowsKeyword.Equals(tn.ThrowsKeyword) {
		return false
	}

	if !t.LParKeyword.Equals(tn.LParKeyword) {
		return false
	}

	if !t.RParKeyword.Equals(tn.RParKeyword) {
		return false
	}

	if len(t.Fields) != len(tn.Fields) {
		return false
	}

	for i := range t.Fields {
		if !t.Fields[i].Equals(tn.Fields[i]) {
			return false
		}
	}

	return true
}

type Function struct {
	LParKeyword          *LParKeyword
	RParKeyword          *RParKeyword
	ListSeparatorKeyword *ListSeparatorKeyword // can be nil
	Name                 *Identifier
	Oneway               *OnewayKeyword // can be nil
	Void                 *VoidKeyword   // can be nil
	FunctionType         *FieldType
	Arguments            []*Field
	Throws               *Throws
	Comments             []*Comment
	EndLineComments      []*Comment
	Annotations          *Annotations

	BadNode bool
	Location
}

func NewFunction(lParKeyword *LParKeyword, rParKeyword *RParKeyword, listSeparatorKeyword *ListSeparatorKeyword, name *Identifier, oneway *OnewayKeyword, void *VoidKeyword, ft *FieldType, args []*Field, throws *Throws, comments []*Comment, endlineComments []*Comment, annotations *Annotations, loc Location) *Function {
	return &Function{
		LParKeyword:          lParKeyword,
		RParKeyword:          rParKeyword,
		ListSeparatorKeyword: listSeparatorKeyword,
		Name:                 name,
		Oneway:               oneway,
		Void:                 void,
		FunctionType:         ft,
		Arguments:            args,
		Throws:               throws,
		Comments:             comments,
		EndLineComments:      endlineComments,
		Annotations:          annotations,
		Location:             loc,
	}
}

func NewBadFunction(loc Location) *Function {
	return &Function{
		BadNode:  true,
		Location: loc,
	}
}

func (f *Function) Children() []Node {
	nodes := []Node{f.LParKeyword, f.RParKeyword}
	if f.Oneway != nil {
		nodes = append(nodes, f.Oneway)
	}
	if f.Void != nil {
		nodes = append(nodes, f.Void)
	}
	if f.ListSeparatorKeyword != nil {
		nodes = append(nodes, f.ListSeparatorKeyword)
	}
	if f.Name != nil {
		nodes = append(nodes, f.Name)
	}
	if f.FunctionType != nil {
		nodes = append(nodes, f.FunctionType)
	}
	for i := range f.Arguments {
		nodes = append(nodes, f.Arguments[i])
	}
	if f.Throws != nil {
		nodes = append(nodes, f.Throws)
	}
	for i := range f.Comments {
		nodes = append(nodes, f.Comments[i])
	}
	for i := range f.EndLineComments {
		nodes = append(nodes, f.EndLineComments[i])
	}
	if f.Annotations != nil {
		nodes = append(nodes, f.Annotations)
	}

	return nodes
}

func (f *Function) Type() string {
	return "Function"
}

func (f *Function) IsBadNode() bool {
	return f.BadNode
}

func (f *Function) ChildrenBadNode() bool {
	children := f.Children()
	for i := range children {
		if children[i].IsBadNode() {
			return true
		}
		if children[i].ChildrenBadNode() {
			return true
		}
	}
	return false
}

func (f *Function) Equals(node Node) bool {
	fn, ok := node.(*Function)
	if !ok {
		return false
	}

	if (f == nil && fn != nil) ||
		(f != nil && fn == nil) {
		return false
	} else if f == nil && fn == nil {
		return true
	}

	if f.BadNode != fn.BadNode {
		return false
	}

	if !f.LParKeyword.Equals(fn.LParKeyword) {
		return false
	}

	if !f.RParKeyword.Equals(fn.RParKeyword) {
		return false
	}

	if !f.ListSeparatorKeyword.Equals(fn.ListSeparatorKeyword) {
		return false
	}

	if !f.Name.Equals(fn.Name) {
		return false
	}

	if !f.Oneway.Equals(fn.Oneway) {
		return false
	}

	if !f.Void.Equals(fn.Void) {
		return false
	}

	if !f.FunctionType.Equals(fn.FunctionType) {
		return false
	}

	if len(f.Arguments) != len(fn.Arguments) {
		return false
	}

	for i := range f.Arguments {
		if !f.Arguments[i].Equals(fn.Arguments[i]) {
			return false
		}
	}

	if !f.Throws.Equals(fn.Throws) {
		return false
	}

	if len(f.Comments) != len(fn.Comments) {
		return false
	}

	for i := range f.Comments {
		if !f.Comments[i].Equals(fn.Comments[i]) {
			return false
		}
	}

	if len(f.EndLineComments) != len(fn.EndLineComments) {
		return false
	}

	for i := range f.EndLineComments {
		if !f.EndLineComments[i].Equals(fn.EndLineComments[i]) {
			return false
		}
	}

	if !f.Annotations.Equals(fn.Annotations) {
		return false
	}

	return true
}

type UnionKeyword struct {
	Keyword
}

func (u *UnionKeyword) Type() string {
	return "UnionKeyword"
}

func (u *UnionKeyword) Equals(node Node) bool {
	un, ok := node.(*UnionKeyword)
	if !ok {
		return false
	}

	if (u == nil && un != nil) ||
		(u != nil && un == nil) {
		return false
	} else if u == nil && un == nil {
		return true
	}

	return u.Keyword.Equals(&un.Keyword)
}

type Union struct {
	UnionKeyword *UnionKeyword
	LCurKeyword  *LCurKeyword
	RCurKeyword  *RCurKeyword
	Name         *Identifier
	Fields       []*Field

	Comments        []*Comment
	EndLineComments []*Comment
	Annotations     *Annotations

	BadNode bool
	Location
}

func NewUnion(unionKeyword *UnionKeyword, lCurKeyword *LCurKeyword, rCurKeyword *RCurKeyword, name *Identifier, fields []*Field, loc Location) *Union {
	return &Union{
		UnionKeyword: unionKeyword,
		LCurKeyword:  lCurKeyword,
		RCurKeyword:  rCurKeyword,
		Name:         name,
		Fields:       fields,
		Location:     loc,
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

func (u *Union) SetComments(comments []*Comment, endLineComments []*Comment) {
	u.Comments = comments
	u.EndLineComments = endLineComments
}

func (u *Union) SetAnnotations(annos *Annotations) {
	u.Annotations = annos
}

func (u *Union) Children() []Node {
	nodes := []Node{u.Name, u.UnionKeyword, u.LCurKeyword, u.RCurKeyword}
	for i := range u.Fields {
		nodes = append(nodes, u.Fields[i])
	}
	for i := range u.Comments {
		nodes = append(nodes, u.Comments[i])
	}
	for i := range u.EndLineComments {
		nodes = append(nodes, u.EndLineComments[i])
	}

	for i := range u.Comments {
		nodes = append(nodes, u.Comments[i])
	}
	for i := range u.EndLineComments {
		nodes = append(nodes, u.EndLineComments[i])
	}
	if u.Annotations != nil {
		nodes = append(nodes, u.Annotations)
	}

	return nodes
}

func (u *Union) IsBadNode() bool {
	return u.BadNode
}

func (u *Union) ChildrenBadNode() bool {
	children := u.Children()
	for i := range children {
		if children[i].IsBadNode() {
			return true
		}
		if children[i].ChildrenBadNode() {
			return true
		}
	}
	return false
}

func (u *Union) SetLocation(loc Location) {
	u.Location = loc
}

func (u *Union) Equals(node Node) bool {
	un, ok := node.(*Union)
	if !ok {
		return false
	}

	if (u == nil && un != nil) ||
		(u != nil && un == nil) {
		return false
	} else if u == nil && un == nil {
		return true
	}

	if !u.UnionKeyword.Equals(un.UnionKeyword) {
		return false
	}

	if !u.LCurKeyword.Equals(un.LCurKeyword) {
		return false
	}

	if !u.RCurKeyword.Equals(un.RCurKeyword) {
		return false
	}

	if !u.Name.Equals(un.Name) {
		return false
	}

	if len(u.Fields) != len(un.Fields) {
		return false
	}

	for i := range u.Fields {
		if !u.Fields[i].Equals(un.Fields[i]) {
			return false
		}
	}

	if len(u.Comments) != len(un.Comments) {
		return false
	}

	for i := range u.Comments {
		if !u.Comments[i].Equals(un.Comments[i]) {
			return false
		}
	}

	if len(u.EndLineComments) != len(un.EndLineComments) {
		return false
	}

	for i := range u.EndLineComments {
		if !u.EndLineComments[i].Equals(un.EndLineComments[i]) {
			return false
		}
	}

	if !u.Annotations.Equals(un.Annotations) {
		return false
	}

	return true
}

type ExceptionKeyword struct {
	Keyword
}

func (e *ExceptionKeyword) Type() string {
	return "ExceptionKeyword"
}

func (e *ExceptionKeyword) Equals(node Node) bool {
	en, ok := node.(*ExceptionKeyword)
	if !ok {
		return false
	}

	if (e == nil && en != nil) ||
		(e != nil && en == nil) {
		return false
	} else if e == nil && en == nil {
		return true
	}

	return e.Keyword.Equals(&en.Keyword)
}

type Exception struct {
	ExceptionKeyword *ExceptionKeyword
	LCurKeyword      *LCurKeyword
	RCurKeyword      *RCurKeyword
	Name             *Identifier
	Fields           []*Field

	Comments        []*Comment
	EndLineComments []*Comment
	Annotations     *Annotations

	BadNode bool
	Location
}

func NewException(exceptionKeyword *ExceptionKeyword, lCurKeyword *LCurKeyword, rCurKeyword *RCurKeyword, name *Identifier, fields []*Field, loc Location) *Exception {
	return &Exception{
		ExceptionKeyword: exceptionKeyword,
		LCurKeyword:      lCurKeyword,
		RCurKeyword:      rCurKeyword,
		Name:             name,
		Fields:           fields,
		Location:         loc,
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

func (e *Exception) SetComments(comments []*Comment, endLineComments []*Comment) {
	e.Comments = comments
	e.EndLineComments = endLineComments
}

func (e *Exception) SetAnnotations(annos *Annotations) {
	e.Annotations = annos
}

func (e *Exception) Children() []Node {
	nodes := []Node{e.Name, e.ExceptionKeyword, e.LCurKeyword, e.RCurKeyword}
	for i := range e.Fields {
		nodes = append(nodes, e.Fields[i])
	}
	for i := range e.Comments {
		nodes = append(nodes, e.Comments[i])
	}
	for i := range e.EndLineComments {
		nodes = append(nodes, e.EndLineComments[i])
	}

	if e.Annotations != nil {
		nodes = append(nodes, e.Annotations)
	}

	return nodes
}

func (e *Exception) IsBadNode() bool {
	return e.BadNode
}

func (e *Exception) ChildrenBadNode() bool {
	children := e.Children()
	for i := range children {
		if children[i].IsBadNode() {
			return true
		}
		if children[i].ChildrenBadNode() {
			return true
		}
	}
	return false
}

func (e *Exception) SetLocation(loc Location) {
	e.Location = loc
}

func (e *Exception) Equals(node Node) bool {
	en, ok := node.(*Exception)
	if !ok {
		return false
	}

	if (e == nil && en != nil) ||
		(e != nil && en == nil) {
		return false
	} else if e == nil && en == nil {
		return true
	}

	if e.BadNode != en.BadNode {
		return false
	}

	if !e.ExceptionKeyword.Equals(en.ExceptionKeyword) {
		return false
	}

	if !e.LCurKeyword.Equals(en.LCurKeyword) {
		return false
	}

	if !e.RCurKeyword.Equals(en.RCurKeyword) {
		return false
	}

	if !e.Name.Equals(en.Name) {
		return false
	}

	if len(e.Fields) != len(en.Fields) {
		return false
	}

	for i := range e.Fields {
		if !e.Fields[i].Equals(en.Fields[i]) {
			return false
		}
	}

	if len(e.Comments) != len(en.Comments) {
		return false
	}

	for i := range e.Comments {
		if !e.Comments[i].Equals(en.Comments[i]) {
			return false
		}
	}

	if len(e.EndLineComments) != len(en.EndLineComments) {
		return false
	}

	for i := range e.EndLineComments {
		if !e.EndLineComments[i].Equals(en.EndLineComments[i]) {
			return false
		}
	}

	if !e.Annotations.Equals(en.Annotations) {
		return false
	}

	return true
}

type IdentifierName struct {
	Text string

	BadNode bool
	Location
}

func NewIdentifierName(name string, loc Location) *IdentifierName {
	return &IdentifierName{
		Text:     name,
		Location: loc,
		BadNode:  name == "",
	}
}

func (i *IdentifierName) Children() []Node {
	return nil
}

func (i *IdentifierName) Type() string {
	return "IdentifierName"
}

func (i *IdentifierName) IsBadNode() bool {
	return i.BadNode
}

func (i *IdentifierName) ChildrenBadNode() bool {
	children := i.Children()
	for i := range children {
		if children[i].IsBadNode() {
			return true
		}
		if children[i].ChildrenBadNode() {
			return true
		}
	}
	return false
}

func (i *IdentifierName) Equals(node Node) bool {
	in, ok := node.(*IdentifierName)
	if !ok {
		return false
	}

	if (i == nil && in != nil) ||
		(i != nil && in == nil) {
		return false
	} else if i == nil && in == nil {
		return true
	}

	if i.BadNode != in.BadNode {
		return false
	}

	if i.Text != in.Text {
		return false
	}

	return true
}

type Identifier struct {
	Name     *IdentifierName
	Comments []*Comment

	BadNode bool
	Location
}

func NewIdentifier(name *IdentifierName, comments []*Comment, loc Location) *Identifier {
	id := &Identifier{
		Name:     name,
		Comments: comments,
		Location: loc,
		BadNode:  name == nil || name.BadNode,
	}

	return id
}

func NewBadIdentifier(loc Location) *Identifier {
	return &Identifier{
		BadNode:  true,
		Location: loc,
	}
}

func (i *Identifier) ToFieldType() *FieldType {
	t := &FieldType{
		TypeName: &TypeName{
			Name:     i.Name.Text,
			Location: i.Name.Location,
		},
		Location: i.Location,
	}

	return t
}

func (i *Identifier) Children() []Node {
	var nodes []Node
	for _, com := range i.Comments {
		nodes = append(nodes, com)
	}
	nodes = append(nodes, i.Name)
	return nodes
}

func (i *Identifier) Type() string {
	return "Identifier"
}

func (i *Identifier) IsBadNode() bool {
	return i.BadNode
}

func (i *Identifier) ChildrenBadNode() bool {
	children := i.Children()
	for i := range children {
		if children[i].IsBadNode() {
			return true
		}
		if children[i].ChildrenBadNode() {
			return true
		}
	}
	return false
}

func (i *Identifier) Equals(node Node) bool {
	in, ok := node.(*Identifier)
	if !ok {
		return false
	}

	if (i == nil && in != nil) ||
		(i != nil && in == nil) {
		return false
	} else if i == nil && in == nil {
		return true
	}

	if i.BadNode != in.BadNode {
		return false
	}

	if !i.Name.Equals(in.Name) {
		return false
	}

	if len(i.Comments) != len(in.Comments) {
		return false
	}

	for n := range i.Comments {
		if !i.Comments[n].Equals(in.Comments[n]) {
			return false
		}
	}

	return true
}

func ConvertPosition(pos position) Position {
	return Position{
		Line:   pos.line,
		Col:    pos.col,
		Offset: pos.offset,
	}
}

type Field struct {
	Index           *FieldIndex
	RequiredKeyword *RequiredKeyword
	FieldType       *FieldType
	Identifier      *Identifier
	ConstValue      *ConstValue

	EqualKeyword         *EqualKeyword         // can be nil
	ListSeparatorKeyword *ListSeparatorKeyword // can be nil

	Comments        []*Comment
	EndLineComments []*Comment
	Annotations     *Annotations

	BadNode bool
	Location
}

func NewField(equalKeyword *EqualKeyword, listSeparatorKeyword *ListSeparatorKeyword, comments []*Comment, endLineComments []*Comment, annotations *Annotations, index *FieldIndex, required *RequiredKeyword, fieldType *FieldType, identifier *Identifier, constValue *ConstValue, loc Location) *Field {
	field := &Field{
		EqualKeyword:         equalKeyword,
		ListSeparatorKeyword: listSeparatorKeyword,
		Comments:             comments,
		EndLineComments:      endLineComments,
		Annotations:          annotations,
		Index:                index,
		RequiredKeyword:      required,
		FieldType:            fieldType,
		Identifier:           identifier,
		ConstValue:           constValue,
		Location:             loc,
	}
	return field
}

func NewBadField(loc Location) *Field {
	return &Field{
		BadNode:  true,
		Location: loc,
	}
}

func (f *Field) Children() []Node {
	var res []Node
	if f.RequiredKeyword != nil {
		res = append(res, f.RequiredKeyword)
	}
	if f.FieldType != nil {
		res = append(res, f.FieldType)
	}
	if f.Identifier != nil {
		res = append(res, f.Identifier)
	}
	if f.ConstValue != nil {
		res = append(res, f.ConstValue)
	}
	if f.EqualKeyword != nil {
		res = append(res, f.EqualKeyword)
	}
	if f.ListSeparatorKeyword != nil {
		res = append(res, f.ListSeparatorKeyword)
	}
	for i := range f.Comments {
		res = append(res, f.Comments[i])
	}
	for i := range f.EndLineComments {
		res = append(res, f.EndLineComments[i])
	}
	if f.Annotations != nil {
		res = append(res, f.Annotations)
	}
	return res
}

func (f *Field) Type() string {
	return "Field"
}

func (f *Field) IsBadNode() bool {
	return f.BadNode
}

func (f *Field) ChildrenBadNode() bool {
	children := f.Children()
	for i := range children {
		if children[i].IsBadNode() {
			return true
		}
		if children[i].ChildrenBadNode() {
			return true
		}
	}
	return false
}

func (f *Field) Equals(node Node) bool {
	fn, ok := node.(*Field)
	if !ok {
		return false
	}

	if (f == nil && fn != nil) ||
		(f != nil && fn == nil) {
		return false
	} else if f == nil && fn == nil {
		return true
	}

	if f.BadNode != fn.BadNode {
		return false
	}

	if !f.Index.Equals(fn.Index) {
		return false
	}

	if !f.RequiredKeyword.Equals(fn.RequiredKeyword) {
		return false
	}

	if !f.FieldType.Equals(fn.FieldType) {
		return false
	}

	if !f.Identifier.Equals(fn.Identifier) {
		return false
	}

	if !f.ConstValue.Equals(fn.ConstValue) {
		return false
	}

	if !f.EqualKeyword.Equals(fn.EqualKeyword) {
		return false
	}

	// 末尾的 , 不影响语义，暂时注释掉
	// if !f.ListSeparatorKeyword.Equals(fn.ListSeparatorKeyword) {
	// 	return false
	// }

	if len(f.Comments) != len(fn.Comments) {
		return false
	}

	for i := range f.Comments {
		if !f.Comments[i].Equals(fn.Comments[i]) {
			return false
		}
	}

	if len(f.EndLineComments) != len(fn.EndLineComments) {
		return false
	}

	for i := range f.EndLineComments {
		if !f.EndLineComments[i].Equals(fn.EndLineComments[i]) {
			return false
		}
	}

	if !f.Annotations.Equals(fn.Annotations) {
		return false
	}

	return true
}

type ColonKeyword struct {
	Keyword
}

func (c *ColonKeyword) Type() string {
	return "ColonKeyword"
}

func (c *ColonKeyword) Equals(node Node) bool {
	cn, ok := node.(*ColonKeyword)
	if !ok {
		return false
	}

	if (c == nil && cn != nil) ||
		(c != nil && cn == nil) {
		return false
	} else if c == nil && cn == nil {
		return true
	}

	return c.Keyword.Equals(&cn.Keyword)
}

type FieldIndex struct {
	ColonKeyword *ColonKeyword
	Value        int

	Comments []*Comment

	BadNode bool
	Location
}

func NewFieldIndex(ColonKeyword *ColonKeyword, v int, comments []*Comment, loc Location) *FieldIndex {
	return &FieldIndex{
		ColonKeyword: ColonKeyword,
		Value:        v,
		Comments:     comments,
		Location:     loc,
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

func (f *FieldIndex) IsBadNode() bool {
	return f.BadNode
}

func (f *FieldIndex) ChildrenBadNode() bool {
	children := f.Children()
	for i := range children {
		if children[i].IsBadNode() {
			return true
		}
		if children[i].ChildrenBadNode() {
			return true
		}
	}
	return false
}

func (f *FieldIndex) Equals(node Node) bool {
	fn, ok := node.(*FieldIndex)
	if !ok {
		return false
	}

	if (f == nil && fn != nil) ||
		(f != nil && fn == nil) {
		return false
	} else if f == nil && fn == nil {
		return true
	}

	if f.BadNode != fn.BadNode {
		return false
	}

	if !f.ColonKeyword.Equals(fn.ColonKeyword) {
		return false
	}

	if f.Value != fn.Value {
		return false
	}

	if len(f.Comments) != len(fn.Comments) {
		return false
	}

	for i := range f.Comments {
		if !f.Comments[i].Equals(fn.Comments[i]) {
			return false
		}
	}

	return true
}

type RequiredKeyword struct {
	Keyword
}

func (r *RequiredKeyword) Type() string {
	return "RequiredKeyword"
}

func (r *RequiredKeyword) Equals(node Node) bool {
	rn, ok := node.(*RequiredKeyword)
	if !ok {
		return false
	}

	if (r == nil && rn != nil) ||
		(r != nil && rn == nil) {
		return false
	} else if r == nil && rn == nil {
		return true
	}

	return r.Keyword.Equals(&rn.Keyword)
}

type LPointKeyword struct {
	Keyword
}

func (l *LPointKeyword) Type() string {
	return "LPointKeyword"
}

func (l *LPointKeyword) Equals(node Node) bool {
	ln, ok := node.(*LPointKeyword)
	if !ok {
		return false
	}

	if (l == nil && ln != nil) ||
		(l != nil && ln == nil) {
		return false
	} else if l == nil && ln == nil {
		return true
	}

	return l.Keyword.Equals(&ln.Keyword)
}

type RPointKeyword struct {
	Keyword
}

func (r *RPointKeyword) Type() string {
	return "RPointKeyword"
}

func (r *RPointKeyword) Equals(node Node) bool {
	rn, ok := node.(*RPointKeyword)
	if !ok {
		return false
	}

	if (r == nil && rn != nil) ||
		(r != nil && rn == nil) {
		return false
	} else if r == nil && rn == nil {
		return true
	}

	return r.Keyword.Equals(&rn.Keyword)
}

type CommaKeyword struct {
	Keyword
}

func (r *CommaKeyword) Type() string {
	return "CommaKeyword"
}

func (r *CommaKeyword) Equals(node Node) bool {
	rn, ok := node.(*CommaKeyword)
	if !ok {
		return false
	}

	if (r == nil && rn != nil) ||
		(r != nil && rn == nil) {
		return false
	} else if r == nil && rn == nil {
		return true
	}

	return r.Keyword.Equals(&rn.Keyword)
}

type CppTypeKeyword struct {
	Keyword
}

func (r *CppTypeKeyword) Type() string {
	return "CppTypeKeyword"
}

func (r *CppTypeKeyword) Equals(node Node) bool {
	rn, ok := node.(*CppTypeKeyword)
	if !ok {
		return false
	}

	if (r == nil && rn != nil) ||
		(r != nil && rn == nil) {
		return false
	} else if r == nil && rn == nil {
		return true
	}

	return r.Keyword.Equals(&rn.Keyword)
}

type CppType struct {
	CppTypeKeyword *CppTypeKeyword
	Literal        *Literal

	BadNode bool
	Location
}

func NewCppType(cppTypeKeyword *CppTypeKeyword, literal *Literal, loc Location) *CppType {
	return &CppType{
		CppTypeKeyword: cppTypeKeyword,
		Literal:        literal,
		Location:       loc,
	}
}

func (c *CppType) Type() string {
	return "CppType"
}

func (c *CppType) Children() []Node {
	return []Node{c.CppTypeKeyword, c.Literal}
}

func (c *CppType) IsBadNode() bool {
	return c.IsBadNode()
}

func (c *CppType) ChildrenBadNode() bool {
	children := c.Children()
	for i := range children {
		if children[i].IsBadNode() {
			return true
		}
		if children[i].ChildrenBadNode() {
			return true
		}
	}
	return false
}

func (c *CppType) Equals(node Node) bool {
	cn, ok := node.(*CppType)
	if !ok {
		return false
	}
	if (c == nil && cn != nil) ||
		(c != nil && cn == nil) {
		return false
	} else if c == nil && cn == nil {
		return true
	}

	if c.BadNode != cn.BadNode {
		return false
	}

	if !c.CppTypeKeyword.Equals(cn.CppTypeKeyword) {
		return false
	}

	if !c.Literal.Equals(cn.Literal) {
		return false
	}

	return true
}

type FieldType struct {
	TypeName *TypeName
	// only exist when TypeName is map or set or list
	KeyType *FieldType
	// only exist when TypeName is map
	ValueType *FieldType

	// only exist in map, set, list. can be nil
	CppType *CppType

	// only exist in map, set, list
	LPointKeyword *LPointKeyword
	// only exist in map, set, list
	RPointKeyword *RPointKeyword
	// only exist in map
	CommaKeyword *CommaKeyword

	Annotations *Annotations

	BadNode bool

	Location
}

func NewFieldType(lpointKeyword *LPointKeyword, rpointKeyword *RPointKeyword, commaKeyword *CommaKeyword, cppType *CppType, typeName *TypeName, keyType *FieldType, valueType *FieldType, loc Location) *FieldType {
	return &FieldType{
		LPointKeyword: lpointKeyword,
		RPointKeyword: rpointKeyword,
		CommaKeyword:  commaKeyword,

		CppType: cppType,

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

func (c *FieldType) IsBadNode() bool {
	return c.BadNode
}

func (c *FieldType) ChildrenBadNode() bool {
	children := c.Children()
	for i := range children {
		if children[i].IsBadNode() {
			return true
		}
		if children[i].ChildrenBadNode() {
			return true
		}
	}
	return false
}

func (c *FieldType) Equals(node Node) bool {
	fn, ok := node.(*FieldType)
	if !ok {
		return false
	}

	if (c == nil && fn != nil) ||
		(c != nil && fn == nil) {
		return false
	} else if c == nil && fn == nil {
		return true
	}

	if c.BadNode != fn.BadNode {
		return false
	}

	if !c.TypeName.Equals(fn.TypeName) {
		return false
	}

	if !c.KeyType.Equals(fn.KeyType) {
		return false
	}

	if !c.ValueType.Equals(fn.ValueType) {
		return false
	}

	if !c.CppType.Equals(fn.CppType) {
		return false
	}

	if !c.LPointKeyword.Equals(fn.LPointKeyword) {
		return false
	}

	if !c.RPointKeyword.Equals(fn.RPointKeyword) {
		return false
	}

	if !c.CommaKeyword.Equals(fn.CommaKeyword) {
		return false
	}

	if !c.Annotations.Equals(fn.Annotations) {
		return false
	}

	return true
}

type TypeName struct {
	// TypeName can be:
	// container type: map, set, list
	// base type: bool, byte, i8, i16, i32, i64, double, string, binary, uuid
	// struct, enum, union, exception, identifier
	Name     string
	Comments []*Comment

	BadNode bool
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
	var nodes []Node
	for i := range t.Comments {
		nodes = append(nodes, t.Comments[i])
	}
	return nodes
}

func (t *TypeName) Type() string {
	return "TypeName"
}

func (t *TypeName) IsBadNode() bool {
	return t.BadNode
}

func (t *TypeName) ChildrenBadNode() bool {
	children := t.Children()
	for i := range children {
		if children[i].IsBadNode() {
			return true
		}
		if children[i].ChildrenBadNode() {
			return true
		}
	}
	return false
}

func (t *TypeName) Equals(node Node) bool {
	tn, ok := node.(*TypeName)
	if !ok {
		return false
	}

	if (t == nil && tn != nil) ||
		(t != nil && tn == nil) {
		return false
	} else if t == nil && tn == nil {
		return true
	}

	if t.BadNode != tn.BadNode {
		return false
	}

	if t.Name != tn.Name {
		return false
	}

	if len(t.Comments) != len(tn.Comments) {
		return false
	}

	for i := range t.Comments {
		if !t.Comments[i].Equals(tn.Comments[i]) {
			return false
		}
	}

	return true
}

type LBrkKeyword struct {
	Keyword
}

func (l *LBrkKeyword) Type() string {
	return "LBrkKeyword"
}

func (l *LBrkKeyword) Equals(node Node) bool {
	ln, ok := node.(*LBrkKeyword)
	if !ok {
		return false
	}

	if (l == nil && ln != nil) ||
		(l != nil && ln == nil) {
		return false
	} else if l == nil && ln == nil {
		return true
	}

	return l.Keyword.Equals(&ln.Keyword)
}

type RBrkKeyword struct {
	Keyword
}

func (l *RBrkKeyword) Type() string {
	return "RBrkKeyword"
}

func (l *RBrkKeyword) Equals(node Node) bool {
	ln, ok := node.(*RBrkKeyword)
	if !ok {
		return false
	}

	if (l == nil && ln != nil) ||
		(l != nil && ln == nil) {
		return false
	} else if l == nil && ln == nil {
		return true
	}

	return l.Keyword.Equals(&ln.Keyword)
}

type ConstValue struct {
	// TypeName can be: list, map, pair, string, identifier, i64, double
	TypeName string
	// Value is the actual value or identifier name
	Value any

	// ValueInText is the user input value
	// it is used for i64 and double type value
	ValueInText string

	// only exist when TypeName is map
	Key any

	// exist in list
	LBrkKeyword *LBrkKeyword
	RBrkKeyword *RBrkKeyword

	// exist in map
	LCurKeyword *LCurKeyword
	RCurKeyword *RCurKeyword

	// exist in list, map item
	ListSeparatorKeyword *ListSeparatorKeyword

	// exist in map item
	ColonKeyword *ColonKeyword

	Comments []*Comment

	BadNode bool
	Location
}

func NewConstValue(typeName string, value any, loc Location) *ConstValue {
	return &ConstValue{
		TypeName: typeName,
		Value:    value,
		Location: loc,
	}
}

func NewBadConstValue(loc Location) *ConstValue {
	return &ConstValue{
		BadNode:  true,
		Location: loc,
	}
}

func NewBadIntConstValue(loc Location) *ConstValue {
	return &ConstValue{
		TypeName: "i64",
		BadNode:  true,
		Value:    int64(0),
		Location: loc,
	}
}

func NewMapConstValue(key, value *ConstValue, loc Location) *ConstValue {
	return &ConstValue{
		TypeName: "pair",
		Key:      key,
		Value:    value,
		Location: loc,
	}
}

func (c *ConstValue) SetComments(comments []*Comment) {
	c.Comments = comments
}

// TODO(jpf): nodes of key, value
func (c *ConstValue) Children() []Node {
	return nil
}

func (c *ConstValue) Type() string {
	return "ConstValue"
}

func (c *ConstValue) IsBadNode() bool {
	return c.BadNode
}

func (c *ConstValue) ChildrenBadNode() bool {
	children := c.Children()
	for i := range children {
		if children[i].IsBadNode() {
			return true
		}
		if children[i].ChildrenBadNode() {
			return true
		}
	}
	return false
}

func (c *ConstValue) Equals(node Node) bool {
	cn, ok := node.(*ConstValue)
	if !ok {
		return false
	}

	if (c == nil && cn != nil) ||
		(c != nil && cn == nil) {
		return false
	} else if c == nil && cn == nil {
		return true
	}

	if c.BadNode != cn.BadNode {
		return false
	}

	if c.TypeName != cn.TypeName {
		return false
	}

	// TODO(jpf): any 的比对
	// if c.Value != cn.Value {
	// 	return false
	// }

	if c.ValueInText != cn.ValueInText {
		return false
	}

	// TODO(jpf): any 的比对
	// if c.Key != cn.Key {
	// 	return false
	// }

	if !c.LBrkKeyword.Equals(cn.LBrkKeyword) {
		return false
	}

	if !c.RBrkKeyword.Equals(cn.RBrkKeyword) {
		return false
	}

	if !c.LCurKeyword.Equals(cn.LCurKeyword) {
		return false
	}

	if !c.RCurKeyword.Equals(cn.RCurKeyword) {
		return false
	}

	if !c.ListSeparatorKeyword.Equals(cn.ListSeparatorKeyword) {
		return false
	}

	if !c.ColonKeyword.Equals(cn.ColonKeyword) {
		return false
	}

	if len(c.Comments) != len(cn.Comments) {
		return false
	}

	for i := range c.Comments {
		if !c.Comments[i].Equals(cn.Comments[i]) {
			return false
		}
	}

	return true
}

type LiteralValue struct {
	Text string

	BadNode bool
	Location
}

func NewLiteralValue(text string, loc Location) *LiteralValue {
	return &LiteralValue{
		Text:     text,
		Location: loc,
	}
}

func NewBadLiteralValue(loc Location) *LiteralValue {
	return &LiteralValue{
		BadNode:  true,
		Location: loc,
	}
}

func (l *LiteralValue) Children() []Node {
	return nil
}

func (l *LiteralValue) Type() string {
	return "LiteralValue"
}

func (l *LiteralValue) IsBadNode() bool {
	return l.BadNode
}

func (l *LiteralValue) ChildrenBadNode() bool {
	children := l.Children()
	for i := range children {
		if children[i].IsBadNode() {
			return true
		}
		if children[i].ChildrenBadNode() {
			return true
		}
	}
	return false
}

func (l *LiteralValue) Equals(node Node) bool {
	ln, ok := node.(*LiteralValue)
	if !ok {
		return false
	}

	if (l == nil && ln != nil) ||
		(l != nil && ln == nil) {
		return false
	} else if l == nil && ln == nil {
		return true
	}

	if l.BadNode != ln.BadNode {
		return false
	}

	if l.Text != ln.Text {
		return false
	}

	return true
}

type Literal struct {
	Value *LiteralValue

	Quote string // single for ', double for "

	Comments []*Comment

	BadNode bool
	Location
}

// TODO: 区分单引号还是双引号?
func NewLiteral(comments []*Comment, v *LiteralValue, quote string, loc Location) *Literal {
	return &Literal{
		Value:    v,
		Quote:    quote,
		Comments: comments,
		Location: loc,
	}
}

func NewBadLiteral(loc Location) *Literal {
	return &Literal{
		Location: loc,
		BadNode:  true,
	}
}

func (l *Literal) Children() []Node {
	var nodes []Node
	for i := range l.Comments {
		nodes = append(nodes, l.Comments[i])
	}
	if l.Value != nil {
		nodes = append(nodes, l.Value)
	}
	return nodes
}

func (l *Literal) Type() string {
	return "Literal"
}

func (l *Literal) IsBadNode() bool {
	return l.BadNode
}

func (l *Literal) ChildrenBadNode() bool {
	children := l.Children()
	for i := range children {
		if children[i].IsBadNode() {
			return true
		}
		if children[i].ChildrenBadNode() {
			return true
		}
	}
	return false
}

func (l *Literal) Equals(node Node) bool {
	ln, ok := node.(*Literal)
	if !ok {
		return false
	}

	if (l == nil && ln != nil) ||
		(l != nil && ln == nil) {
		return false
	} else if l == nil && ln == nil {
		return true
	}

	if l.BadNode != ln.BadNode {
		return false
	}

	if !l.Value.Equals(ln.Value) {
		return false
	}

	if l.Quote != ln.Quote {
		return false
	}

	if len(l.Comments) != len(ln.Comments) {
		return false
	}

	for i := range l.Comments {
		if !l.Comments[i].Equals(ln.Comments[i]) {
			return false
		}
	}

	return true
}

type Annotations struct {
	Annotations []*Annotation
	LParKeyword *LParKeyword
	RParKeyword *RParKeyword

	BadNode bool
	Location
}

func NewAnnotations(lpar *LParKeyword, rpar *RParKeyword, annos []*Annotation, loc Location) *Annotations {
	return &Annotations{
		LParKeyword: lpar,
		RParKeyword: rpar,
		Annotations: annos,
		Location:    loc,
	}
}

func (a *Annotations) Type() string {
	return "Annotations"
}

func (a *Annotations) Children() []Node {
	nodes := []Node{a.LParKeyword, a.RParKeyword}
	for i := range a.Annotations {
		nodes = append(nodes, a.Annotations[i])
	}

	return nodes
}

func (a *Annotations) IsBadNode() bool {
	return a.BadNode
}

func (a *Annotations) ChildrenBadNode() bool {
	children := a.Children()
	for i := range children {
		if children[i].IsBadNode() {
			return true
		}
		if children[i].ChildrenBadNode() {
			return true
		}
	}
	return false
}

func (a *Annotations) Equals(node Node) bool {
	an, ok := node.(*Annotations)
	if !ok {
		return false
	}

	if (a == nil && an != nil) ||
		(a != nil && an == nil) {
		return false
	} else if a == nil && an == nil {
		return true
	}

	if a.BadNode != an.BadNode {
		return false
	}

	if len(a.Annotations) != len(an.Annotations) {
		return false
	}

	for i := range a.Annotations {
		if !a.Annotations[i].Equals(an.Annotations[i]) {
			return false
		}
	}

	if !a.LParKeyword.Equals(an.LParKeyword) {
		return false
	}

	if !a.RParKeyword.Equals(an.RParKeyword) {
		return false
	}

	return true

}

type Annotation struct {
	EqualKeyword         *EqualKeyword
	ListSeparatorKeyword *ListSeparatorKeyword

	Identifier *Identifier
	Value      *Literal

	BadNode bool
	Location
}

func NewAnnotation(equalKeyword *EqualKeyword, listSeparatorKeyword *ListSeparatorKeyword, id *Identifier, value *Literal, loc Location) *Annotation {
	return &Annotation{
		EqualKeyword:         equalKeyword,
		ListSeparatorKeyword: listSeparatorKeyword,
		Identifier:           id,
		Value:                value,
		Location:             loc,
	}
}

func NewBadAnnotation(loc Location) *Annotation {
	return &Annotation{
		BadNode:  true,
		Location: loc,
	}
}

func (a *Annotation) Children() []Node {
	nodes := []Node{a.Identifier, a.Value, a.EqualKeyword}
	if a.ListSeparatorKeyword != nil {
		nodes = append(nodes, a.ListSeparatorKeyword)
	}

	return nodes
}

func (a *Annotation) Type() string {
	return "Annotation"
}

func (a *Annotation) IsBadNode() bool {
	return a.BadNode
}

func (a *Annotation) ChildrenBadNode() bool {
	children := a.Children()
	for i := range children {
		if children[i].IsBadNode() {
			return true
		}
		if children[i].ChildrenBadNode() {
			return true
		}
	}
	return false
}

func (a *Annotation) Equals(node Node) bool {
	an, ok := node.(*Annotation)
	if !ok {
		return false
	}

	if (a == nil && an != nil) ||
		(a != nil && an == nil) {
		return false
	} else if a == nil && an == nil {
		return true
	}

	if a.BadNode != an.BadNode {
		return false
	}

	if !a.EqualKeyword.Equals(an.EqualKeyword) {
		return false
	}

	if !a.ListSeparatorKeyword.Equals(an.ListSeparatorKeyword) {
		return false
	}

	if !a.Identifier.Equals(an.Identifier) {
		return false
	}

	if !a.Value.Equals(an.Value) {
		return false
	}

	return true

}

type CommentStyle string

const (
	CommentStyleShell      CommentStyle = "shell"
	CommentStyleMultiLine  CommentStyle = "multiline"
	CommentStyleSingleLine CommentStyle = "singleline"
)

type Comment struct {
	Text  string
	Style CommentStyle // shell: #xxx, multiline: /* *** */, singleline: // xxxxx

	BadNode bool
	Location
}

func NewComment(text string, style CommentStyle, loc Location) *Comment {
	return &Comment{
		Text:     text,
		Style:    style,
		Location: loc,
	}
}

func NewBadComment(loc Location) *Comment {
	return &Comment{
		BadNode:  true,
		Location: loc,
	}
}

func (c *Comment) Children() []Node {
	return nil
}

func (c *Comment) Type() string {
	return "Comment"
}

func (c *Comment) IsBadNode() bool {
	return c.BadNode
}

func (c *Comment) ChildrenBadNode() bool {
	return false
}

func (c *Comment) Equals(node Node) bool {
	cn, ok := node.(*Comment)
	if !ok {
		return false
	}

	if (c == nil && cn != nil) ||
		(c != nil && cn == nil) {
		return false
	} else if c == nil && cn == nil {
		return true
	}

	if c.BadNode != cn.BadNode {
		return false
	}

	aLines := strings.Split(c.Text, "\n")
	bLines := strings.Split(cn.Text, "\n")

	if len(aLines) != len(bLines) {
		return false
	}

	for i := range aLines {
		if strings.TrimSpace(aLines[i]) != strings.TrimSpace(bLines[i]) {
			return false
		}
	}

	if c.Style != cn.Style {
		return false
	}

	return true
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
	// TODO(jpf): ut
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

	nLine := strings.Count(text, "\n") // "\r\nline 1", this will start with line 1,0 in parsed ast
	if startPos.col == 0 {
		nLine = nLine - 1
	}
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
