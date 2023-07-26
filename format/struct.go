package format

import (
	"bytes"
	"fmt"

	"github.com/joyme123/thrift-ls/parser"
)

func MustFormatStruct(st *parser.Struct) string {
	buf := bytes.NewBufferString(fmt.Sprintf("struct %s {\n", st.Identifier.Name))
	for i := range st.Fields {
		buf.WriteString(fmt.Sprintf("  %s\n", MustFormatField(st.Fields[i])))
	}

	buf.WriteString("}\n")

	return buf.String()
}
