package format

import (
	"bytes"
	"text/template"

	"github.com/joyme123/thrift-ls/parser"
)

const (
	Indent = "    "
)

func MustFormat(tplText string, formatter any) string {
	tpl, err := template.New("default").Parse(tplText)
	if err != nil {
		panic(err)
	}

	buf := bytes.NewBuffer(nil)
	err = tpl.Execute(buf, formatter)
	if err != nil {
		panic(err)
	}

	return buf.String()
}

func formatCommentsAndAnnos(comments []*parser.Comment, annotations *parser.Annotations, indent string) (string, string) {
	commentsStr := ""
	if len(comments) > 0 {
		commentsStr = MustFormatComments(comments, indent) + "\n"
	}
	annos := ""
	if annotations != nil && len(annotations.Annotations) > 0 {
		annos = " " + MustFormatAnnotations(annotations) + " "
	}

	return commentsStr, annos
}

func formatListSeparator(sep *parser.ListSeparatorKeyword) string {
	if sep == nil {
		return ""
	}

	return MustFormatKeyword(sep.Keyword)
}

func lineDistance(preNode parser.Node, currentNode parser.Node) int {
	return currentNode.Pos().Line - preNode.End().Line
}
