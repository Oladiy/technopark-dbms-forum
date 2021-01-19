package post

import "technopark-dbms-forum/internal/pkg/post/models"

type Repository interface {
	GetPost(id int) (*models.Post, error)
	UpdatePost(id int, message string) (*models.Post, error)
}
