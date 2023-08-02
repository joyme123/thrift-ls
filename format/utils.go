package format

import (
	"bytes"
	"fmt"
	"text/template"
	"unicode"

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
		annos = " " + MustFormatAnnotations(annotations)
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

// FormatedEquals is used to judge if documents has been changed after formated
// implementation: ignore all space charactor to compare two strings
func EqualsAfterFormat(doc1, doc2 string) error {
	cur1, cur2 := 0, 0
	runes1 := []rune(doc1)
	runes2 := []rune(doc2)

	for cur1 < len(runes1) && cur2 < len(runes2) {
		for cur1 < len(runes1) && unicode.IsSpace(runes1[cur1]) {
			cur1++
		}
		for cur2 < len(runes2) && unicode.IsSpace(runes2[cur2]) {
			cur2++
		}

		if cur1 >= len(runes1) || cur2 >= len(runes2) {
			break
		}

		if runes1[cur1] != runes2[cur2] {
			return fmt.Errorf("different at doc1: %s as %d, doc2: %s at %d, str1: %s, str2: %s", string(runes1[cur1]), cur1, string(runes2[cur2]), cur2, showStringContext(runes1, cur1, 40), showStringContext(runes2, cur2, 40))
		}

		cur1++
		cur2++
	}

	for cur1 < len(runes1) {
		for !unicode.IsSpace(runes1[cur1]) {
			return fmt.Errorf("")
		}
		cur1++
	}

	for cur2 < len(runes2) {
		for !unicode.IsSpace(runes2[cur2]) {
			return fmt.Errorf("")
		}
		cur2++
	}

	return nil
}

func showStringContext(text []rune, offset int, n int) string {
	start := offset - n
	if start < 0 {
		start = 0
	}

	end := offset + n
	if end >= len(text) {
		end = len(text) - 1
	}

	return string(text[start : end+1])
}
