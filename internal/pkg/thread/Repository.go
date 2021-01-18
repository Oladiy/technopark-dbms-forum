package thread

import (
	"technopark-dbms-forum/internal/pkg/post"
)

type Repository interface {
	CreateThreadPosts(forumSlug string, threadId int, posts *[]post.RequestBody) (*[]post.Post, error)
	GetThread(forumSlug string, threadId int, threadSlug string) (*Thread, error)
}
