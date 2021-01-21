package repository

import (
	"github.com/go-openapi/strfmt"
	"github.com/jackc/pgx"
	customErrors "technopark-dbms-forum/internal/pkg/common/custom_errors"
	"technopark-dbms-forum/internal/pkg/post/models"
)

type PostRepository struct {
	connectionDB *pgx.ConnPool
}

func NewPostRepository(connectionDB *pgx.ConnPool) *PostRepository {
	return &PostRepository {
		connectionDB: connectionDB,
	}
}

func (t *PostRepository) GetPost(id int) (*models.Post, error) {
	querySelect := `SELECT id, parent, author, message, isEdited, forum, thread, created
					FROM Post
					WHERE id = $1`
	p := new(models.Post)
	var created strfmt.DateTime

	err := t.connectionDB.QueryRow(querySelect, id).Scan(&p.Id, &p.Parent, &p.Author, &p.Message, &p.IsEdited, &p.Forum,
														 &p.Thread, &created)
	if err != nil {
		return nil, customErrors.PostNotFound
	}
	p.Created = created.String()

	return p, nil
}

func (t *PostRepository) UpdatePost(id int, message string) (*models.Post, error) {
	queryUpdate := `UPDATE Post 
					SET message = $1, isEdited = true
					WHERE id = $2
					RETURNING id, parent, author, message, isEdited, forum, thread, created`
	p := new(models.Post)
	var created strfmt.DateTime
	err := t.connectionDB.QueryRow(queryUpdate, message, id).Scan(&p.Id, &p.Parent, &p.Author, &p.Message, &p.IsEdited,
																  &p.Forum, &p.Thread, &created)
	if err != nil {
		return nil, customErrors.PostNotFound
	}

	p.Created = created.String()

	return p, nil
}
