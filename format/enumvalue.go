package format

import (
	"bytes"
	"fmt"
	"text/tabwriter"

	"github.com/joyme123/thrift-ls/parser"
)

type enumValueGroup []string

func MustFormatEnumValues(values []*parser.EnumValue, indent string) string {
	buf := bytes.NewBuffer(nil)

	fmtCtx := &fmtContext{}

	var enumValueGroups []enumValueGroup
	var eg enumValueGroup

	for i, v := range values {
		if needAddtionalLineForEnumValues(fmtCtx.preNode, values[i]) {
			enumValueGroups = append(enumValueGroups, eg)
			eg = make(enumValueGroup, 0)
		}
		space := " "
		if Align == AlignTypeField {
			space = "\t"
		}
		eg = append(eg, MustFormatEnumValue(v, space, indent))
		fmtCtx.preNode = values[i]
	}

	if len(eg) > 0 {
		enumValueGroups = append(enumValueGroups, eg)
	}

	for i, eg := range enumValueGroups {
		w := new(tabwriter.Writer)
		w.Init(buf, 1, 8, 1, ' ', tabwriter.TabIndent)
		for j := range eg {
			fmt.Fprintln(w, eg[j])
		}
		w.Flush()

		if i < len(enumValueGroups)-1 {
			buf.WriteString("\n")
		}
	}

	return buf.String()
}

func MustFormatEnumValue(enumValue *parser.EnumValue, space, indent string) string {
	comments, annos := formatCommentsAndAnnos(enumValue.Comments, enumValue.Annotations, indent)

	if len(comments) > 0 && lineDistance(enumValue.Comments[len(enumValue.Comments)-1], enumValue.Name) > 1 {
		comments = comments + "\n"
	}

	buf := bytes.NewBufferString(comments)
	buf.WriteString(indent + MustFormatIdentifier(enumValue.Name, ""))
	if enumValue.ValueNode != nil {
		equalSpace := space
		if Align == AlignTypeAssign {
			equalSpace = "\t"
		}
		buf.WriteString(fmt.Sprintf("%s%s%s%s", equalSpace, MustFormatKeyword(enumValue.EqualKeyword.Keyword), equalSpace, MustFormatConstValue(enumValue.ValueNode, indent, false)))
	}

	buf.WriteString(annos)

	if FieldLineComma == FieldLineCommaAdd {
		buf.WriteString(",")
	} else if FieldLineComma == FieldLineCommaDisable {
		if enumValue.ListSeparatorKeyword != nil {
			buf.WriteString(MustFormatKeyword(enumValue.ListSeparatorKeyword.Keyword))
		}
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
