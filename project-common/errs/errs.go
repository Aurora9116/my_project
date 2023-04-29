package errs

type ErrorCode int

type BError struct {
	Code ErrorCode
	Msg  string
}

func (e *BError) Error() string {
	return e.Msg
}

func NewError(code ErrorCode, msg string) *BError {
	return &BError{
		Code: code,
		Msg:  msg,
	}
}
