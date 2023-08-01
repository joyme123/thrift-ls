package format

import (
	"bytes"
	"fmt"

	"github.com/joyme123/thrift-ls/parser"
)

func MustFormatAnnotations(annotations *parser.Annotations) string {
	buf := bytes.NewBuffer(nil)

	buf.WriteString(MustFormatKeyword(annotations.LParKeyword.Keyword))

	for i, anno := range annotations.Annotations {
		buf.WriteString(MustFormatAnnotation(anno, i == len(annotations.Annotations)-1))
	}

	buf.WriteString(MustFormatKeyword(annotations.RParKeyword.Keyword))

	return buf.String()
}

func MustFormatAnnotation(anno *parser.Annotation, isLast bool) string {
	sep := ""
	if !isLast {
		sep = MustFormatKeyword(anno.ListSeparatorKeyword.Keyword)
	}

	space := ""
	if !isLast {
		space = " "
	}

	// a = "xxxx",
	return fmt.Sprintf("%s %s %s%s%s", MustFormatIdentifier(anno.Identifier), MustFormatKeyword(anno.EqualKeyword.Keyword), MustFormatLiteral(anno.Value), sep, space)
}
