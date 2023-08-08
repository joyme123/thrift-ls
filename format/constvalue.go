package format

import (
	"bytes"
	"fmt"

	"github.com/joyme123/thrift-ls/parser"
)

func MustFormatConstValue(cv *parser.ConstValue) string {
	buf := bytes.NewBuffer(nil)
	if cv.Comments != nil {
		buf.WriteString(MustFormatComments(cv.Comments, Indent))
	}
	sep := ""
	if cv.ListSeparatorKeyword != nil {
		sep = MustFormatKeyword(cv.ListSeparatorKeyword.Keyword) + " "
	}

	switch cv.TypeName {
	case "list":
		values := cv.Value.([]*parser.ConstValue)
		buf.WriteString(MustFormatKeyword(cv.LBrkKeyword.Keyword))
		for i := range values {
			buf.WriteString(MustFormatConstValue(values[i]))
		}
		buf.WriteString(MustFormatKeyword(cv.RBrkKeyword.Keyword))
	case "map":
		values := cv.Value.([]*parser.ConstValue)
		buf.WriteString(MustFormatKeyword(cv.LCurKeyword.Keyword))
		for i := range values {
			buf.WriteString(MustFormatConstValue(values[i]))
		}
		buf.WriteString(MustFormatKeyword(cv.RCurKeyword.Keyword))
	case "pair":
		key := cv.Key.(*parser.ConstValue)
		value := cv.Value.(*parser.ConstValue)

		buf.WriteString(fmt.Sprintf("%s%s %s%s", MustFormatConstValue(key), MustFormatKeyword(cv.ColonKeyword.Keyword), MustFormatConstValue(value), sep))
	case "identifier":
		buf.WriteString(fmt.Sprintf("%s%s", cv.Value.(string), sep))
	case "string":
		fmt.Println("value:", cv.Value)
		val := ""
		if _, ok := cv.Value.(string); ok {
			val = cv.Value.(string)
		} else {
			fmt.Println("type literal")
			val = cv.Value.(*parser.Literal).Value.Text
		}
		buf.WriteString(fmt.Sprintf("%q%s", val, sep))
	case "i64":
		buf.WriteString(fmt.Sprintf("%s%s", cv.ValueInText, sep))
	case "double":
		buf.WriteString(fmt.Sprintf("%s%s", cv.ValueInText, sep))
	}

	return buf.String()
}
