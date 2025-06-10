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
	comments, annos := formatCommentsAndAnnos(td.Comments, td.Annotations, "")

	if len(td.Comments) > 0 && lineDistance(td.Comments[len(td.Comments)-1], td.TypedefKeyword) > 1 {
		comments = comments + "\n"
	}

	f := &TypedefFormatter{
		Comments:        comments,
		Typedef:         MustFormatKeyword(td.TypedefKeyword.Keyword),
		Type:            MustFormatFieldType(td.T),
		Name:            MustFormatIdentifier(td.Alias, ""),
		Annotations:     annos,
		EndLineComments: MustFormatEndLineComments(td.EndLineComments, ""),
	}

	return MustFormat(typedefOneLineTpl, f)
}
