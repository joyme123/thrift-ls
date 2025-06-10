package format

import (
	"github.com/joyme123/thrift-ls/parser"
)

const (
	serviceOneLineTpl = `{{.Comments}}{{.Service}} {{.Identifier}}{{.Extends}}{{.ExtendServiceName}} {{.LCUR}}{{.RCUR}}{{.Annotations}}{{.EndLineComments}}`

	serviceMultiLineTpl = `{{.Comments}}{{.Service}} {{.Identifier}}{{.Extends}}{{.ExtendServiceName}} {{.LCUR}}
{{.Functions}}
{{.RCUR}}{{.Annotations}}{{.EndLineComments}}
`
)

type ServiceFormatter struct {
	Comments          string
	Service           string
	Identifier        string
	LCUR              string
	Functions         string
	RCUR              string
	Annotations       string
	EndLineComments   string
	Extends           string
	ExtendServiceName string
}

func MustFormatService(svc *parser.Service) string {
	comments, annos := formatCommentsAndAnnos(svc.Comments, svc.Annotations, "")
	if len(svc.Comments) > 0 && lineDistance(svc.Comments[len(svc.Comments)-1], svc.ServiceKeyword) > 1 {
		comments = comments + "\n"
	}

	f := ServiceFormatter{
		Comments:        comments,
		Service:         MustFormatKeyword(svc.ServiceKeyword.Keyword),
		Identifier:      MustFormatIdentifier(svc.Name, ""),
		LCUR:            MustFormatKeyword(svc.LCurKeyword.Keyword),
		Functions:       MustFormatFunctions(svc.Functions, Indent),
		RCUR:            MustFormatKeyword(svc.RCurKeyword.Keyword),
		Annotations:     annos,
		EndLineComments: MustFormatEndLineComments(svc.EndLineComments, ""),
	}

	if svc.ExtendsKeyword != nil {
		f.Extends = " " + MustFormatKeyword(svc.ExtendsKeyword.Keyword)
	}
	if svc.Extends != nil {
		f.ExtendServiceName = " " + MustFormatIdentifier(svc.Extends, "")
	}

	if len(svc.Functions) > 0 {
		return MustFormat(serviceMultiLineTpl, f)
	}

	return MustFormat(serviceOneLineTpl, f)
}
