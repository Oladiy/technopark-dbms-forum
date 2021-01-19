package repository

import (
	"database/sql"
	"log"
	customErrors "technopark-dbms-forum/internal/pkg/common/custom_errors"
	"technopark-dbms-forum/internal/pkg/forum"
	"technopark-dbms-forum/internal/pkg/thread"
	"technopark-dbms-forum/internal/pkg/user"
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

func (t *ForumRepository) CreateForumThread(slug string, requestBody *thread.RequestBody) (*thread.Thread, error) {
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

		return th, customErrors.ThreadAlreadyExists
	}

	return th, nil
}

func (t *ForumRepository) GetForumThreads(slug string, limit int, since string, desc bool) (*[]thread.Thread, error) {
	var querySelect string
	var querySelectResult *sql.Rows
	var err error

	selection := make([]thread.Thread, 0)
	lengthSince := len(since)

	if desc {
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

func (t *ForumRepository) GetForumUsers(slug string, limit int, since string, desc bool) (*[]user.User, error) {
	var querySelect string

	log.Println("slug: ", slug)
	log.Println("limit: ", limit)
	log.Println("since: ", since)
	log.Println("desc: ", desc)

	if len(since) != 0 && !desc {
		querySelect = `SELECT nickname, fullname, about, email
					FROM ForumUsers fu JOIN Users u ON(fu.user_nickname = u.nickname)
					WHERE fu.forum_slug = $1 AND u.nickname > $2 COLLATE "C"
					ORDER BY u.nickname COLLATE "C" `
	} else if len(since) != 0 {
		querySelect = `SELECT nickname, fullname, about, email
					FROM ForumUsers fu JOIN Users u ON(fu.user_nickname = u.nickname)
					WHERE fu.forum_slug = $1 AND u.nickname < $2 COLLATE "C"
					ORDER BY u.nickname COLLATE "C" `
	} else {
		querySelect = `SELECT nickname, fullname, about, email
					FROM ForumUsers fu JOIN Users u ON(fu.user_nickname = u.nickname)
					WHERE fu.forum_slug = $1 AND u.nickname > $2 COLLATE "C"
					ORDER BY u.nickname COLLATE "C" `
	}

	if desc {
		querySelect += "DESC "
	}
	querySelect += "LIMIT $3"

	querySelectResult, err := t.connectionDB.Query(querySelect, slug, since, limit)
	if err != nil {
		return nil, customErrors.ForumSlugNotFound
	}

	selection := make([]user.User, 0)

	for querySelectResult.Next() {
		u := new(user.User)
		_ = querySelectResult.Scan(&u.Nickname, &u.FullName, &u.About, &u.Email)
		selection = append(selection, *u)
	}

	querySelectResult.Close()
	return &selection, nil
}
