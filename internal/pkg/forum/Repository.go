package forum

type Repository interface {
	CreateForum(requestBody *RequestBody) (*Forum, error)
	GetForumDetails(slug string) (*Forum, error)
}
