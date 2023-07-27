package format

import (
	"bytes"
	"fmt"

	"github.com/joyme123/thrift-ls/parser"
)

func MustFormatFunction(fn *parser.Function) string {
	fnType := "void"
	if fn.Void == nil {
		fnType = MustFormatFieldType(fn.FunctionType)
	}
	buf := bytes.NewBufferString(fmt.Sprintf("%s %s(", fnType, fn.Name.Name))

	for i := range fn.Arguments {
		if i != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(MustFormatField(fn.Arguments[i]))
	}
	buf.WriteString(")")

	if fn.Throws != nil {
		buf.WriteString(" throws (")
		for i := range fn.Throws.Fields {
			if i != 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(MustFormatField(fn.Throws.Fields[i]))
		}
		buf.WriteString(")")
	}

	return buf.String()
}
