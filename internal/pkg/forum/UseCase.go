package forum

type UseCase interface {
	CreateForum(requestBody *RequestBody) (*Forum, error)
	GetForumDetails(slug string) (*Forum, error)
}
