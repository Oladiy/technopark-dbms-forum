package usecase

import (
	customErrors "technopark-dbms-forum/internal/pkg/common/custom_errors"
	"technopark-dbms-forum/internal/pkg/post"
	"technopark-dbms-forum/internal/pkg/thread"
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

	isParent := false
	for _, element := range *posts {
		if element.Parent == 0 {
			isParent = true
			break
		}
	}
	if !isParent {
		return nil, customErrors.ThreadParentNotFound
	}

	forumSlug := th.Forum
	threadId = th.Id
	return t.ThreadRepository.CreateThreadPosts(forumSlug, threadId, posts)
}

func (t *ThreadUseCase) GetThread(forumSlug string, threadId int, threadSlug string) (*thread.Thread, error) {
	return t.ThreadRepository.GetThread(forumSlug, threadId, threadSlug)
}
