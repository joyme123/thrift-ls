package format

import (
	"github.com/joyme123/thrift-ls/parser"
)

const (
	enumOneLineTpl = `{{.Comments}}{{.Enum}} {{.Identifier}} {{.LCUR}}{{.RCUR}}{{.Annotations}}{{.EndLineComments}}`

	enumMultiLineTpl = `{{.Comments}}{{.Enum}} {{.Identifier}} {{.LCUR}}
{{.EnumValues}}
{{.RCUR}}{{.Annotations}}{{.EndLineComments}}
`
)

type EnumFormatter struct {
	Comments        string
	Enum            string
	Identifier      string
	LCUR            string
	EnumValues      string
	RCUR            string
	Annotations     string
	EndLineComments string
}

func MustFormatEnum(enum *parser.Enum) string {
	comments, annos := formatCommentsAndAnnos(enum.Comments, enum.Annotations, "")
	if len(enum.Comments) > 0 && lineDistance(enum.Comments[len(enum.Comments)-1], enum.EnumKeyword) > 1 {
		comments = comments + "\n"
	}

	f := EnumFormatter{
		Comments:        comments,
		Enum:            MustFormatKeyword(enum.EnumKeyword.Keyword),
		Identifier:      MustFormatIdentifier(enum.Name),
		LCUR:            MustFormatKeyword(enum.LCurKeyword.Keyword),
		EnumValues:      MustFormatEnumValues(enum.Values, Indent),
		RCUR:            MustFormatKeyword(enum.RCurKeyword.Keyword),
		Annotations:     annos,
		EndLineComments: MustFormatEndLineComments(enum.EndLineComments, ""),
	}

	if len(enum.Values) > 0 {
		return MustFormat(enumMultiLineTpl, f)
	}

	return MustFormat(enumOneLineTpl, f)
}
