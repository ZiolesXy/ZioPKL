package respons

type JsonResponse struct {
	Message string `json:"messege"`
	Data any `json:"data"`
}

type Entries[T any] struct {
	Entries []T `json:"entries"`
}

func NewJsonResponse(message string, data any) JsonResponse {
	return JsonResponse{
		Message: message,
		Data: data,
	}
}

func NewEntries[T any](data []T) Entries[T] {
	return Entries[T]{
		Entries: data,
	}
}