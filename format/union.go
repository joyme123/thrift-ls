package format

import (
	"github.com/joyme123/thrift-ls/parser"
)

const (
	unionOneLineTpl = `{{.Comments}}{{.Union}} {{.Identifier}} {{.LCUR}}{{.RCUR}}{{.Annotations}}{{.EndLineComments}}`

	unionMultiLineTpl = `{{.Comments}}{{.Union}} {{.Identifier}} {{.LCUR}}
{{.Fields}}
{{.RCUR}}{{.Annotations}}{{.EndLineComments}}
`
)

type UnionFormatter struct {
	Comments        string
	Union           string
	Identifier      string
	LCUR            string
	Fields          string
	RCUR            string
	Annotations     string
	EndLineComments string
}

func MustFormatUnion(union *parser.Union) string {
	comments, annos := formatCommentsAndAnnos(union.Comments, union.Annotations)

	f := UnionFormatter{
		Comments:        comments,
		Union:           MustFormatKeyword(union.UnionKeyword.Keyword),
		Identifier:      MustFormatIdentifier(union.Name),
		LCUR:            MustFormatKeyword(union.LCurKeyword.Keyword),
		Fields:          MustFormatFields(union.Fields, Indent),
		RCUR:            MustFormatKeyword(union.RCurKeyword.Keyword),
		Annotations:     annos,
		EndLineComments: MustFormatComments(union.EndLineComments),
	}

	if len(union.Fields) > 0 {
		return MustFormat(unionMultiLineTpl, f)
	}

	return MustFormat(unionOneLineTpl, f)
}
