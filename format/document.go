package format

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/joyme123/thrift-ls/parser"
)

type fmtContext struct {
	// preNode record previous print node. we can use preNode as print context
	// if preNodex is const or typdef, and current node is const or typedef, '\n' should be ignore
	preNode parser.Node
}

func FormatDocument(doc *parser.Document) (string, error) {
	return FormatDocumentWithValidation(doc, false)
}

func FormatDocumentWithValidation(doc *parser.Document, selfValidation bool) (string, error) {
	if doc.ChildrenBadNode() {
		return "", BadNodeError
	}

	buf := bytes.NewBuffer(nil)

	fmtCtx := &fmtContext{}

	writeBuf := func(node parser.Node, addtionalLine bool) {
		if addtionalLine {
			if len(buf.Bytes()) > 0 && buf.Bytes()[buf.Len()-1] != '\n' {
				// if preNode doesn't have \n at end of line, set \n for it
				buf.WriteString("\n")
			}
			buf.WriteString("\n")
		}

		switch node.Type() {
		case "Include":
			buf.WriteString(MustFormatInclude(node.(*parser.Include)))
		case "CPPInclude":
			buf.WriteString(MustFormatCPPInclude(node.(*parser.CPPInclude)))
		case "Namespace":
			buf.WriteString(MustFormatNamespace(node.(*parser.Namespace)))
		case "Struct":
			buf.WriteString(MustFormatStruct(node.(*parser.Struct)))
		case "Union":
			buf.WriteString(MustFormatUnion(node.(*parser.Union)))
		case "Exception":
			buf.WriteString(MustFormatException(node.(*parser.Exception)))
		case "Service":
			buf.WriteString(MustFormatService(node.(*parser.Service)))
		case "Typedef":
			buf.WriteString(MustFormatTypedef(node.(*parser.Typedef)))
		case "Const":
			buf.WriteString(MustFormatConst(node.(*parser.Const)))
		case "Enum":
			buf.WriteString(MustFormatEnum(node.(*parser.Enum)))
		}

	}

	for _, node := range doc.Nodes {
		addtionalLine := needAddtionalLineInDocument(fmtCtx.preNode, node)
		writeBuf(node, addtionalLine)
		fmtCtx.preNode = node
	}

	if len(doc.Comments) > 0 {
		buf.WriteString(MustFormatComments(doc.Comments, ""))
	}
	res := buf.String()

	res = strings.TrimSpace(res)

	if selfValidation {
		psr := parser.PEGParser{}
		formattedAst, err := psr.Parse("formated.thrift", []byte(res))
		if err != nil {
			return "", fmt.Errorf("format error: format result failed to parse, error msg: %v. Please report bug to author at https://github.com/joyme123/thrift-ls/issues", err)
		}

		if !doc.Equals(formattedAst) {
			return "", fmt.Errorf("format error: format result failed to pass self validation. Please report bug to author at https://github.com/joyme123/thrift-ls/issues")
		}
	}

	return res, nil
}

var (
	header = map[string]struct{}{
		"Include":    {},
		"CPPInclude": {},
		"Namespace":  {},
	}
	onelineDefinition = map[string]struct{}{
		"Const":   {},
		"Typedef": {},
	}
	multiLineDefinition = map[string]struct{}{
		"Struct":    {},
		"Union":     {},
		"Exception": {},
		"Service":   {},
		"Typedef":   {},
		"Const":     {},
		"Enum":      {},
	}
)

func isHeader(node parser.Node) bool {
	_, ok := header[node.Type()]
	return ok
}

func isOneLineDefinition(node parser.Node) bool {
	_, ok := onelineDefinition[node.Type()]
	return ok
}

func isMultiLineDefinition(node parser.Node) bool {
	_, ok := multiLineDefinition[node.Type()]
	return ok
}

func needAddtionalLineInDocument(preNode parser.Node, currentNode parser.Node) bool {
	if preNode == nil {
		return false
	}

	if isHeader(preNode) && isHeader(currentNode) {
		if preNode.Type() == currentNode.Type() {
			if lineDistance(preNode, currentNode) > 1 {
				return true
			}
			return false
		}
		return true
	}

	if isOneLineDefinition(preNode) && isOneLineDefinition(currentNode) {
		// if preNode and currentNode has one or more empty lines between them, we should reserve
		// one empty line
		if lineDistance(preNode, currentNode) > 1 {
			return true
		}
		return false
	}

	return true
}
