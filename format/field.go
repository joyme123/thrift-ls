package format

import (
	"fmt"

	"github.com/joyme123/thrift-ls/parser"
)

func MustFormatField(field *parser.Field) string {
	required := ""
	if field.Required != nil {
		if !field.Required.Required {
			required = " optional"
		} else {
			required = " required"
		}
	}

	value := ""
	if field.ConstValue != nil {
		value = fmt.Sprintf(" = %s", MustFormatConstValue(field.ConstValue))
	}
	return fmt.Sprintf("%d:%s %s %s%s", field.Index.Value, required, MustFormatFieldType(field.FieldType), field.Identifier.Name, value)
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
