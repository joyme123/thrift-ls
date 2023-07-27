package format

import (
	"bytes"
	"fmt"

	"github.com/joyme123/thrift-ls/parser"
)

func MustFormatEnumValues(values []*parser.EnumValue, indent string) string {
	buf := bytes.NewBuffer(nil)

	for i, v := range values {
		buf.WriteString(indent + MustFormatEnumValue(v))
		if i < len(values)-1 {
			buf.WriteString("\n")
		}
	}

	return buf.String()
}

func MustFormatEnumValue(enumValue *parser.EnumValue) string {
	buf := bytes.NewBufferString(enumValue.Name.Name)
	if enumValue.ValueNode != nil {
		buf.WriteString(fmt.Sprintf(" = %s", MustFormatConstValue(enumValue.ValueNode)))
	}

	return buf.String()
}
