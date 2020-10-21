package validator

type basicError struct {
	message string
}

func (e basicError) Error() string {
	return e.message
}

// WrongTypeError means message format is wrong
type WrongTypeError struct {
	basicError
}

func newWrongTypeError(msg string) WrongTypeError {
	return WrongTypeError{
		basicError: basicError{
			message: msg,
		},
	}
}

// NotExistError means parameter do not exist
type NotExistError struct {
	basicError
}

func newNotExistError(msg string) NotExistError {
	return NotExistError{
		basicError: basicError{
			message: msg,
		},
	}
}
