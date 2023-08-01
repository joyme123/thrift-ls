package format

import (
	"bytes"
	"fmt"

	"github.com/joyme123/thrift-ls/parser"
)

func MustFormatEnumValues(values []*parser.EnumValue, indent string) string {
	buf := bytes.NewBuffer(nil)

	fmtCtx := &fmtContext{}

	for i, v := range values {
		if needAddtionalLineForEnumValues(fmtCtx.preNode, values[i]) {
			buf.WriteString("\n")
		}
		buf.WriteString(MustFormatEnumValue(v, indent))
		if i < len(values)-1 {
			buf.WriteString("\n")
		}
		fmtCtx.preNode = values[i]
	}

	return buf.String()
}

// TODO(jpf): comments
func MustFormatEnumValue(enumValue *parser.EnumValue, indent string) string {
	comments, annos := formatCommentsAndAnnos(enumValue.Comments, enumValue.Annotations, indent)

	if len(comments) > 0 && lineDistance(enumValue.Comments[len(enumValue.Comments)-1], enumValue.Name) > 1 {
		comments = comments + "\n"
	}

	buf := bytes.NewBufferString(comments)
	buf.WriteString(indent + MustFormatIdentifier(enumValue.Name))
	if enumValue.ValueNode != nil {
		buf.WriteString(fmt.Sprintf(" %s %s", MustFormatKeyword(enumValue.EqualKeyword.Keyword), MustFormatConstValue(enumValue.ValueNode)))
	}

	buf.WriteString(annos)

	if enumValue.ListSeparatorKeyword != nil {
		buf.WriteString(MustFormatKeyword(enumValue.ListSeparatorKeyword.Keyword))
	}

	buf.WriteString(MustFormatEndLineComments(enumValue.EndLineComments, ""))

	return buf.String()
}

func needAddtionalLineForEnumValues(preNode, curNode parser.Node) bool {
	if preNode == nil {
		return false
	}

	curValue := curNode.(*parser.EnumValue)

	var curStartLine int
	if len(curValue.Comments) > 0 {
		curStartLine = curValue.Comments[0].Pos().Line
	} else {
		curStartLine = curValue.Name.Pos().Line
	}

	return curStartLine-preNode.End().Line > 1
}
