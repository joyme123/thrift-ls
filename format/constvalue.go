package format

import (
	"bytes"
	"fmt"

	"github.com/joyme123/thrift-ls/parser"
)

func MustFormatConstValue(cv *parser.ConstValue, indent string) string {
	buf := bytes.NewBuffer(nil)
	if len(cv.Comments) > 0 {
		buf.WriteString(MustFormatComments(cv.Comments, indent))
	}
	sep := ""
	if cv.ListSeparatorKeyword != nil {
		sep = MustFormatKeyword(cv.ListSeparatorKeyword.Keyword) + " "
	}

	switch cv.TypeName {
	case "list":
		values := cv.Value.([]*parser.ConstValue)

		if len(cv.Comments) > 0 && len(values) > 0 {
			if lineDistance(cv.Comments[len(cv.Comments)-1], values[0]) >= 1 {
				buf.WriteString("\n")
			}
		}

		buf.WriteString(MustFormatKeyword(cv.LBrkKeyword.Keyword))
		for i := range values {
			buf.WriteString(MustFormatConstValue(values[i], ""))
		}
		buf.WriteString(MustFormatKeyword(cv.RBrkKeyword.Keyword))
	case "map":
		values := cv.Value.([]*parser.ConstValue)

		if len(cv.Comments) > 0 && len(values) > 0 {
			if lineDistance(cv.Comments[len(cv.Comments)-1], values[0]) >= 1 {
				buf.WriteString("\n")
			}
		}

		var preNode parser.Node
		buf.WriteString(MustFormatKeyword(cv.LCurKeyword.Keyword))
		preNode = cv.LCurKeyword
		for i := range values {
			distance := lineDistance(preNode, values[i])
			if distance >= 1 {
				buf.WriteString("\n")
			}
			buf.WriteString(MustFormatConstValue(values[i], ""))
			preNode = values[i]
		}
		if lineDistance(cv.RCurKeyword, preNode) >= 1 {
			buf.WriteString("\n")
		}
		buf.WriteString(MustFormatKeyword(cv.RCurKeyword.Keyword))
	case "pair":
		key := cv.Key.(*parser.ConstValue)
		value := cv.Value.(*parser.ConstValue)

		if len(cv.Comments) > 0 {
			if lineDistance(cv.Comments[len(cv.Comments)-1], key) >= 1 {
				buf.WriteString("\n")
			}
		}

		if cv.ListSeparatorKeyword != nil {
			sep = MustFormatKeyword(cv.ListSeparatorKeyword.Keyword)
		}
		buf.WriteString(fmt.Sprintf("%s%s %s%s", MustFormatConstValue(key, indent+Indent), MustFormatKeyword(cv.ColonKeyword.Keyword), MustFormatConstValue(value, ""), sep))
	case "identifier":
		if len(cv.Comments) > 0 {
			// special case for iline distance
			if lineDistance(cv.Comments[len(cv.Comments)-1], cv) >= 1 {
				buf.WriteString("\n")
			}
		}
		buf.WriteString(indent + fmt.Sprintf("%s%s", cv.Value.(string), sep))
	case "string":
		val := ""
		if _, ok := cv.Value.(string); ok {
			val = cv.Value.(string)
			buf.WriteString(fmt.Sprintf("%q%s", val, sep))
		} else {
			literal := cv.Value.(*parser.Literal)
			if len(cv.Comments) > 0 {
				if lineDistance(cv.Comments[len(cv.Comments)-1], literal) >= 1 {
					buf.WriteString("\n")
				}
			}
			val = MustFormatLiteral(literal)
			buf.WriteString(fmt.Sprintf("%s%s", val, sep))
		}
	case "i64":
		buf.WriteString(fmt.Sprintf("%s%s", cv.ValueInText, sep))
	case "double":
		buf.WriteString(fmt.Sprintf("%s%s", cv.ValueInText, sep))
	}

	return buf.String()
}
