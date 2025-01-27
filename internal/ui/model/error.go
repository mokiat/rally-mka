package model

func NewErrorModel() *ErrorModel {
	return &ErrorModel{}
}

type ErrorModel struct {
	err error
}

func (e *ErrorModel) Error() error {
	return e.err
}

func (e *ErrorModel) SetError(err error) {
	e.err = err
}
