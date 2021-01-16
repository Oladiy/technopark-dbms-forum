package thread

type Repository interface {
	CreateThread(slug string, requestBody *RequestBody) (*Thread, error)
}
