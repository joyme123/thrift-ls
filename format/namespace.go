package format

import (
	"bytes"
	"fmt"

	"github.com/joyme123/thrift-ls/parser"
)

func MustFormatNamespace(ns *parser.Namespace) string {
	buf := bytes.NewBuffer([]byte(fmt.Sprintf("namespace %s %s", ns.Language.Name, ns.Name.Name)))
	if ns.Annotations != nil {
		buf.WriteString(fmt.Sprintf(" %s\n", MustFormatAnnotations(ns.Annotations)))
	}

	return buf.String()
}
