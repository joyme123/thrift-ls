package format

import (
	"fmt"

	"github.com/joyme123/thrift-ls/parser"
)

func MustFormatTypedef(td *parser.Typedef) string {
	return fmt.Sprintf("typedef %s %s\n", MustFormatFieldType(td.T), td.Alias.Name)
}
