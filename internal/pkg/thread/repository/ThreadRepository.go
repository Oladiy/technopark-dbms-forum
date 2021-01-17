package repository

import (
	"database/sql"
)

type ThreadRepository struct {
	connectionDB *sql.DB
}

func NewThreadRepository(connectionDB *sql.DB) *ThreadRepository {
	return &ThreadRepository {
		connectionDB: connectionDB,
	}
}
