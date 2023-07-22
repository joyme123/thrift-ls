package errors

import "fmt"

type AggregateError struct {
	errs []error
}

func NewAggregate(errs []error) *AggregateError {
	return &AggregateError{
		errs: errs,
	}
}

func (e *AggregateError) Error() string {
	if len(e.errs) == 0 {
		return ""
	}
	return fmt.Sprintf("%v", e.errs)
}
