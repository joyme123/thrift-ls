package format

import (
	"bytes"

	"github.com/joyme123/thrift-ls/parser"
)

func MustFormatKeyword(kw parser.Keyword) string {
	if len(kw.Comments) > 0 {
		buf := bytes.NewBuffer(nil)
		buf.WriteString(MustFormatComments(kw.Comments, ""))

		if lineDistance(kw.Comments[len(kw.Comments)-1], kw.Literal) >= 1 {
			buf.WriteString("\n")
		} else {
			buf.WriteString(" ")
		}

		buf.WriteString(kw.Literal.Text)

		return buf.String()
	}

	return kw.Literal.Text
}
