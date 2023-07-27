package format

import (
	"bytes"
	"fmt"

	"github.com/joyme123/thrift-ls/parser"
)

func MustFormatFields(fields []*parser.Field, indent string) string {
	buf := bytes.NewBuffer(nil)
	for i, field := range fields {
		buf.WriteString(indent + MustFormatField(field))
		if i < len(fields)-1 {
			buf.WriteString("\n")
		}
	}

	return buf.String()
}

func MustFormatOneLineFields(fields []*parser.Field) string {
	buf := bytes.NewBuffer(nil)
	for i, field := range fields {
		buf.WriteString(MustFormatField(field))
		if i < len(fields)-1 {
			buf.WriteString(", ")
		}
	}

	return buf.String()
}

func MustFormatField(field *parser.Field) string {
	comments, annos := formatCommentsAndAnnos(field.Comments, field.Annotations)

	buf := bytes.NewBuffer([]byte(comments))
	required := ""
	if field.RequiredKeyword != nil {
		required = MustFormatKeyword(field.RequiredKeyword.Keyword) + " "
	}

	value := ""
	if field.ConstValue != nil {
		value = fmt.Sprintf(" %s %s", MustFormatKeyword(field.EqualKeyword.Keyword), MustFormatConstValue(field.ConstValue))
	}
	str := fmt.Sprintf("%d: %s%s %s%s", field.Index.Value, required, MustFormatFieldType(field.FieldType), field.Identifier.Name, value)
	buf.WriteString(str)
	buf.WriteString(annos)
	buf.WriteString(formatListSeparator(field.ListSeparatorKeyword))
	if len(field.EndLineComments) > 0 {
		buf.WriteString(" " + MustFormatComments(field.EndLineComments))
	}

	return buf.String()
}

func MustFormatFieldType(ft *parser.FieldType) string {
	switch ft.TypeName.Name {
	case "map":
		return fmt.Sprintf("map<%s,%s>", MustFormatFieldType(ft.KeyType), MustFormatFieldType(ft.ValueType))
	case "set":
		return fmt.Sprintf("set<%s>", MustFormatFieldType(ft.KeyType))
	case "list":
		return fmt.Sprintf("list<%s>", MustFormatFieldType(ft.KeyType))
	default:
		return ft.TypeName.Name
	}
}
