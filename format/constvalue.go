package format

import (
	"bytes"
	"fmt"

	"github.com/joyme123/thrift-ls/parser"
)

func MustFormatConstValue(cv *parser.ConstValue) string {
	switch cv.TypeName {
	case "list":
		values := cv.Value.([]*parser.ConstValue)
		buf := bytes.NewBufferString("[")
		for i := range values {
			if i != 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(MustFormatConstValue(values[i]))
		}
		buf.WriteByte(']')
		return buf.String()
	case "map":
		values := cv.Value.([]*parser.ConstValue)
		buf := bytes.NewBufferString("{")
		for i := range values {
			if i != 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(MustFormatConstValue(values[i]))
		}
		buf.WriteByte('}')
		return buf.String()
	case "pair":
		key := cv.Key.(*parser.ConstValue)
		value := cv.Value.(*parser.ConstValue)
		return fmt.Sprintf("%s: %s", MustFormatConstValue(key), MustFormatConstValue(value))
	case "identifier":
		return cv.Value.(string)
	case "string":
		return fmt.Sprintf("%q", cv.Value)
	case "i64":
		return cv.ValueInText
	case "double":
		return cv.ValueInText
	}

	return ""
}
