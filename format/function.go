package format

import (
	"bytes"
	"fmt"

	"github.com/joyme123/thrift-ls/parser"
)

func MustFormatFunction(fn *parser.Function) string {
	fnType := "void"
	if !fn.Void {
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

	if len(fn.Throws) > 0 {
		buf.WriteString(" throws (")
		for i := range fn.Throws {
			if i != 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(MustFormatField(fn.Throws[i]))
		}
		buf.WriteString(")")
	}

	return buf.String()
}
