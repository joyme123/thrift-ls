package format

import (
	"bytes"
	"fmt"

	"github.com/joyme123/thrift-ls/parser"
)

func MustFormatConstValue(cv *parser.ConstValue, indent string, newLine bool) string {
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
			// TODO(jpf): 优化显示
			newLine = false
			buf.WriteString(MustFormatConstValue(values[i], indent, newLine))
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
				newLine = true
			} else {
				buf.WriteString(" ")
				newLine = false
			}
			buf.WriteString(MustFormatConstValue(values[i], indent, newLine))
			preNode = values[i]
		}
		if lineDistance(preNode, cv.RCurKeyword) >= 1 {
			buf.WriteString("\n")
			buf.WriteString(Indent)
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
		buf.WriteString(fmt.Sprintf("%s%s %s%s",
			MustFormatConstValue(key, indent+Indent, newLine),
			MustFormatKeyword(cv.ColonKeyword.Keyword),
			MustFormatConstValue(value, indent, false),
			sep))
	case "identifier":
		if len(cv.Comments) > 0 {
			// special case for iline distance
			if lineDistance(cv.Comments[len(cv.Comments)-1], cv) >= 1 {
				buf.WriteString("\n")
				newLine = true
			}
		}

		if newLine {
			buf.WriteString(indent)
		}

		buf.WriteString(fmt.Sprintf("%s%s", cv.Value.(string), sep))
	case "string":
		val := ""
		if _, ok := cv.Value.(string); ok {
			if len(cv.Comments) > 0 {
				if lineDistance(cv.Comments[len(cv.Comments)-1], cv) >= 1 {
					buf.WriteString("\n")
					newLine = true
				}
			}

			if newLine {
				buf.WriteString(indent)
			}
			val = cv.Value.(string)
			buf.WriteString(fmt.Sprintf("%q%s", val, sep))
		} else {
			literal := cv.Value.(*parser.Literal)
			if len(cv.Comments) > 0 {
				if lineDistance(cv.Comments[len(cv.Comments)-1], literal) >= 1 {
					buf.WriteString("\n")
					newLine = true
				}
			}

			if !newLine {
				indent = ""
			}

			val = MustFormatLiteral(literal, indent)
			buf.WriteString(fmt.Sprintf("%s%s", val, sep))
		}
	case "i64":
		if len(cv.Comments) > 0 {
			if lineDistance(cv.Comments[len(cv.Comments)-1], cv) >= 1 {
				buf.WriteString("\n")
				newLine = true
			}
		}

		if newLine {
			buf.WriteString(indent)
		}
		buf.WriteString(fmt.Sprintf("%s%s", cv.ValueInText, sep))
	case "double":
		if len(cv.Comments) > 0 {
			if lineDistance(cv.Comments[len(cv.Comments)-1], cv) >= 1 {
				buf.WriteString("\n")
				newLine = true
			}
		}

		if newLine {
			buf.WriteString(indent)
		}
		buf.WriteString(fmt.Sprintf("%s%s", cv.ValueInText, sep))
	}

	return buf.String()
}
