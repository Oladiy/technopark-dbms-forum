package forum

import (
	"technopark-dbms-forum/internal/pkg/thread"
)

type Repository interface {
	CreateForum(requestBody *RequestBody) (*Forum, error)
	CreateForumThread(slug string, requestBody *thread.RequestBody) (*thread.Thread, error)
	GetForumDetails(slug string) (*Forum, error)
	GetForumThreads(slug string, limit int, since string, desc bool) (*[]thread.Thread, error)
}
