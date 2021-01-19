package repository

import (
	"database/sql"
	customErrors "technopark-dbms-forum/internal/pkg/common/custom_errors"
	"technopark-dbms-forum/internal/pkg/post/models"
)

type PostRepository struct {
	connectionDB *sql.DB
}

func NewPostRepository(connectionDB *sql.DB) *PostRepository {
	return &PostRepository {
		connectionDB: connectionDB,
	}
}

func (t *PostRepository) GetPost(id int) (*models.Post, error) {
	querySelect := `SELECT id, parent, author, message, isEdited, forum, thread, created
					FROM Post
					WHERE id = $1`
	p := new(models.Post)

	err := t.connectionDB.QueryRow(querySelect, id).Scan(&p.Id, &p.Parent, &p.Author, &p.Message, &p.IsEdited, &p.Forum,
														 &p.Thread, &p.Created)
	if err != nil {
		return nil, customErrors.PostNotFound
	}

	return p, nil
}

func (t *PostRepository) UpdatePost(id int, message string) (*models.Post, error) {
	queryUpdate := `UPDATE Post 
					SET message = $1, isEdited = true
					WHERE id = $2
					RETURNING id, parent, author, message, isEdited, forum, thread, created`
	p := new(models.Post)
	err := t.connectionDB.QueryRow(queryUpdate, message, id).Scan(&p.Id, &p.Parent, &p.Author, &p.Message, &p.IsEdited,
																  &p.Forum, &p.Thread, &p.Created)
	if err != nil {
		return nil, customErrors.PostNotFound
	}

	return p, nil
}
