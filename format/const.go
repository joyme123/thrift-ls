package format

import (
	"github.com/joyme123/thrift-ls/parser"
)

const constOneLineTpl = `{{.Comments}}{{.Const}} {{.Type}} {{.Name}} {{.Equal}} {{.Value}}{{.Annotations}}{{.EndLineComments}}
`

type ConstFormatter struct {
	Comments        string
	Const           string
	Type            string
	Name            string
	Annotations     string
	Equal           string
	Value           string
	EndLineComments string
}

func MustFormatConst(cst *parser.Const) string {
	comments, annos := formatCommentsAndAnnos(cst.Comments, cst.Annotations)

	f := &ConstFormatter{
		Comments:        comments,
		Const:           MustFormatKeyword(cst.ConstKeyword.Keyword),
		Type:            MustFormatFieldType(cst.ConstType),
		Name:            MustFormatIdentifier(cst.Name),
		Annotations:     annos,
		Equal:           MustFormatKeyword(cst.EqualKeyword.Keyword),
		Value:           MustFormatConstValue(cst.Value),
		EndLineComments: MustFormatComments(cst.EndLineComments),
	}

	return MustFormat(constOneLineTpl, f)
}
