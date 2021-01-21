package repository

import (
	"github.com/go-openapi/strfmt"
	"github.com/jackc/pgx"
	customErrors "technopark-dbms-forum/internal/pkg/common/custom_errors"
	forumModels "technopark-dbms-forum/internal/pkg/forum/models"
	threadModels "technopark-dbms-forum/internal/pkg/thread/models"
	userModels "technopark-dbms-forum/internal/pkg/user/models"
)

type ForumRepository struct {
	connectionDB *pgx.ConnPool
}

func NewForumRepository(connectionDB *pgx.ConnPool) *ForumRepository {
	return &ForumRepository {
		connectionDB: connectionDB,
	}
}

func (t *ForumRepository) CreateForum(requestBody *forumModels.RequestBody) (*forumModels.Forum, error) {
	queryInsert := `INSERT INTO Forum (title, author, slug) 
					VALUES($1, $2, $3) 
					returning title, author, slug;`
	querySelect := `SELECT title, author, slug, posts, threads 
					FROM Forum 
					WHERE slug = $1;`
	f := new(forumModels.Forum)

	row := t.connectionDB.QueryRow(queryInsert, requestBody.Title, requestBody.User, requestBody.Slug)
	err := row.Scan(&f.Title, &f.User, &f.Slug)
	if err != nil {
		selection := t.connectionDB.QueryRow(querySelect, requestBody.Slug)
		err = selection.Scan(&f.Title, &f.User, &f.Slug, &f.Posts, &f.Threads)
		if err != nil {
			return nil, customErrors.IncorrectInputData
		}

		return f, customErrors.ForumAlreadyExists
	}

	return f, nil
}

func (t *ForumRepository) GetForumDetails(slug string) (*forumModels.Forum, error) {
	querySelect := `SELECT title, author, slug, posts, threads 
					FROM Forum 
					WHERE slug = $1;`
	selection := t.connectionDB.QueryRow(querySelect, slug)
	if selection == nil {
		return nil, customErrors.ForumSlugNotFound
	}

	f := new(forumModels.Forum)
	if err := selection.Scan(&f.Title, &f.User, &f.Slug, &f.Posts, &f.Threads); err != nil {
		return nil, customErrors.ForumSlugNotFound
	}

	return f, nil
}

func (t *ForumRepository) CreateForumThread(slug string, requestBody *threadModels.RequestBody) (*threadModels.Thread, error) {
	var created strfmt.DateTime
	var queryInsert string
	var querySelect string
	var row *pgx.Row
	var err error

	lengthRBCreated := len(requestBody.Created)
	lengthRBSlug := len(requestBody.Slug)
	th := new(threadModels.Thread)


	if lengthRBCreated != 0 {
		queryInsert = `	INSERT INTO Thread (title, author, forum, message, slug, created) 
						VALUES($1, $2, $3, $4, $5, $6) 
						RETURNING id, slug, created;`
		querySelect = `	SELECT id, title, author, forum, message, slug, created 
						FROM Thread 
						WHERE slug = $1;`
		row = t.connectionDB.QueryRow(queryInsert, requestBody.Title, requestBody.Author, requestBody.Forum, requestBody.Message, requestBody.Slug, requestBody.Created)
		err = row.Scan(&th.Id, &th.Slug, &created)
	} else {
		queryInsert = `	INSERT INTO Thread (title, author, forum, message, slug) 
						VALUES($1, $2, $3, $4, $5) 
						RETURNING id, slug, created;`
		querySelect = `	SELECT id, title, author, forum, message, slug 
						FROM Thread 
						WHERE slug = $1;`
		row = t.connectionDB.QueryRow(queryInsert, requestBody.Title, requestBody.Author, requestBody.Forum, requestBody.Message, requestBody.Slug)
		err = row.Scan(&th.Id, &th.Slug, &created)
	}

	if err != nil {
		var selection *pgx.Row

		selection = t.connectionDB.QueryRow(querySelect, requestBody.Slug)
		if lengthRBCreated != 0 && lengthRBSlug != 0 {
			err = selection.Scan(&th.Id, &th.Title, &th.Author, &th.Forum, &th.Message, &th.Slug, &created)
		} else {
			err = selection.Scan(&th.Id, &th.Title, &th.Author, &th.Forum, &th.Message, &th.Slug, &created)
		}

		if err != nil {
			return nil, customErrors.IncorrectInputData
		}

		return th, customErrors.ThreadAlreadyExists
	}

	th.Title = requestBody.Title
	th.Message = requestBody.Message
	th.Author = requestBody.Author
	th.Created = created.String()
	th.Forum = slug

	return th, nil
}

func (t *ForumRepository) GetForumThreads(slug string, limit int, since string, desc bool) (*[]threadModels.Thread, error) {
	var querySelect string
	var querySelectResult *pgx.Rows
	var err error

	selection := make([]threadModels.Thread, 0)
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

	if err != nil {
		return nil, customErrors.IncorrectInputData
	}

	var created strfmt.DateTime
	for querySelectResult.Next() {
		th := new(threadModels.Thread)
		_ = querySelectResult.Scan(&th.Id, &th.Title, &th.Author, &th.Forum, &th.Message, &th.Slug, &th.Votes, &created)
		th.Created = created.String()
		selection = append(selection, *th)
	}

	querySelectResult.Close()

	return &selection, nil
}

func (t *ForumRepository) GetForumUsers(slug string, limit int, since string, desc bool) (*[]userModels.User, error) {
	var querySelect string

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

	selection := make([]userModels.User, 0)

	for querySelectResult.Next() {
		u := new(userModels.User)
		_ = querySelectResult.Scan(&u.Nickname, &u.FullName, &u.About, &u.Email)
		selection = append(selection, *u)
	}

	querySelectResult.Close()
	return &selection, nil
}
