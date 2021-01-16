package repository

import (
	"database/sql"
	customErrors "technopark-dbms-forum/internal/pkg/common/custom_errors"
	"technopark-dbms-forum/internal/pkg/forum"
)

type ForumRepository struct {
	connectionDB *sql.DB
}

func NewForumRepository(connectionDB *sql.DB) *ForumRepository {
	return &ForumRepository {
		connectionDB: connectionDB,
	}
}

func (t *ForumRepository) CreateForum(requestBody *forum.RequestBody) (*forum.Forum, error) {
	queryInsert := `INSERT INTO Forum (title, author, slug) 
					VALUES($1, $2, $3) 
					returning title, author, slug;`
	querySelect := `SELECT title, author, slug, posts, threads 
					FROM Forum 
					WHERE slug = $1;`
	f := new(forum.Forum)

	row := t.connectionDB.QueryRow(queryInsert, requestBody.Title, requestBody.User, requestBody.Slug)
	err := row.Scan(&f.Title, &f.User, &f.Slug)
	if err != nil {
		selection := t.connectionDB.QueryRow(querySelect, requestBody.Slug)
		err := selection.Scan(&f.Title, &f.User, &f.Slug, &f.Posts, &f.Threads)
		if selection.Err() != nil || err != nil {
			return nil, customErrors.IncorrectInputData
		}

		return f, customErrors.ForumAlreadyExists
	}

	return f, nil
}

func (t *ForumRepository) GetForumDetails(slug string) (*forum.Forum, error) {
	querySelect := `SELECT title, author, slug, posts, threads 
					FROM Forum 
					WHERE slug = $1;`
	selection := t.connectionDB.QueryRow(querySelect, slug)
	if selection == nil {
		return nil, customErrors.ForumSlugNotFound
	}

	f := new(forum.Forum)
	if err := selection.Scan(&f.Title, &f.User, &f.Slug, &f.Posts, &f.Threads); err != nil {
		return nil, customErrors.ForumSlugNotFound
	}

	return f, nil
}
