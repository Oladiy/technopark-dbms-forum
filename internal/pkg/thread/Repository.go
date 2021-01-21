package thread

import (
	"technopark-dbms-forum/internal/pkg/post/models"
	models2 "technopark-dbms-forum/internal/pkg/thread/models"
	models3 "technopark-dbms-forum/internal/pkg/vote/models"
)

type Repository interface {
	CreateThreadPosts(forumSlug string, threadId int, posts *[]models.RequestBody) (*[]models.Post, error)
	GetThreadPosts(threadId int, limit int, since int, sort string, desc bool) (*[]models.Post, error)
	GetThread(forumSlug string, threadId int, threadSlug string) (*models2.Thread, error)
	UpdateThread(threadId int, threadToUpdate *models2.RequestBody) (*models2.Thread, error)
	ThreadVote(threadId int, userVote *models3.Vote) (*models2.Thread, error)
}
