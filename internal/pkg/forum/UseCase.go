package forum

import (
	"technopark-dbms-forum/internal/pkg/thread"
)

type UseCase interface {
	CreateForum(requestBody *RequestBody) (*Forum, error)
	GetForumDetails(slug string) (*Forum, error)
	GetForumThreadList(slug string, limit int, since string, desc bool) (*[]thread.Thread, error)
}
