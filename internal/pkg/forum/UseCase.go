package forum

import (
	models3 "technopark-dbms-forum/internal/pkg/forum/models"
	"technopark-dbms-forum/internal/pkg/thread/models"
	models2 "technopark-dbms-forum/internal/pkg/user/models"
)

type UseCase interface {
	CreateForum(requestBody *models3.RequestBody) (*models3.Forum, error)
	CreateForumThread(slug string, requestBody *models.RequestBody) (*models.Thread, error)
	GetForumDetails(slug string) (*models3.Forum, error)
	GetForumThreads(slug string, limit int, since string, desc bool) (*[]models.Thread, error)
	GetForumUsers(slug string, limit int, since string, desc bool) (*[]models2.User, error)
}
