package parser

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ParseRecursively(t *testing.T) {
	baseFile := `
struct User {
	1: required string Name,
	2: optional i32 Age,
}
struct GetUserRequest {
	1: required string Name,
}
exception Error {
    1: string Code,
    2: string Message,
    3: map<string, string> Data,
}
`
	serviceFile := `
include "base.thrift"

service APIService {
	base.User GetUser(1: base.GetUserRequest req) throws (1: base.Error err),
}
`

	parser := &PEGParser{}
	parseResult := parser.ParseRecursively("service.thrift", []byte(serviceFile), 10, func(include string) (filename string, content []byte, err error) {
		if include == "base.thrift" {
			return "base.thrift", []byte(baseFile), nil
		}
		return "", nil, errors.New("file not found")
	})

	assert.Len(t, parseResult, 2)
	for _, res := range parseResult {
		if res.Doc.Filename == "service.thrift" {
			assert.Len(t, res.Errors, 0)
			assert.Len(t, res.Doc.Services, 1)
			assert.Len(t, res.Doc.Services[0].Functions, 1)
		} else if res.Doc.Filename == "base.thrift" {
			assert.Len(t, res.Errors, 0)
			assert.Len(t, res.Doc.Structs, 2)
			assert.Len(t, res.Doc.Exceptions, 1)
		}
	}
}
