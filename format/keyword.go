package format

import (
	"fmt"

	"github.com/joyme123/thrift-ls/parser"
)

func MustFormatKeyword(kw parser.Keyword) string {
	if len(kw.Comments) > 0 {
		return fmt.Sprintf("%s %s", MustFormatComments(kw.Comments, ""), kw.Literal.Text)
	}

	return kw.Literal.Text
}
