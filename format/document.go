package format

import (
	"bytes"

	"github.com/joyme123/thrift-ls/parser"
)

func FormatDocument(doc *parser.Document) (string, error) {
	if doc.ChildrenBadNode() {
		return "", BadNodeError
	}

	buf := bytes.NewBuffer(nil)

	var headers []parser.Node
	var definitions []parser.Node
	for i := range doc.Nodes {
		node := doc.Nodes[i]
		switch node.Type() {
		case "Include", "CPPInclude", "Namespace":
			headers = append(headers, node)
		default:
			definitions = append(definitions, node)
		}
	}

	writeBuf := func(node parser.Node) {
		switch node.Type() {
		case "Include":
			buf.WriteString(MustFormatInclude(node.(*parser.Include)))
			buf.WriteString("\n")
		case "CPPInclude":
			buf.WriteString(MustFormatCPPInclude(node.(*parser.CPPInclude)))
			buf.WriteString("\n")
		case "Namespace":
			buf.WriteString(MustFormatNamespace(node.(*parser.Namespace)))
			buf.WriteString("\n")
		case "Struct":
			buf.WriteString(MustFormatStruct(node.(*parser.Struct)))
			buf.WriteString("\n")
		case "Union":
			buf.WriteString(MustFormatUnion(node.(*parser.Union)))
			buf.WriteString("\n")
		case "Exception":
			buf.WriteString(MustFormatException(node.(*parser.Exception)))
			buf.WriteString("\n")
		case "Service":
			buf.WriteString(MustFormatService(node.(*parser.Service)))
			buf.WriteString("\n")
		case "Typedef":
			buf.WriteString(MustFormatTypedef(node.(*parser.Typedef)))
			buf.WriteString("\n")
		case "Const":
			buf.WriteString(MustFormatConst(node.(*parser.Const)))
			buf.WriteString("\n")
		case "Enum":
			buf.WriteString(MustFormatEnum(node.(*parser.Enum)))
			buf.WriteString("\n")
		}
	}

	for _, node := range headers {
		writeBuf(node)
	}
	// addition line between headers and definitions
	buf.WriteString("\n")
	for _, node := range definitions {
		writeBuf(node)
	}

	return buf.String(), nil
}
