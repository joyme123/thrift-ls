package format

import (
	"bytes"
	"strings"

	"github.com/joyme123/thrift-ls/parser"
)

// TODO(jpf): 多行注释，换行上还需要优化
// MustFormatComments formats comments
// return string doesn't include '\n' end of line
func MustFormatComments(comments []*parser.Comment, indent string) string {
	fmtCtx := &fmtContext{}
	buf := bytes.NewBuffer(nil)
	for _, c := range comments {
		if fmtCtx.preNode != nil {
			buf.WriteString("\n")
			if lineDistance(fmtCtx.preNode, c) > 1 {
				buf.WriteString("\n")
			}
		}

		buf.WriteString(formatMultiLineComment(c.Text, indent))

		fmtCtx.preNode = c
	}

	return buf.String()
}

func formatMultiLineComment(comment string, indent string) string {
	comment = strings.TrimSpace(comment)

	if strings.HasPrefix(comment, "//") {
		return indent + comment
	}

	lines := strings.Split(comment, "\n")
	if len(lines) == 1 {
		return indent + comment
	}

	buf := bytes.NewBuffer(nil)
	for i, line := range lines {
		line = strings.TrimSpace(line)
		space := ""
		if strings.HasPrefix(line, "*") {
			space = " "
		}
		if i == 0 {
			buf.WriteString(indent + space + line)
		} else {
			buf.WriteString("\n" + indent + space + line)
		}
	}

	return buf.String()

}
