package validator

const (
	errWrongType = iota
	errNotExist
)

// Error is error type of router
type Error struct {
	message string
	errno   int
}

func (e Error) Error() string {
	return e.message
}

// IsWrongType means type of val is wrong
func (e Error) IsWrongType() bool {
	return e.errno == errWrongType
}

// IsNotExist means val doesn't exist
func (e Error) IsNotExist() bool {
	return e.errno == errNotExist
}
