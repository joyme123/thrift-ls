package format

import (
	"github.com/joyme123/thrift-ls/parser"
)

const (
	serviceOneLineTpl = `{{.Comments}}{{.Service}} {{.Identifier}} {{.LCUR}}{{.RCUR}}{{.Annotations}}{{.EndLineComments}}`

	serviceMultiLineTpl = `{{.Comments}}{{.Service}} {{.Identifier}} {{.LCUR}}
{{.Functions}}
{{.RCUR}}{{.Annotations}}{{.EndLineComments}}
`
)

type ServiceFormatter struct {
	Comments        string
	Service         string
	Identifier      string
	LCUR            string
	Functions       string
	RCUR            string
	Annotations     string
	EndLineComments string
}

func MustFormatService(svc *parser.Service) string {
	comments, annos := formatCommentsAndAnnos(svc.Comments, svc.Annotations)

	f := ServiceFormatter{
		Comments:        comments,
		Service:         MustFormatKeyword(svc.ServiceKeyword.Keyword),
		Identifier:      MustFormatIdentifier(svc.Name),
		LCUR:            MustFormatKeyword(svc.LCurKeyword.Keyword),
		Functions:       MustFormatFunctions(svc.Functions, Indent),
		RCUR:            MustFormatKeyword(svc.RCurKeyword.Keyword),
		Annotations:     annos,
		EndLineComments: MustFormatComments(svc.EndLineComments),
	}

	if len(svc.Functions) > 0 {
		return MustFormat(serviceMultiLineTpl, f)
	}

	return MustFormat(serviceOneLineTpl, f)
}
