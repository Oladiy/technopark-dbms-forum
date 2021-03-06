package thread

import (
	"technopark-dbms-forum/internal/pkg/post/models"
	models2 "technopark-dbms-forum/internal/pkg/thread/models"
	models3 "technopark-dbms-forum/internal/pkg/vote/models"
)

type UseCase interface {
	CreateThreadPosts(forumSlug string, threadId int, posts *[]models.RequestBody) (*[]models.Post, error)
	GetThreadPosts(threadId int, threadSlug string, limit int, since int, sort string, desc bool) (*[]models.Post, error)
	GetThread(forumSlug string, threadId int, threadSlug string) (*models2.Thread, error)
	UpdateThread(threadId int, threadSlug string, threadToUpdate *models2.RequestBody) (*models2.Thread, error)
	ThreadVote(threadId int, threadSlug string, userVote *models3.Vote) (*models2.Thread, error)
}
