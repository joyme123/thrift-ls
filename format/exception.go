package format

import (
	"bytes"
	"fmt"

	"github.com/joyme123/thrift-ls/parser"
)

func MustFormatException(excep *parser.Exception) string {
	buf := bytes.NewBufferString(fmt.Sprintf("exception %s {\n", excep.Name.Name))
	for i := range excep.Fields {
		buf.WriteString(fmt.Sprintf("  %s\n", MustFormatField(excep.Fields[i])))
	}

	buf.WriteString("}\n")

	return buf.String()
}
