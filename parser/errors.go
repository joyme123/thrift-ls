package parser

import (
	"errors"
)

// ErrorLister is the public interface to access the inner errors
// included in a errList
type ErrorLister interface {
	Errors() []error
}

func (e errList) Errors() []error {
	return e
}

// ParserError is the public interface to errors of type parserError
type ParserError interface {
	Error() string
	InnerError() error
	Pos() (int, int, int)
	Expected() []string
}

func (p *parserError) InnerError() error {
	return p.Inner
}

func (p *parserError) Pos() (line, col, offset int) {
	return p.pos.line, p.pos.col, p.pos.offset
}

func (p *parserError) Expected() []string {
	return p.expected
}

var (
	RequiredError          error = errors.New("expecting 'required' or 'optional'")
	InvalidFieldTypeError  error = errors.New("expecting a valid field type")
	InvalidFieldIndexError error = errors.New("expecting a valid int16 field index")

	InvalidStructIdentifierError error = errors.New("expecting a valid struct identifier")
	InvalidStructBlockLCURError  error = errors.New("expecting a starting '{' of struct block")
	InvalidStructBlockRCURError  error = errors.New("expecting an ending '}' of struct block")
	InvalidStructFieldError      error = errors.New("expecting a valid struct field")

	InvalidUnionIdentifierError error = errors.New("expecting a valid union identifier")
	InvalidUnionBlockLCURError  error = errors.New("expecting a starting '{' of union block")
	InvalidUnionBlockRCURError  error = errors.New("expecting a ending '}' of union block")
	InvalidUnionFieldError      error = errors.New("expecting a valid union field")

	InvalidIdentifierError error = errors.New("expecting a valid identifier")

	InvalidLiteral1MissingRightError error = errors.New("expecting a right \" ")
	InvalidLiteral1Error             error = errors.New("expecting a valid literal")
	InvalidLiteral2MissingRightError error = errors.New("expecting a right ' ")
	InvalidLiteral2Error             error = errors.New("expecting a valid literal")

	InvalidHeaderError error = errors.New("expecting a valid header")

	InvalidIncludeError    error = errors.New("expecting a valid include header")
	InvalidCppIncludeError error = errors.New("expecting a valid cpp include header")
	InvalidNamespaceError  error = errors.New("expecting a valid namespace header")

	InvalidDefinitionError error = errors.New("expecting a valid definition")

	InvalidConstError error = errors.New("expecting a valid const definition")

	InvalidTypedefError error = errors.New("expecting a valid typedef definition")

	InvalidEnumError error = errors.New("expecting a valid enum definition")

	InvalidServiceError error = errors.New("expecting a valid service definition")

	InvalidStructError error = errors.New("expecting a valid struct definition")

	InvalidUnionError error = errors.New("expecting a valid union definition")

	InvalidExceptionError error = errors.New("expecting a valid exception definition")
)
