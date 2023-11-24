package format

import (
	"github.com/joyme123/thrift-ls/parser"
)

const includeTpl = "{{.Comments}}{{.Include}} {{.Path}}{{.EndLineComments}}\n"

type IncludeFormatter struct {
	Comments        string
	Include         string
	Path            string
	EndLineComments string
}

func MustFormatInclude(inc *parser.Include) string {
	comments, _ := formatCommentsAndAnnos(inc.Comments, nil, "")
	if len(inc.Comments) > 0 && lineDistance(inc.Comments[len(inc.Comments)-1], inc.IncludeKeyword) > 1 {
		comments = comments + "\n"
	}

	f := &IncludeFormatter{
		Comments:        comments,
		Include:         MustFormatKeyword(inc.IncludeKeyword.Keyword),
		Path:            MustFormatLiteral(inc.Path, ""),
		EndLineComments: MustFormatComments(inc.EndLineComments, ""),
	}

	return MustFormat(includeTpl, f)
}

func MustFormatCPPInclude(inc *parser.CPPInclude) string {
	comments, _ := formatCommentsAndAnnos(inc.Comments, nil, "")
	if len(inc.Comments) > 0 && lineDistance(inc.Comments[len(inc.Comments)-1], inc.CPPIncludeKeyword) > 1 {
		comments = comments + "\n"
	}

	f := &IncludeFormatter{
		Comments:        comments,
		Include:         MustFormatKeyword(inc.CPPIncludeKeyword.Keyword),
		Path:            MustFormatLiteral(inc.Path, ""),
		EndLineComments: MustFormatComments(inc.EndLineComments, ""),
	}

	return MustFormat(includeTpl, f)
}
