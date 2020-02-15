package terror

type Terror interface {
	Error() string
	Code() int
}

type errorStringWithCode struct {
	msg  string
	code int
}

func (e errorStringWithCode) Error() string {
	return e.msg
}

func (e errorStringWithCode) Code() int {
	return e.code
}

func New(code int, msg string) Terror {
	return errorStringWithCode{
		msg:  msg,
		code: code,
	}
}
