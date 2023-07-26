package format

import (
	"bytes"
	"fmt"

	"github.com/joyme123/thrift-ls/parser"
)

func MustFormatNamespace(ns *parser.Namespace) string {
	buf := bytes.NewBuffer([]byte(fmt.Sprintf("namespace %s %s", ns.Language, ns.Name)))
	if len(ns.Annotations) > 0 {
		buf.WriteString(fmt.Sprintf(" %s\n", MustFormatAnnotations(ns.Annotations)))
	}

	return buf.String()
}
