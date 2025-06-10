package format

import (
	"github.com/joyme123/thrift-ls/parser"
)

const (
	structOneLineTpl = `{{.Comments}}{{.Struct}} {{.Identifier}} {{.LCUR}}{{.RCUR}}{{.Annotations}}{{.EndLineComments}}`

	structMultiLineTpl = `{{.Comments}}{{.Struct}} {{.Identifier}} {{.LCUR}}
{{.Fields}}{{.RCUR}}{{.Annotations}}{{.EndLineComments}}
`
)

type StructFormatter struct {
	Comments        string
	Struct          string
	Identifier      string
	LCUR            string
	Fields          string
	RCUR            string
	Annotations     string
	EndLineComments string
}

func MustFormatStruct(st *parser.Struct) string {
	comments, annos := formatCommentsAndAnnos(st.Comments, st.Annotations, "")

	if len(st.Comments) > 0 && lineDistance(st.Comments[len(st.Comments)-1], st.StructKeyword) > 1 {
		comments = comments + "\n"
	}

	f := StructFormatter{
		Comments:        comments,
		Struct:          MustFormatKeyword(st.StructKeyword.Keyword),
		Identifier:      MustFormatIdentifier(st.Identifier, ""),
		LCUR:            MustFormatKeyword(st.LCurKeyword.Keyword),
		Fields:          MustFormatFields(st.Fields, Indent),
		RCUR:            MustFormatKeyword(st.RCurKeyword.Keyword),
		Annotations:     annos,
		EndLineComments: MustFormatEndLineComments(st.EndLineComments, ""),
	}

	if len(st.Fields) > 0 {
		return MustFormat(structMultiLineTpl, f)
	}

	return MustFormat(structOneLineTpl, f)
}
