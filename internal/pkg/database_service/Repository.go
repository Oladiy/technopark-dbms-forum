package database_service

import "technopark-dbms-forum/internal/pkg/database_service/models"

type Repository interface {
	GetStatus() (*models.Status, error)
	ClearDatabase() error
}
