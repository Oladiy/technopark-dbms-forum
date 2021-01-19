package thread

import (
	"technopark-dbms-forum/internal/pkg/post"
	"technopark-dbms-forum/internal/pkg/vote"
)

type Repository interface {
	CreateThreadPosts(forumSlug string, threadId int, posts *[]post.RequestBody) (*[]post.Post, error)
	GetThreadPosts(threadId int, threadSlug string, limit int, since int, sort string, desc bool) (*[]post.Post, error)
	GetThread(forumSlug string, threadId int, threadSlug string) (*Thread, error)
	UpdateThread(threadId int, threadToUpdate *RequestBody) (*Thread, error)
	ThreadVote(threadId int, threadSlug string, userVote *vote.Vote) (*Thread, error)
}
