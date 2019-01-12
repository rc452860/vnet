package api

type Response struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

const (
	SUCCESS = "success"
	FAILURE = "failure"
	ERROR   = "error"
)

func Success(data interface{}, param ...string) *Response {
	r := &Response{
		Code: SUCCESS,
		Data: data,
	}
	if len(param) > 0 {
		r.Message = param[0]
	}
	return r
}

func Failure(message string) *Response {
	r := &Response{
		Code:    FAILURE,
		Message: message,
	}
	return r
}

func Error(err error) *Response {
	r := &Response{
		Code:    ERROR,
		Message: err.Error(),
	}
	return r
}
