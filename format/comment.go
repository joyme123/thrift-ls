package format

import (
	"bytes"

	"github.com/joyme123/thrift-ls/parser"
)

// TODO(jpf): 多行注释，换行上还需要优化
func MustFormatComments(comments []*parser.Comment) string {
	buf := bytes.NewBuffer(nil)
	for i, c := range comments {
		buf.WriteString(c.Text)
		if i < len(comments)-1 {
			buf.WriteString("\n")
		}
	}

	return buf.String()
}
