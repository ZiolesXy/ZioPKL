package request

type JsonResponse struct {
	Message string `json:"messege"`
	Data any `json:"data"`
}

func NewJsonResponse(message string, data any) JsonResponse {
	return JsonResponse{
		Message: message,
		Data: data,
	}
}