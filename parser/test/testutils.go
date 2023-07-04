package test

import "github.com/joyme123/thrift-ls/parser"

func containsError(errs []error, target error) bool {
	for _, err := range errs {
		err = err.(parser.ParserError).InnerError()
		if err == target {
			return true
		}
	}
	return false
}
