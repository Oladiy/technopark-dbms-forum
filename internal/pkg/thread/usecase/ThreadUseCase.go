package usecase

import (
	customErrors "technopark-dbms-forum/internal/pkg/common/custom_errors"
	"technopark-dbms-forum/internal/pkg/post/models"
	"technopark-dbms-forum/internal/pkg/thread"
	models2 "technopark-dbms-forum/internal/pkg/thread/models"
	"technopark-dbms-forum/internal/pkg/user"
	models3 "technopark-dbms-forum/internal/pkg/vote/models"
)

type ThreadUseCase struct {
	ThreadRepository thread.Repository
	UserRepository   user.Repository
}

func NewThreadUseCase(threadRepository thread.Repository, userRepository user.Repository) *ThreadUseCase {
	return &ThreadUseCase {
		ThreadRepository: threadRepository,
		UserRepository: userRepository,
	}
}

func (t *ThreadUseCase) CreateThreadPosts(threadSlug string, threadId int, posts *[]models.RequestBody) (*[]models.Post, error) {
	th, err := t.GetThread("", threadId, threadSlug)
	if err != nil {
		return nil, customErrors.ThreadSlugNotFound
	}

	if len(*posts) == 0 {
		emptyResult := make([]models.Post, 0)
		return &emptyResult, nil
	}

	forumSlug := th.Forum
	threadId = th.Id
	return t.ThreadRepository.CreateThreadPosts(forumSlug, threadId, posts)
}

func (t *ThreadUseCase) GetThread(forumSlug string, threadId int, threadSlug string) (*models2.Thread, error) {
	return t.ThreadRepository.GetThread(forumSlug, threadId, threadSlug)
}

func (t *ThreadUseCase) ThreadVote(threadId int, threadSlug string, userVote *models3.Vote) (*models2.Thread, error) {
	_, err := t.GetThread("", threadId, threadSlug)
	if err != nil {
		return nil, customErrors.ThreadSlugNotFound
	}

	u, err := t.UserRepository.GetUserProfile(userVote.Nickname)
	if err != nil || u == nil {
		return nil, customErrors.UserNotFound
	}

	return t.ThreadRepository.ThreadVote(threadId, threadSlug, userVote)
}

func (t *ThreadUseCase) GetThreadPosts(threadId int, threadSlug string, limit int, since int, sort string, desc bool) (*[]models.Post, error) {
	th, err := t.GetThread("", threadId, threadSlug)
	if err != nil {
		return nil, customErrors.ThreadSlugNotFound
	}

	threadId = th.Id

	return t.ThreadRepository.GetThreadPosts(threadId, threadSlug, limit, since, sort, desc)
}

func (t *ThreadUseCase) UpdateThread(threadId int, threadSlug string, threadToUpdate *models2.RequestBody) (*models2.Thread, error) {
	th, err := t.GetThread("", threadId, threadSlug)
	if err != nil {
		return nil, customErrors.ThreadSlugNotFound
	}
	threadId = th.Id

	if len(threadToUpdate.Title) == 0 && len(threadToUpdate.Message) == 0 {
		return th, nil
	}

	return t.ThreadRepository.UpdateThread(threadId, threadToUpdate)
}
