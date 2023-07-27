package format

import (
	"github.com/joyme123/thrift-ls/parser"
)

const typedefOneLineTpl = `{{.Comments}}{{.Typedef}} {{.Type}} {{.Name}}{{.Annotations}}{{.EndLineComments}}
`

type TypedefFormatter struct {
	Comments        string
	Typedef         string
	Type            string
	Name            string
	Annotations     string
	EndLineComments string
}

func MustFormatTypedef(td *parser.Typedef) string {
	comments, annos := formatCommentsAndAnnos(td.Comments, td.Annotations)

	f := &TypedefFormatter{
		Comments:        comments,
		Typedef:         MustFormatKeyword(td.TypedefKeyword.Keyword),
		Type:            MustFormatFieldType(td.T),
		Name:            MustFormatIdentifier(td.Alias),
		Annotations:     annos,
		EndLineComments: MustFormatComments(td.EndLineComments),
	}

	return MustFormat(typedefOneLineTpl, f)
}
