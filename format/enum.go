package format

import (
	"bytes"
	"fmt"

	"github.com/joyme123/thrift-ls/parser"
)

func MustFormatEnum(enum *parser.Enum) string {
	buf := bytes.NewBufferString(fmt.Sprintf("enum %s {\n", enum.Name.Name))
	for i := range enum.Values {
		buf.WriteString(fmt.Sprintf("  %s\n", MustFormatEnumValue(enum.Values[i])))
	}

	buf.WriteString("}\n")

	return buf.String()
}

func MustFormatEnumValue(enumValue *parser.EnumValue) string {
	buf := bytes.NewBufferString(enumValue.Name.Name)
	if enumValue.ValueNode != nil {
		buf.WriteString(fmt.Sprintf(" = %s", MustFormatConstValue(enumValue.ValueNode)))
	}

	return buf.String()
}
