package format

import (
	"bytes"

	"github.com/joyme123/thrift-ls/parser"
)

func MustFormatFunctions(fns []*parser.Function, indent string) string {
	buf := bytes.NewBuffer(nil)
	for i := range fns {
		buf.WriteString(indent + MustFormatFunction(fns[i]))
		if i < len(fns)-1 {
			buf.WriteString("\n")
		}
	}

	return buf.String()
}

const functionTpl = "{{.Comments}} {{.Oneway}}{{.FunctionType}} {{.Identifier}} {{.LPAR}}{{.Args}}{{.RPAR}} {{.Throws}}{{.Annotations}}{{.EndLineComments}}"

type FunctionFormatter struct {
	Comments        string
	Oneway          string
	FunctionType    string
	Identifier      string
	LPAR            string
	Args            string
	RPAR            string
	Throws          string
	Annotations     string
	EndLineComments string
}

func MustFormatFunction(fn *parser.Function) string {
	comments, annos := formatCommentsAndAnnos(fn.Comments, fn.Annotations)
	oneway := ""
	if fn.Oneway != nil {
		oneway = "oneway "
	}
	args := ""
	if len(fn.Arguments) > 0 {
		args = " " + MustFormatOneLineFields(fn.Arguments) + " "
	}

	ft := ""
	if fn.Void != nil {
		ft = MustFormatKeyword(fn.Void.Keyword)
	} else {
		ft = MustFormatFieldType(fn.FunctionType)
	}

	f := &FunctionFormatter{
		Comments:        comments,
		Oneway:          oneway,
		FunctionType:    ft,
		Identifier:      MustFormatIdentifier(fn.Name),
		LPAR:            MustFormatKeyword(fn.LParKeyword.Keyword),
		Args:            args,
		RPAR:            MustFormatKeyword(fn.RParKeyword.Keyword),
		Throws:          MustFormatThrows(fn.Throws),
		Annotations:     annos,
		EndLineComments: MustFormatComments(fn.EndLineComments),
	}

	return MustFormat(functionTpl, f)

}

const throwTpl = "{{.Throw}} {{.LPAR}}{{.Fields}}{{.RPAR}} "

type ThrowFormatter struct {
	Throw  string
	LPAR   string
	Fields string
	RPAR   string
}

func MustFormatThrows(throws *parser.Throws) string {
	if throws == nil {
		return ""
	}

	args := ""
	if len(throws.Fields) > 0 {
		args = " " + MustFormatOneLineFields(throws.Fields) + " "
	}

	f := &ThrowFormatter{
		Throw:  MustFormatKeyword(throws.ThrowsKeyword.Keyword),
		LPAR:   MustFormatKeyword(throws.LParKeyword.Keyword),
		Fields: args,
		RPAR:   MustFormatKeyword(throws.RParKeyword.Keyword),
	}

	return MustFormat(throwTpl, f)

}
