package format

import (
	"fmt"

	"github.com/joyme123/thrift-ls/parser"
)

func MustFormatIdentifier(id *parser.Identifier) string {
	return fmt.Sprintf("%s %s", MustFormatComments(id.Comments), id.Name)
}
