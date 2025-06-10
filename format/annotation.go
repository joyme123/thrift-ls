package format

import (
	"bytes"
	"fmt"

	"github.com/joyme123/thrift-ls/parser"
)

func MustFormatAnnotations(annotations *parser.Annotations) string {
	buf := bytes.NewBuffer(nil)

	buf.WriteString(MustFormatKeyword(annotations.LParKeyword.Keyword))

	var preNode parser.Node
	preNode = annotations.LParKeyword

	indent := ""
	isNewLine := false

	for i, anno := range annotations.Annotations {
		if lineDistance(preNode, annotations.Annotations[i]) >= 1 {
			buf.WriteString("\n")
			isNewLine = true
			indent = Indent + Indent
		}
		buf.WriteString(MustFormatAnnotation(anno, i == len(annotations.Annotations)-1, i == 0, indent, isNewLine))
		preNode = annotations.Annotations[i]
		isNewLine = false
		indent = ""
	}

	if lineDistance(preNode, annotations.RParKeyword) >= 1 {
		buf.WriteString("\n")
		buf.WriteString(Indent)
	}
	buf.WriteString(MustFormatKeyword(annotations.RParKeyword.Keyword))

	return buf.String()
}

func MustFormatAnnotation(anno *parser.Annotation, isLast bool, isFirst bool, indent string, isNewLine bool) string {
	sep := ""
	if (!isLast) && anno.ListSeparatorKeyword != nil {
		sep = MustFormatKeyword(anno.ListSeparatorKeyword.Keyword)
	}

	space := ""
	if (!isFirst) && (!isNewLine) {
		space = " "
	}

	// a = "xxxx",
	return fmt.Sprintf("%s%s %s %s%s", space, MustFormatIdentifier(anno.Identifier, indent), MustFormatKeyword(anno.EqualKeyword.Keyword), MustFormatLiteral(anno.Value, ""), sep)
}
