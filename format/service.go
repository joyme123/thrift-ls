package format

import (
	"bytes"
	"fmt"

	"github.com/joyme123/thrift-ls/parser"
)

func MustFormatService(svc *parser.Service) string {
	extends := ""
	if svc.Extends != nil {
		extends = fmt.Sprintf("extends %s ", extends)
	}

	buf := bytes.NewBufferString(fmt.Sprintf("service %s %s{\n", svc.Name.Name, extends))

	for _, fn := range svc.Functions {
		buf.WriteString(fmt.Sprintf("  %s\n", MustFormatFunction(fn)))
	}

	buf.WriteString("}\n")
	return buf.String()
}
