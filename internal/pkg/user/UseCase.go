package user

import "technopark-dbms-forum/internal/pkg/user/models"

type UseCase interface {
	CreateUser(nickname string, user *models.RequestBody) (*[]models.User, error)
	GetUserProfile(nickname string) (*models.User, error)
	UpdateUserProfile(nickname string, user *models.RequestBody) (*models.User, error)
}
