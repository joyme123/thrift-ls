package format

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/joyme123/thrift-ls/parser"
)

func MustFormatFields(fields []*parser.Field, indent string) string {
	buf := bytes.NewBuffer(nil)

	fmtCtx := &fmtContext{}

	for i, field := range fields {
		if needAddtionalLineForFields(fmtCtx.preNode, field) {
			buf.WriteString("\n")
		}
		buf.WriteString(MustFormatField(field, indent))
		if i < len(fields)-1 {
			buf.WriteString("\n")
		}
		fmtCtx.preNode = field
	}

	return buf.String()
}

func MustFormatOneLineFields(fields []*parser.Field) string {
	buf := bytes.NewBuffer(nil)
	for i, field := range fields {
		buf.WriteString(MustFormatField(field, ""))
		if i < len(fields)-1 {
			buf.WriteString(" ")
		}
	}

	return buf.String()
}

func MustFormatField(field *parser.Field, indent string) string {
	comments, annos := formatCommentsAndAnnos(field.Comments, field.Annotations, indent)
	if len(field.Comments) > 0 && lineDistance(field.Comments[len(field.Comments)-1], field.Index) > 1 {
		comments = comments + "\n"
	}

	buf := bytes.NewBuffer([]byte(comments))
	required := ""
	if field.RequiredKeyword != nil {
		required = MustFormatKeyword(field.RequiredKeyword.Keyword) + " "
	}

	value := ""
	if field.ConstValue != nil {
		value = fmt.Sprintf(" %s %s", MustFormatKeyword(field.EqualKeyword.Keyword), MustFormatConstValue(field.ConstValue))
	}
	str := fmt.Sprintf("%s%d: %s%s %s%s", indent, field.Index.Value, required, MustFormatFieldType(field.FieldType), field.Identifier.Name, value)
	buf.WriteString(str)
	buf.WriteString(annos)
	buf.WriteString(formatListSeparator(field.ListSeparatorKeyword))
	if len(field.EndLineComments) > 0 {
		buf.WriteString(" " + MustFormatComments(field.EndLineComments, ""))
	}

	// remove space at end of line
	return strings.TrimRight(buf.String(), " ")
}

func MustFormatFieldType(ft *parser.FieldType) string {
	annos := ""
	if ft.Annotations != nil {
		annos = MustFormatAnnotations(ft.Annotations)
		if len(ft.Annotations.Annotations) > 0 {
			annos = " " + annos
		}
	}

	tn := MustFormatTypeName(ft.TypeName)

	switch ft.TypeName.Name {
	case "map":
		return fmt.Sprintf("%s<%s,%s>%s", tn, MustFormatFieldType(ft.KeyType), MustFormatFieldType(ft.ValueType), annos)
	case "set":
		return fmt.Sprintf("%s<%s>%s", tn, MustFormatFieldType(ft.KeyType), annos)
	case "list":
		return fmt.Sprintf("%s<%s>%s", tn, MustFormatFieldType(ft.KeyType), annos)
	default:
		return tn + annos
	}
}

func MustFormatTypeName(tn *parser.TypeName) string {
	comments := MustFormatComments(tn.Comments, "")
	if len(tn.Comments) > 0 {
		comments = comments + " "
	}

	return comments + tn.Name
}

func needAddtionalLineForFields(preNode, curNode parser.Node) bool {
	if preNode == nil {
		return false
	}

	curField := curNode.(*parser.Field)

	var curStartLine int
	if len(curField.Comments) > 0 {
		curStartLine = curField.Comments[0].Pos().Line
	} else {
		if curField.Index != nil {
			curStartLine = curField.Index.Pos().Line
		} else if curField.RequiredKeyword != nil {
			curStartLine = curField.RequiredKeyword.Pos().Line
		} else {
			curStartLine = curField.FieldType.Pos().Line
		}
	}

	return curStartLine-preNode.End().Line >= 1
}
