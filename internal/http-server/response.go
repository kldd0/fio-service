package http_server

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

func OK() Response {
	return Response{
		Status: StatusOK,
	}
}

func Error(message string) Response {
	return Response{
		Status: StatusError,
		Error:  message,
	}
}
