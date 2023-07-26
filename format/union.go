package format

import (
	"bytes"
	"fmt"

	"github.com/joyme123/thrift-ls/parser"
)

func MustFormatUnion(union *parser.Union) string {
	buf := bytes.NewBufferString(fmt.Sprintf("union %s {\n", union.Name.Name))
	for i := range union.Fields {
		buf.WriteString(fmt.Sprintf("  %s\n", MustFormatField(union.Fields[i])))
	}

	buf.WriteString("}\n")

	return buf.String()
}
