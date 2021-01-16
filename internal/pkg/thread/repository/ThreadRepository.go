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
		queryInsert = `	INSERT INTO Thread (title, author, forum, message, slug, created) 
						VALUES($1, $2, $3, $4, $5, $6) 
						RETURNING id, title, author, forum, message, slug, created;`
		querySelect = `	SELECT id, title, author, forum, message, slug, created 
						FROM Thread 
						WHERE slug = $1;`
		row = t.connectionDB.QueryRow(queryInsert, requestBody.Title, requestBody.Author, requestBody.Forum, requestBody.Message, requestBody.Slug, requestBody.Created)
		err = row.Scan(&th.Id, &th.Title, &th.Author, &th.Forum, &th.Message, &th.Slug, &th.Created)
	} else if lengthRBCreated != 0 {
		queryInsert = `	INSERT INTO Thread (title, author, forum, message, slug, created) 
						VALUES($1, $2, $3, $4, $5, $6) 
						RETURNING id, title, author, forum, message, created;`
		querySelect = `	SELECT id, title, author, forum, message, created 
						FROM Thread 
						WHERE slug = $1;`
		row = t.connectionDB.QueryRow(queryInsert, requestBody.Title, requestBody.Author, requestBody.Forum, requestBody.Message, slug, requestBody.Created)
		err = row.Scan(&th.Id, &th.Title, &th.Author, &th.Forum, &th.Message, &th.Created)
	} else {
		queryInsert = `	INSERT INTO Thread (title, author, forum, message, slug) 
						VALUES($1, $2, $3, $4, $5) 
						RETURNING id, title, author, forum, message, slug, created;`
		querySelect = `	SELECT id, title, author, forum, message, slug 
						FROM Thread 
						WHERE slug = $1;`
		row = t.connectionDB.QueryRow(queryInsert, requestBody.Title, requestBody.Author, requestBody.Forum, requestBody.Message, requestBody.Slug)
		err = row.Scan(&th.Id, &th.Title, &th.Author, &th.Forum, &th.Message, &th.Slug, &th.Created)
	}

	if err != nil {
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

	return th, nil
}

func (t *ThreadRepository) GetThreadList(slug string, limit int, since string, desc bool) (*[]thread.Thread, error) {
	var querySelect string
	var querySelectResult *sql.Rows
	var err error

	selection := make([]thread.Thread, 0)
	lengthSince := len(since)

	if desc {
		log.Println(desc)
		if lengthSince != 0 && limit != 0 {
			querySelect = `	SELECT id, title, author, forum, message, slug, votes, created 
							FROM Thread 
							WHERE forum = $1 and created <= $2
							GROUP BY id, created
							ORDER BY created DESC
							LIMIT $3;`
			querySelectResult, err = t.connectionDB.Query(querySelect, slug, since, limit)
		} else if lengthSince != 0 {
			querySelect = `	SELECT id, title, author, forum, message, slug, votes, created 
							FROM Thread 
							WHERE forum = $1 and created <= $2
							ORDER BY created DESC;`
			querySelectResult, err = t.connectionDB.Query(querySelect, slug, since)
		} else if limit != 0 {
			querySelect = `	SELECT id, title, author, forum, message, slug, votes, created 
							FROM Thread 
							WHERE forum = $1
							ORDER BY created DESC
							LIMIT $2;`
			querySelectResult, err = t.connectionDB.Query(querySelect, slug, limit)
		} else {
			querySelect = `	SELECT id, title, author, forum, message, slug, votes, created 
							FROM Thread 
							WHERE forum = $1
							ORDER BY created DESC;`
			querySelectResult, err = t.connectionDB.Query(querySelect, slug)
		}
	} else {
		log.Println(desc)
		if lengthSince != 0 && limit != 0 {
			querySelect = `	SELECT id, title, author, forum, message, slug, votes, created 
							FROM Thread 
							WHERE forum = $1 and created >= $2
							ORDER BY created
							LIMIT $3;`
			querySelectResult, err = t.connectionDB.Query(querySelect, slug, since, limit)
		} else if lengthSince != 0 {
			querySelect = `	SELECT id, title, author, forum, message, slug, votes, created 
							FROM Thread 
							WHERE forum = $1 and created >= $2
							ORDER BY created;`
			querySelectResult, err = t.connectionDB.Query(querySelect, slug, since)
		} else if limit != 0 {
			querySelect = `	SELECT id, title, author, forum, message, slug, votes, created 
							FROM Thread 
							WHERE forum = $1
							GROUP BY id, created
							ORDER BY created
							LIMIT $2;`
			querySelectResult, err = t.connectionDB.Query(querySelect, slug, limit)
		} else {
			querySelect = `	SELECT id, title, author, forum, message, slug, votes, created 
							FROM Thread 
							WHERE forum = $1
							ORDER BY created;`
			querySelectResult, err = t.connectionDB.Query(querySelect, slug)
		}
	}

	if err != nil || querySelectResult.Err() != nil {
		return nil, customErrors.IncorrectInputData
	}

	for querySelectResult.Next() {
		t := new(thread.Thread)
		_ = querySelectResult.Scan(&t.Id, &t.Title, &t.Author, &t.Forum, &t.Message, &t.Slug, &t.Votes, &t.Created)
		selection = append(selection, *t)
	}
	err = querySelectResult.Close()
	if err != nil {
		return nil, err
	}

	return &selection, nil
}
