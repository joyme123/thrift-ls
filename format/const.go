package format

import (
	"fmt"

	"github.com/joyme123/thrift-ls/parser"
)

func MustFormatConst(cst *parser.Const) string {
	res := fmt.Sprintf("const %s %s = %s\n",
		MustFormatFieldType(cst.ConstType),
		cst.Name.Name, MustFormatConstValue(cst.Value))

	return res
}
