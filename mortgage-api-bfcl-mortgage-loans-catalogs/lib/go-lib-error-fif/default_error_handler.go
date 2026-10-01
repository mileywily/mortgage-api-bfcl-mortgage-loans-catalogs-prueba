package error_fif

func DefaultErrorHandler() ErrorHandler {
	return func(err error) (int, interface{}) {
		return CreateInternalServerError(err.Error())
	}
}
