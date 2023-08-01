package format

import (
	"github.com/joyme123/thrift-ls/parser"
)

const (
	exceptionOneLineTpl = `{{.Comments}}{{.Exception}} {{.Identifier}} {{.LCUR}}{{.RCUR}}{{.Annotations}}{{.EndLineComments}}`

	exceptionMultiLineTpl = `{{.Comments}}{{.Exception}} {{.Identifier}} {{.LCUR}}
{{.Fields}}
{{.RCUR}}{{.Annotations}}{{.EndLineComments}}
`
)

type ExceptionFormatter struct {
	Comments        string
	Exception       string
	Identifier      string
	LCUR            string
	Fields          string
	RCUR            string
	Annotations     string
	EndLineComments string
}

func MustFormatException(excep *parser.Exception) string {
	comments, annos := formatCommentsAndAnnos(excep.Comments, excep.Annotations, "")
	if len(excep.Comments) > 0 && lineDistance(excep.Comments[len(excep.Comments)-1], excep.ExceptionKeyword) > 1 {
		comments = comments + "\n"
	}
	f := ExceptionFormatter{
		Comments:        comments,
		Exception:       MustFormatKeyword(excep.ExceptionKeyword.Keyword),
		Identifier:      MustFormatIdentifier(excep.Name),
		LCUR:            MustFormatKeyword(excep.LCurKeyword.Keyword),
		Fields:          MustFormatFields(excep.Fields, Indent),
		RCUR:            MustFormatKeyword(excep.RCurKeyword.Keyword),
		Annotations:     annos,
		EndLineComments: MustFormatComments(excep.EndLineComments, ""),
	}

	if len(excep.Fields) > 0 {
		return MustFormat(exceptionMultiLineTpl, f)
	}

	return MustFormat(exceptionOneLineTpl, f)
}
