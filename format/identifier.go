package format

import (
	"fmt"

	"github.com/joyme123/thrift-ls/parser"
)

func MustFormatIdentifier(id *parser.Identifier, indent string) string {
	comments := MustFormatComments(id.Comments, indent)
	if comments != "" {
		comments = comments
		if lineDistance(id.Comments[len(id.Comments)-1], id.Name) >= 1 {
			comments = comments + "\n"
		} else {
			comments = comments + " "
		}
	}
	return fmt.Sprintf("%s%s", comments, indent+id.Name.Text)
}
