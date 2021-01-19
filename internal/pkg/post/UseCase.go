package post

import "technopark-dbms-forum/internal/pkg/post/models"

type UseCase interface {
	GetPost(id int, isUserRelated bool, isThreadRelated bool, isForumRelated bool) (*models.Details, error)
	UpdatePost(id int, message string) (*models.Post, error)
}
