package repository

import (
	"database/sql"
	"fmt"
	"log"
	customErrors "technopark-dbms-forum/internal/pkg/common/custom_errors"
	"technopark-dbms-forum/internal/pkg/post"
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

func (t *ThreadRepository) CreateThreadPosts(forumSlug string, threadId int, posts *[]post.RequestBody) (*[]post.Post, error) {
	isFirst := true
	queryInsert := `INSERT INTO Post (parent, author, message, forum, thread) VALUES `
	querySelect := `SELECT id, author, created, forum, message, thread 
					FROM Post
					WHERE thread = $1 and forum = $2;`
	selection := make([]post.Post, 0)

	for _, value := range *posts {
		if isFirst {
			queryInsert += fmt.Sprintf(`(%d, '%s', '%s', '%s', %d)`, value.Parent, value.Author, value.Message, forumSlug, threadId)
			isFirst = false
			continue
		}
		queryInsert += fmt.Sprintf(`, (%d, '%s', '%s', '%s', %d)`, value.Parent, value.Author, value.Message, forumSlug, threadId)
	}
	queryInsert += fmt.Sprint(";")

	if _, err := t.connectionDB.Exec(queryInsert); err != nil {
		log.Println("queryInsert")
		log.Println(err)
		return nil, customErrors.IncorrectInputData
	}

	querySelectResult, err := t.connectionDB.Query(querySelect, threadId, forumSlug)
	if err != nil || querySelectResult.Err() != nil {
		log.Println("querySelectResult")
		log.Println(err)
		return nil, customErrors.IncorrectInputData
	}

	for querySelectResult.Next() {
		p := new(post.Post)
		_ = querySelectResult.Scan(&p.Id, &p.Author, &p.Created, &p.Forum, &p.Message, &p.Thread)
		selection = append(selection, *p)
	}

	return &selection, nil
}

func (t *ThreadRepository) GetThread(forumSlug string, threadId int, threadSlug string) (*thread.Thread, error) {
	var querySelect string
	var selection *sql.Row

	if len(forumSlug) != 0 {
		querySelect = `	SELECT *
						FROM Thread 
						WHERE forum = $1;`
		selection = t.connectionDB.QueryRow(querySelect, forumSlug)
	} else if len(threadSlug) != 0 {
		querySelect = `	SELECT *
						FROM Thread 
						WHERE slug = $1 or forum = $1;`
		selection = t.connectionDB.QueryRow(querySelect, threadSlug)
	} else {
		querySelect = `	SELECT *
						FROM Thread 
						WHERE id = $1;`
		selection = t.connectionDB.QueryRow(querySelect, threadId)
	}

	if selection == nil {
		return nil, customErrors.ThreadSlugNotFound
	}

	th := new(thread.Thread)
	if err := selection.Scan(&th.Id, &th.Title, &th.Author, &th.Forum, &th.Message, &th.Votes, &th.Slug, &th.Created); err != nil {
		return nil, customErrors.ThreadSlugNotFound
	}

	return th, nil
}