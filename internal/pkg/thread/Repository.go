package thread

type Repository interface {
	CreateThread(slug string, requestBody *RequestBody) (*Thread, error)
	GetThreadList(slug string, limit int, since string, desc bool) (*[]Thread, error)
}
