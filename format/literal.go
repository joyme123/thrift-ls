package format

import (
	"fmt"

	"github.com/joyme123/thrift-ls/parser"
)

func MustFormatLiteral(l *parser.Literal) string {
	if len(l.Comments) > 0 {
		return fmt.Sprintf("%s %s%s%s", MustFormatComments(l.Comments, ""), l.Quote, l.Value.Text, l.Quote)
	}
	return fmt.Sprintf("%s%s%s", l.Quote, l.Value.Text, l.Quote)
}
