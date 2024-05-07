package domain

type SuccessData struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func NewMessageSuccess(msg string, data any) SuccessData {
	return SuccessData{
		Message: msg,
		Data:    data,
	}
}
