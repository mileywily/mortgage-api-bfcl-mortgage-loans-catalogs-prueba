package error_fif

type ErrorHandler func(error) (int, interface{})

type ErrorHandlerMiddleware func(ErrorHandler) ErrorHandler
