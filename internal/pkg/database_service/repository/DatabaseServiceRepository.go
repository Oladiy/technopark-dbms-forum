package repository

import (
	"database/sql"
	"technopark-dbms-forum/internal/pkg/database_service/models"
)

type DatabaseServiceRepository struct {
	connectionDB *sql.DB
}

func NewDatabaseServiceRepository(connectionDB *sql.DB) *DatabaseServiceRepository {
	return &DatabaseServiceRepository {
		connectionDB: connectionDB,
	}
}

func (t *DatabaseServiceRepository) GetStatus() (*models.Status, error) {
	st := new(models.Status)

	querySelect := `SELECT
					(SELECT COUNT(id) FROM Forum),
					(SELECT COUNT(id) FROM Post),
					(SELECT COUNT(id) FROM Thread),
					(SELECT COUNT(id) FROM Users)`

	err := t.connectionDB.QueryRow(querySelect).Scan(&st.Forum, &st.Post, &st.Thread, &st.User)
	if err != nil {
		return nil, err
	}

	return st, nil
}

func (t *DatabaseServiceRepository) ClearDatabase() error {
	queryClear := `TRUNCATE Users, Forum, Thread, Post, Vote CASCADE;`

	_, err := t.connectionDB.Exec(queryClear)
	if err != nil {
		return err
	}

	return nil
}
