package error

type ErrorWrapper struct {
	err error
}

func NewErrorWrapper(err error) ErrorWrapper {
	return ErrorWrapper{err: err}
}

func (ew ErrorWrapper) GetError() error {
	return ew.err
}

func (ew ErrorWrapper) Error() string {
	return ew.GetError().Error()
}
