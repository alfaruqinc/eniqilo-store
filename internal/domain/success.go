package domain

type SuccessData struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func NewMessageSuccess(data any) SuccessData {
	return SuccessData{
		Message: "success",
		Data:    data,
	}
}
