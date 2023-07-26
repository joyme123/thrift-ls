package format

import (
	"fmt"

	"github.com/joyme123/thrift-ls/parser"
)

func MustFormatInclude(inc *parser.Include) string {
	return fmt.Sprintf(`include "%s"\n`, inc.Path.Value)
}

func MustFormatCPPInclude(inc *parser.CPPInclude) string {
	return fmt.Sprintf(`cppinclude "%s"\n`, inc.Path.Value)
}
