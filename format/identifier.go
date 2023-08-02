package format

import (
	"fmt"

	"github.com/joyme123/thrift-ls/parser"
)

func MustFormatIdentifier(id *parser.Identifier) string {
	comments := MustFormatComments(id.Comments, "")
	if comments != "" {
		comments = comments + " "
	}
	return fmt.Sprintf("%s%s", comments, id.Name.Text)
}
