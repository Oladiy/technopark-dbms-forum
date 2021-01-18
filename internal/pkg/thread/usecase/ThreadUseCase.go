package usecase

import (
	customErrors "technopark-dbms-forum/internal/pkg/common/custom_errors"
	"technopark-dbms-forum/internal/pkg/post"
	"technopark-dbms-forum/internal/pkg/thread"
	"technopark-dbms-forum/internal/pkg/vote"
)

type ThreadUseCase struct {
	ThreadRepository thread.Repository
}

func NewThreadUseCase(threadRepository thread.Repository) *ThreadUseCase {
	return &ThreadUseCase {
		ThreadRepository: threadRepository,
	}
}

func (t *ThreadUseCase) CreateThreadPosts(threadSlug string, threadId int, posts *[]post.RequestBody) (*[]post.Post, error) {
	th, err := t.GetThread("", threadId, threadSlug)
	if err != nil {
		return nil, customErrors.ThreadSlugNotFound
	}

	if len(*posts) == 0 {
		emptyResult := make([]post.Post, 0)
		return &emptyResult, nil
	}

	forumSlug := th.Forum
	threadId = th.Id
	return t.ThreadRepository.CreateThreadPosts(forumSlug, threadId, posts)
}

func (t *ThreadUseCase) GetThread(forumSlug string, threadId int, threadSlug string) (*thread.Thread, error) {
	return t.ThreadRepository.GetThread(forumSlug, threadId, threadSlug)
}

func (t *ThreadUseCase) ThreadVote(threadId int, threadSlug string, userVote *vote.Vote) (*thread.Thread, error) {
	_, err := t.GetThread("", threadId, threadSlug)
	if err != nil {
		return nil, customErrors.ThreadSlugNotFound
	}

	return t.ThreadRepository.ThreadVote(threadId, threadSlug, userVote)
}
