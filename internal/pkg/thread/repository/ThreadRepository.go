package repository

import (
	"database/sql"
	"fmt"
	"strings"
	customErrors "technopark-dbms-forum/internal/pkg/common/custom_errors"
	"technopark-dbms-forum/internal/pkg/post"
	"technopark-dbms-forum/internal/pkg/thread"
	"technopark-dbms-forum/internal/pkg/vote"
	"time"
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
	queryInsert := `INSERT INTO Post (parent, author, message, forum, thread, created) VALUES %s`
	selection := make([]post.Post, 0)

	values := make([]interface{}, 0)
	queries := make([]string, 0)
	createdTime := time.Now()
	counter := 0
	for _, value := range *posts {
		queries = append(queries, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)\n",
			counter * 6 + 1, counter * 6 + 2, counter * 6 + 3, counter * 6 + 4, counter * 6 + 5, counter * 6 + 6))
		values = append(values, value.Parent, value.Author, value.Message, forumSlug, threadId, createdTime)
		counter++
	}
	queryInsert = fmt.Sprintf(queryInsert, strings.Join(queries, ","))
	queryInsert += " RETURNING id, created;"

	queryInsertResult, err := t.connectionDB.Query(queryInsert, values...)
	if err != nil {
		return nil, customErrors.IncorrectInputData
	}

	for index := 0; queryInsertResult.Next(); index++ {
		p := new(post.Post)
		_ = queryInsertResult.Scan(&p.Id, &p.Created)
		p.Author = (*posts)[index].Author
		p.Message = (*posts)[index].Message
		p.Parent = (*posts)[index].Parent
		p.Forum = forumSlug
		p.Thread = threadId
		selection = append(selection, *p)
	}

	return &selection, nil
}

func (t *ThreadRepository) GetThread(forumSlug string, threadId int, threadSlug string) (*thread.Thread, error) {
	var querySelect string
	var selection *sql.Row

	if len(forumSlug) != 0 {
		querySelect = `	SELECT id, title, author, forum, message, votes, slug, created
						FROM Thread 
						WHERE forum = $1;`
		selection = t.connectionDB.QueryRow(querySelect, forumSlug)
	} else if len(threadSlug) != 0 {
		querySelect = `	SELECT id, title, author, forum, message, votes, slug, created 
						FROM Thread 
						WHERE slug = $1;`
		selection = t.connectionDB.QueryRow(querySelect, threadSlug)
	} else {
		querySelect = `	SELECT id, title, author, forum, message, votes, slug, created
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

func (t *ThreadRepository) ThreadVote(threadId int, threadSlug string, userVote *vote.Vote) (*thread.Thread, error) {
	var querySelectThread string
	var queryUpdateThread string
	var selection *sql.Row
	var err error

	if len(threadSlug) != 0 {
		querySelectThread = `	SELECT author, created, forum, id, message, slug, title, votes 
								FROM Thread 
								WHERE slug = $1;`
	} else {
		querySelectThread = `	SELECT author, created, forum, id, message, slug, title, votes 
								FROM Thread 
								WHERE id = $1;`
	}
	th := new(thread.Thread)

	v := new(vote.Vote)
	querySelectVote := `SELECT nickname, voice
						FROM Vote
						WHERE nickname = $1;`
	selection = t.connectionDB.QueryRow(querySelectVote, userVote.Nickname)
	err = selection.Scan(&v.Nickname, &v.Voice)
	if err == nil  {
		queryUpdateVote := `UPDATE Vote 
							SET voice = $1 
							WHERE nickname = $2;`
		_, err = t.connectionDB.Exec(queryUpdateVote, userVote.Voice, userVote.Nickname)
		if err != nil {
			return nil, customErrors.IncorrectInputData
		}

		if v.Voice == userVote.Voice {
			if len(threadSlug) != 0 {
				selection = t.connectionDB.QueryRow(querySelectThread, threadSlug)
			} else {
				selection = t.connectionDB.QueryRow(querySelectThread, threadId)
			}

			err = selection.Scan(&th.Author, &th.Created, &th.Forum, &th.Id, &th.Message, &th.Slug, &th.Title, &th.Votes)
			if selection.Err() != nil || err != nil {
				return nil, customErrors.ThreadSlugNotFound
			}

			return th, nil
		}

		if len(threadSlug) != 0 {
			if userVote.Voice == 1 {
				queryUpdateThread = `	UPDATE Thread
										SET votes = votes + 2
										WHERE slug = $1;`
			} else {
				queryUpdateThread = `	UPDATE Thread
										SET votes = votes - 2
										WHERE slug = $1;`
			}
		} else {
			if userVote.Voice == 1 {
				queryUpdateThread = `	UPDATE Thread
										SET votes = votes + 2
										WHERE id = $1;`
			} else {
				queryUpdateThread = `	UPDATE Thread
										SET votes = votes - 2
										WHERE id = $1;`
			}
		}

		if len(threadSlug) != 0 {
			if _, err = t.connectionDB.Exec(queryUpdateThread, threadSlug); err != nil {
				return nil, customErrors.IncorrectInputData
			}
			selection = t.connectionDB.QueryRow(querySelectThread, threadSlug)
		} else {
			if _, err = t.connectionDB.Exec(queryUpdateThread, threadId); err != nil {
				return nil, customErrors.IncorrectInputData
			}
			selection = t.connectionDB.QueryRow(querySelectThread, threadId)
		}

		err = selection.Scan(&th.Author, &th.Created, &th.Forum, &th.Id, &th.Message, &th.Slug, &th.Title, &th.Votes)
		if selection.Err() != nil || err != nil {
			return nil, customErrors.ThreadSlugNotFound
		}

		return th, nil
	}

	if len(threadSlug) != 0 {
		if userVote.Voice == 1 {
			queryUpdateThread = `	UPDATE Thread
									SET votes = votes + 1
									WHERE slug = $1;`
		} else {
			queryUpdateThread = `	UPDATE Thread
									SET votes = votes - 1
									WHERE slug = $1;`
		}
	} else {
		if userVote.Voice == 1 {
			queryUpdateThread = `	UPDATE Thread
									SET votes = votes + 1
									WHERE id = $1;`
		} else {
			queryUpdateThread = `	UPDATE Thread
									SET votes = votes - 1
									WHERE id = $1;`
		}
	}

	queryInsertVote := `INSERT INTO Vote (nickname, voice) 
						VALUES ($1, $2);`

	_, err = t.connectionDB.Exec(queryInsertVote, userVote.Nickname, userVote.Voice)
	if err != nil {
		return nil, customErrors.IncorrectInputData
	}

	if len(threadSlug) != 0 {
		if _, err = t.connectionDB.Exec(queryUpdateThread, threadSlug); err != nil {
			return nil, customErrors.IncorrectInputData
		}
		selection = t.connectionDB.QueryRow(querySelectThread, threadSlug)
	} else {
		if _, err = t.connectionDB.Exec(queryUpdateThread, threadId); err != nil {
			return nil, customErrors.IncorrectInputData
		}
		selection = t.connectionDB.QueryRow(querySelectThread, threadId)
	}

	err = selection.Scan(&th.Author, &th.Created, &th.Forum, &th.Id, &th.Message, &th.Slug, &th.Title, &th.Votes)
	if selection.Err() != nil || err != nil {
		return nil, customErrors.ThreadSlugNotFound
	}

	return th, nil
}
