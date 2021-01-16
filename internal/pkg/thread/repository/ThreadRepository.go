package repository

import (
	"database/sql"
	"log"
	customErrors "technopark-dbms-forum/internal/pkg/common/custom_errors"
	"technopark-dbms-forum/internal/pkg/thread"
)

type ThreadRepository struct {
	connectionDB *sql.DB
}

func NewThreadRepository(connectionDB *sql.DB) *ThreadRepository {
	return &ThreadRepository {
		connectionDB: connectionDB,
	}
}

func (t *ThreadRepository) CreateThread(slug string, requestBody *thread.RequestBody) (*thread.Thread, error) {
	var queryInsert string
	var querySelect string
	var row *sql.Row
	var err error

	lengthRBCreated := len(requestBody.Created)
	lengthRBSlug := len(requestBody.Slug)
	th := new(thread.Thread)

	if lengthRBCreated != 0 && lengthRBSlug != 0 {
		log.Println(1)
		queryInsert = "INSERT INTO Thread (title, author, forum, message, slug, created) VALUES($1, $2, $3, $4, $5, $6) returning id, title, author, forum, message, slug, created;"
		querySelect = "SELECT id, title, author, forum, message, slug, created FROM Thread WHERE slug = $1;"
		row = t.connectionDB.QueryRow(queryInsert, requestBody.Title, requestBody.Author, requestBody.Forum, requestBody.Message, requestBody.Slug, requestBody.Created)
		err = row.Scan(&th.Id, &th.Title, &th.Author, &th.Forum, &th.Message, &th.Slug, &th.Created)
	} else if lengthRBCreated != 0 {
		log.Println(2)
		queryInsert = "INSERT INTO Thread (title, author, forum, message, slug, created) VALUES($1, $2, $3, $4, $5, $6) returning id, title, author, forum, message, created;"
		querySelect = "SELECT id, title, author, forum, message, created FROM Thread WHERE slug = $1;"
		row = t.connectionDB.QueryRow(queryInsert, requestBody.Title, requestBody.Author, requestBody.Forum, requestBody.Message, slug, requestBody.Created)
		err = row.Scan(&th.Id, &th.Title, &th.Author, &th.Forum, &th.Message, &th.Created)
	} else {
		log.Println(3)
		queryInsert = "INSERT INTO Thread (title, author, forum, message, slug) VALUES($1, $2, $3, $4, $5) returning id, title, author, forum, message, slug, created;"
		querySelect = "SELECT id, title, author, forum, message, slug FROM Thread WHERE slug = $1;"
		row = t.connectionDB.QueryRow(queryInsert, requestBody.Title, requestBody.Author, requestBody.Forum, requestBody.Message, requestBody.Slug)
		err = row.Scan(&th.Id, &th.Title, &th.Author, &th.Forum, &th.Message, &th.Slug, &th.Created)
	}

	if err != nil {
		log.Println(err)
		var selection *sql.Row

		if len(requestBody.Created) != 0 {
			selection = t.connectionDB.QueryRow(querySelect, slug)
		} else {
			selection = t.connectionDB.QueryRow(querySelect, requestBody.Slug)
		}

		if lengthRBCreated != 0 && lengthRBSlug != 0 {
			err = selection.Scan(&th.Id, &th.Title, &th.Author, &th.Forum, &th.Message, &th.Slug, &th.Created)
		} else if len(requestBody.Created) != 0 {
			err = selection.Scan(&th.Id, &th.Title, &th.Author, &th.Forum, &th.Message, &th.Created)
		} else {
			err = selection.Scan(&th.Id, &th.Title, &th.Author, &th.Forum, &th.Message, &th.Slug)
		}

		if selection.Err() != nil || err != nil {
			return nil, customErrors.IncorrectInputData
		}

		return th, customErrors.ForumAlreadyExists
	}
	// log.Println(th.Id)
	// log.Println(th.Author)
	log.Println(th.Slug)
	log.Println(th.Created)

	return th, nil
}
