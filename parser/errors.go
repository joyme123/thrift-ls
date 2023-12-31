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

	InvalidStructError           error = errors.New("expecting a valid struct definition")
	InvalidStructIdentifierError error = errors.New("expecting a valid struct identifier")
	InvalidStructBlockLCURError  error = errors.New("expecting a starting '{' of struct block")
	InvalidStructBlockRCURError  error = errors.New("expecting an ending '}' of struct block")
	InvalidStructFieldError      error = errors.New("expecting a valid struct field")

	InvalidUnionError           error = errors.New("expecting a valid union definition")
	InvalidUnionIdentifierError error = errors.New("expecting a valid union identifier")
	InvalidUnionBlockLCURError  error = errors.New("expecting a starting '{' of union block")
	InvalidUnionBlockRCURError  error = errors.New("expecting a ending '}' of union block")
	InvalidUnionFieldError      error = errors.New("expecting a valid union field")

	InvalidExceptionError           error = errors.New("expecting a valid exception definition")
	InvalidExceptionIdentifierError error = errors.New("expecting a valid exception identifier")
	InvalidExceptionBlockLCURError  error = errors.New("expecting a starting '{' of exception block")
	InvalidExceptionBlockRCURError  error = errors.New("expecting a ending '}' of exception block")
	InvalidExceptionFieldError      error = errors.New("expecting a valid exception field")

	InvalidEnumError                 error = errors.New("expecting a valid enum definition")
	InvalidEnumIdentifierError       error = errors.New("expecting a valid enum identifier")
	InvalidEnumBlockLCURError        error = errors.New("expecting a starting '{' of enum block")
	InvalidEnumBlockRCURError        error = errors.New("expecting a ending '}' of enum block")
	InvalidEnumValueError            error = errors.New("expecting a valid enum field")
	InvalidEnumValueIntConstantError error = errors.New("expecting a valid int contant")

	InvalidTypedefError           error = errors.New("expecting a valid typedef definition")
	InvalidTypedefIdentifierError error = errors.New("expecting a valid typedef identifier")

	InvalidConstError             error = errors.New("expecting a valid const definition")
	InvalidConstConstValueError   error = errors.New("expecting a valid const value")
	InvalidConstMissingValueError error = errors.New("expecting a const value")
	InvalidConstIdentifierError   error = errors.New("expecting a valid const identifier")

	InvalidServiceIdentifierError error = errors.New("expecting a valid service identifier")
	InvalidServiceBlockRCURError  error = errors.New("expecting a ending '}' of service block")
	InvalidServiceFunctionError   error = errors.New("expecting a valid service function")

	InvalidFunctionIdentifierError error = errors.New("expecting a valid function identifier")
	InvalidFunctionArgumentError   error = errors.New("expecting a valid function argument")

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

	InvalidServiceError error = errors.New("expecting a valid service definition")
)
