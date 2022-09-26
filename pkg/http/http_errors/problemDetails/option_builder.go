package problemDetails

import "reflect"

type OptionBuilder struct {
	internalErrors map[reflect.Type]func(err error) ProblemDetailErr
}

func NewOptionBuilder() *OptionBuilder {
	return &OptionBuilder{}
}

func (p *OptionBuilder) Map(srcErrorType reflect.Type, problem ProblemDetailFunc[error]) *OptionBuilder {
	internalErrorMaps[srcErrorType] = problem

	return p
}

func (p *OptionBuilder) Build() map[reflect.Type]func(err error) ProblemDetailErr {
	return p.internalErrors
}
