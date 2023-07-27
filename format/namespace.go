package format

import (
	"github.com/joyme123/thrift-ls/parser"
)

const namespaceOneLineTpl = `{{.Comments}}{{.Namespace}} {{.Language}} {{.Name}}{{.Annotations}}{{.EndLineComments}}`

type NamespaceFormatter struct {
	Comments        string
	Namespace       string
	Language        string
	Name            string
	Annotations     string
	EndLineComments string
}

func MustFormatNamespace(ns *parser.Namespace) string {
	comments, annos := formatCommentsAndAnnos(ns.Comments, ns.Annotations)

	f := &NamespaceFormatter{
		Comments:        comments,
		Namespace:       MustFormatKeyword(ns.NamespaceKeyword.Keyword),
		Language:        MustFormatIdentifier(&ns.Language.Identifier),
		Name:            MustFormatIdentifier(ns.Name),
		Annotations:     annos,
		EndLineComments: MustFormatComments(ns.EndLineComments),
	}

	return MustFormat(namespaceOneLineTpl, f)
}
