package forum

import (
	"technopark-dbms-forum/internal/pkg/thread"
	"technopark-dbms-forum/internal/pkg/user"
)

type UseCase interface {
	CreateForum(requestBody *RequestBody) (*Forum, error)
	CreateForumThread(slug string, requestBody *thread.RequestBody) (*thread.Thread, error)
	GetForumDetails(slug string) (*Forum, error)
	GetForumThreads(slug string, limit int, since string, desc bool) (*[]thread.Thread, error)
	GetForumUsers(slug string, limit int, since string, desc bool) (*[]user.User, error)
}
