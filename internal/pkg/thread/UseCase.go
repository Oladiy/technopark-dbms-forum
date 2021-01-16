package thread

type UseCase interface {
	CreateThread(slug string, requestBody *RequestBody) (*Thread, error)
}
