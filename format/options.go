package format

import (
	"flag"
	"strconv"
	"strings"
)

var Indent = "    "
var Align = "field"

type Options struct {
	// Do not print reformatted sources to standard output.
	// If a file's formatting is different from thriftls's, overwrite it
	// with thrfitls's version.
	Write bool `yaml:"rewrite"`

	// Indent to use. Support: nspace(s), ntab(s). example: 4spaces, 1tab, tab
	// if indent format is invalid or not specified, default is 4spaces
	Indent string `yaml:"indent"`

	// Do not print reformatted sources to standard output.
	// If a file's formatting is different than gofmt's, print diffs
	// to standard output.
	Diff bool `yaml:"diff"`

	// Align enables align option for struct/enum/exception/union fields
	// Options: "field", "assign", "disable"
	// Default is "field" if not set
	Align string `yaml:"alignByAssign"`
}

func (o *Options) SetFlags() {
	flag.BoolVar(&o.Write, "w", false, "Do not print reformatted sources to standard output. If a file's formatting is different from thriftls's, overwrite it with thrfitls's version.")
	flag.BoolVar(&o.Diff, "d", false, "Do not print reformatted sources to standard output. If a file's formatting is different than gofmt's, print diffs to standard output.")
	flag.StringVar(&o.Indent, "indent", "4spaces", "Indent to use. Support: num*space, num*tab. example: 4spaces, 1tab, tab")
	flag.StringVar(&o.Align, "align", "field", `Align enables align option for struct/enum/exception/union fields, Options: "field", "assign", "disable", Default is "field" if not set.`)
}

func (o *Options) InitDefault() {
	Indent = o.GetIndent()

	if o.Align == "" || (o.Align != AlignTypeField && o.Align != AlignTypeAssign && o.Align != AlignTypeDisable) {
		o.Align = "field"
	}

	Align = o.Align
}

func (o *Options) GetIndent() string {
	if o.Indent == "" {
		o.Indent = "4spaces"
	}

	indent := o.Indent
	suffixes := []string{"spaces", "space", "tabs", "tab"}
	for _, suffix := range suffixes {
		if strings.HasSuffix(indent, suffix) {
			char := ""
			if strings.HasPrefix(suffix, "tab") {
				char = "	"
			} else {
				char = " "
			}
			num := 1
			numStr := strings.TrimSuffix(indent, suffix)
			if len(numStr) == 0 {
				num = 1
			} else {
				num, _ = strconv.Atoi(numStr)
				if num == 0 {
					num = 4
					char = " "
				}
			}

			return strings.Repeat(char, num)
		}
	}

	return "    "

}
