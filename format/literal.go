package format

import (
	"bytes"
	"fmt"

	"github.com/joyme123/thrift-ls/parser"
)

func MustFormatLiteral(l *parser.Literal, indent string) string {
	if len(l.Comments) > 0 {
		buf := bytes.NewBuffer(nil)
		buf.WriteString(MustFormatComments(l.Comments, indent))

		if lineDistance(l.Comments[len(l.Comments)-1], l.Value) >= 1 {
			buf.WriteString("\n")
			buf.WriteString(indent)
		} else {
			buf.WriteString(" ")
		}

		buf.WriteString(fmt.Sprintf("%s%s%s", l.Quote, l.Value.Text, l.Quote))

		return buf.String()
	}
	return indent + fmt.Sprintf("%s%s%s", l.Quote, l.Value.Text, l.Quote)
}
