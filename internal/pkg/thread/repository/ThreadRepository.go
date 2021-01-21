package repository

import (
	"database/sql"
	"fmt"
	"strings"
	customErrors "technopark-dbms-forum/internal/pkg/common/custom_errors"
	postModel "technopark-dbms-forum/internal/pkg/post/models"
	threadModel "technopark-dbms-forum/internal/pkg/thread/models"
	voteModel "technopark-dbms-forum/internal/pkg/vote/models"
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

func (t *ThreadRepository) CreateThreadPosts(forumSlug string, threadId int, posts *[]postModel.RequestBody) (*[]postModel.Post, error) {
	queryInsert := `INSERT INTO Post (parent, author, message, forum, thread, created) VALUES %s`
	selection := make([]postModel.Post, 0)

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
		if err.Error() == `pq: insert or update on table "forumusers" violates foreign key constraint "forumusers_user_nickname_fkey"` {
			return nil, customErrors.PostNotFound
		}
		if queryInsertResult == nil {
			return nil, customErrors.ThreadParentConflict
		}

		return nil, customErrors.IncorrectInputData
	}

	for index := 0; queryInsertResult.Next(); index++ {
		p := new(postModel.Post)
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

func (t *ThreadRepository) GetThread(forumSlug string, threadId int, threadSlug string) (*threadModel.Thread, error) {
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

	th := new(threadModel.Thread)
	if err := selection.Scan(&th.Id, &th.Title, &th.Author, &th.Forum, &th.Message, &th.Votes, &th.Slug, &th.Created); err != nil {
		return nil, customErrors.ThreadSlugNotFound
	}

	return th, nil
}

func (t *ThreadRepository) ThreadVote(threadId int, userVote *voteModel.Vote) (*threadModel.Thread, error) {
	var queryUpdateThread string
	var selection *sql.Row
	var err error

	querySelectThread := `	SELECT author, created, forum, id, message, slug, title, votes 
							FROM Thread 
							WHERE id = $1;`
	th := new(threadModel.Thread)

	v := new(voteModel.Vote)
	querySelectVote := `SELECT nickname, voice
						FROM Vote
						WHERE nickname = $1 AND thread = $2;`
	selection = t.connectionDB.QueryRow(querySelectVote, userVote.Nickname, threadId)
	err = selection.Scan(&v.Nickname, &v.Voice)
	if err == nil  {
		queryUpdateVote := `UPDATE Vote 
							SET voice = $1 
							WHERE nickname = $2 AND thread = $3;`
		_, err = t.connectionDB.Exec(queryUpdateVote, userVote.Voice, userVote.Nickname, threadId)
		if err != nil {
			return nil, customErrors.IncorrectInputData
		}

		if v.Voice == userVote.Voice {
			selection = t.connectionDB.QueryRow(querySelectThread, threadId)

			err = selection.Scan(&th.Author, &th.Created, &th.Forum, &th.Id, &th.Message, &th.Slug, &th.Title, &th.Votes)
			if err != nil {
				return nil, customErrors.ThreadSlugNotFound
			}

			return th, nil
		}

		if userVote.Voice == 1 {
			queryUpdateThread = `	UPDATE Thread
									SET votes = votes + 2
									WHERE id = $1;`
		} else {
			queryUpdateThread = `	UPDATE Thread
									SET votes = votes - 2
									WHERE id = $1;`
		}

		if _, err = t.connectionDB.Exec(queryUpdateThread, threadId); err != nil {
			return nil, customErrors.IncorrectInputData
		}
		selection = t.connectionDB.QueryRow(querySelectThread, threadId)

		err = selection.Scan(&th.Author, &th.Created, &th.Forum, &th.Id, &th.Message, &th.Slug, &th.Title, &th.Votes)
		if err != nil {
			return nil, customErrors.ThreadSlugNotFound
		}

		return th, nil
	}

	if userVote.Voice == 1 {
		queryUpdateThread = `	UPDATE Thread
								SET votes = votes + 1
								WHERE id = $1;`
	} else {
		queryUpdateThread = `	UPDATE Thread
								SET votes = votes - 1
								WHERE id = $1;`
	}

	queryInsertVote := `INSERT INTO Vote (nickname, voice, thread) 
						VALUES ($1, $2, $3);`

	_, err = t.connectionDB.Exec(queryInsertVote, userVote.Nickname, userVote.Voice, threadId)
	if err != nil {
		return nil, customErrors.IncorrectInputData
	}

	if _, err = t.connectionDB.Exec(queryUpdateThread, threadId); err != nil {
		return nil, customErrors.IncorrectInputData
	}
	selection = t.connectionDB.QueryRow(querySelectThread, threadId)

	err = selection.Scan(&th.Author, &th.Created, &th.Forum, &th.Id, &th.Message, &th.Slug, &th.Title, &th.Votes)
	if err != nil {
		return nil, customErrors.ThreadSlugNotFound
	}

	return th, nil
}

func (t *ThreadRepository) UpdateThread(threadId int, threadToUpdate *threadModel.RequestBody) (*threadModel.Thread, error) {
	var selection *sql.Row
	isFirst := true
	fieldsToUpdate := make([]interface{}, 0)

	queryUpdate := `UPDATE Thread
					SET `

	counter := 1
	if len(threadToUpdate.Title) != 0 {
		queryUpdate += fmt.Sprintf("title = $%d ", counter)
		isFirst = false
		fieldsToUpdate = append(fieldsToUpdate, threadToUpdate.Title)
		counter++
	}

	if len(threadToUpdate.Message) != 0 {
		if isFirst {
			queryUpdate += fmt.Sprintf("message = $%d ", counter)
		} else {
			queryUpdate += fmt.Sprintf(", message = $%d ", counter)
		}

		counter++
		fieldsToUpdate = append(fieldsToUpdate, threadToUpdate.Message)
	}
	queryUpdate += fmt.Sprintf("WHERE id = $%d ", counter)

	fieldsToUpdate = append(fieldsToUpdate, threadId)
	queryUpdate += `RETURNING author, created, forum, id, message, slug, title, votes;`
	selection = t.connectionDB.QueryRow(queryUpdate, fieldsToUpdate...)

	th := new(threadModel.Thread)
	err := selection.Scan(&th.Author, &th.Created, &th.Forum, &th.Id, &th.Message, &th.Slug, &th.Title, &th.Votes)
	if err != nil {
		return nil, customErrors.IncorrectInputData
	}

	return th, nil
}

func (t *ThreadRepository) GetThreadPosts(threadId int, limit int, since int, sort string, desc bool) (*[]postModel.Post, error) {
	querySelect := `SELECT id, parent, author, message, isEdited, forum, thread, created
					FROM Post
					WHERE thread = $1 `
	switch sort {
	case "tree":
		if since > 0 {
			if desc {
				querySelect += `AND path < (SELECT path FROM Post WHERE id = $2)
								ORDER BY path DESC `
			} else {
				querySelect += `AND path > (SELECT path FROM Post WHERE id = $2)
								ORDER BY path `
			}
		} else {
			if desc {
				querySelect += `AND path[1] > $2
							ORDER BY path DESC `
			} else {
				querySelect += `AND path[1] > $2
							ORDER BY path `
			}
		}
		querySelect += "LIMIT $3 "
		break
	case "parent_tree":
		if since > 0 {
			if desc {
				querySelect += `AND path[1] 
								IN (SELECT path[1] 
									FROM Post 
									WHERE thread = $1 AND parent = 0 AND path[1] < (SELECT path[1] 
																					FROM Post 
																					WHERE id = $2)
									ORDER BY path[1] DESC 
									LIMIT $3)
								ORDER BY path[1] DESC, path`
			} else {
				querySelect += `AND path[1] 
								IN (SELECT path[1] 
									FROM Post 
									WHERE thread = $1 AND parent = 0 AND path[1] > (SELECT path[1] 
																					FROM Post 
																					WHERE id = $2)
									ORDER BY path[1] 
									LIMIT $3)
								ORDER BY path[1], path`
			}
		} else {
			if desc {
				querySelect += `AND path[1] > $2 AND path[1] 
								IN (SELECT path[1]
									FROM Post
									WHERE thread = $1 and id > $2 AND parent = 0
									ORDER BY path[1] DESC
									LIMIT $3)
								ORDER BY path[1] DESC, path`
			} else {
				querySelect += `AND path[1] > $2 AND path[1] 
								IN (SELECT path[1]
									FROM Post
									WHERE thread = $1 and id > $2 AND parent = 0
									ORDER BY path[1] 
									LIMIT $3)
								ORDER BY path[1], path`
			}
		}
		break
	default:
		if since > 0 {
			if desc {
				querySelect += `AND id < $2
								ORDER BY created DESC, id DESC `
			} else {
				querySelect += `AND id > $2
								ORDER BY created, id `
			}
		} else {
			if desc {
				querySelect += `AND id > $2
								ORDER BY created DESC, id DESC `
			} else {
				querySelect += `AND id > $2
								ORDER BY created, id `
			}
		}
		querySelect += "LIMIT $3 "
	}

	querySelectResult, err := t.connectionDB.Query(querySelect, threadId, since, limit)
	if err != nil || querySelectResult == nil {
		return nil, customErrors.ThreadSlugNotFound
	}

	selection := make([]postModel.Post, 0)
	for querySelectResult.Next() {
		p := new(postModel.Post)
		err = querySelectResult.Scan(&p.Id, &p.Parent, &p.Author, &p.Message, &p.IsEdited, &p.Forum, &p.Thread, &p.Created)
		if err != nil {
			return nil, err
		}
		selection = append(selection, *p)
	}

	return &selection, nil
}
