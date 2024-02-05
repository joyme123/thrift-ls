package format

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/joyme123/thrift-ls/parser"
)

type fieldGroup []string

func MustFormatFields(fields []*parser.Field, indent string) string {
	buf := bytes.NewBuffer(nil)

	fmtCtx := &fmtContext{}

	var fieldGroups []fieldGroup
	var fg fieldGroup
	for _, field := range fields {
		if needAddtionalLineForFields(fmtCtx.preNode, field) {
			fieldGroups = append(fieldGroups, fg)
			fg = make(fieldGroup, 0)
		}
		space := " "
		if Align == AlignTypeField {
			space = "\t"
		}
		fg = append(fg, MustFormatField(field, space, indent))
		fmtCtx.preNode = field
	}

	if len(fg) > 0 {
		fieldGroups = append(fieldGroups, fg)
	}

	for i, fg := range fieldGroups {
		w := new(tabwriter.Writer)
		w.Init(buf, 1, 8, 1, ' ', tabwriter.TabIndent)
		for j := range fg {
			fmt.Fprintln(w, fg[j])
		}
		w.Flush()

		if i < len(fieldGroups)-1 {
			buf.WriteString("\n")
		}
	}

	return buf.String()
}

func MustFormatOneLineFields(fields []*parser.Field) string {
	buf := bytes.NewBuffer(nil)
	for i, field := range fields {
		buf.WriteString(MustFormatField(field, " ", ""))
		if i < len(fields)-1 {
			buf.WriteString(" ")
		}
	}

	return buf.String()
}

func MustFormatField(field *parser.Field, space string, indent string) string {
	comments, annos := formatCommentsAndAnnos(field.Comments, field.Annotations, indent)
	if len(field.Comments) > 0 && lineDistance(field.Comments[len(field.Comments)-1], field.Index) > 1 {
		comments = comments + "\n"
	}

	buf := bytes.NewBuffer([]byte(comments))
	required := ""
	if field.RequiredKeyword != nil {
		required = MustFormatKeyword(field.RequiredKeyword.Keyword) + space
	}

	value := ""
	if field.ConstValue != nil {
		equalSpace := space
		if Align == AlignTypeAssign {
			equalSpace = "\t"
		}
		value = fmt.Sprintf("%s%s%s%s", equalSpace, MustFormatKeyword(field.EqualKeyword.Keyword), equalSpace, MustFormatConstValue(field.ConstValue, indent, false))
	}
	str := fmt.Sprintf("%s%d:%s%s%s%s%s%s", indent, field.Index.Value, space, required, MustFormatFieldType(field.FieldType), space, field.Identifier.Name.Text, value)
	buf.WriteString(str)
	buf.WriteString(annos)
	if FieldLineComma == FieldLineCommaAdd {
		buf.WriteString(",")
	} else if FieldLineComma == FieldLineCommaDisable {
		buf.WriteString(formatListSeparator(field.ListSeparatorKeyword))
	}

	if len(field.EndLineComments) > 0 {
		buf.WriteString(MustFormatEndLineComments(field.EndLineComments, ""))
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

	return curStartLine-preNode.End().Line > 1
}
