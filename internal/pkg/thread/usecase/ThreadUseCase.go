package usecase

import (
	customErrors "technopark-dbms-forum/internal/pkg/common/custom_errors"
	"technopark-dbms-forum/internal/pkg/forum"
	"technopark-dbms-forum/internal/pkg/thread"
	"technopark-dbms-forum/internal/pkg/user"
)

type ThreadUseCase struct {
	ThreadRepository thread.Repository
	ForumRepository forum.Repository
	UserRepository 	user.Repository
}

func NewForumUseCase(threadRepository thread.Repository, forumRepository forum.Repository, userRepository user.Repository) *ThreadUseCase {
	return &ThreadUseCase {
		ThreadRepository: threadRepository,
		ForumRepository: forumRepository,
		UserRepository:  userRepository,
	}
}

func (t *ThreadUseCase) CreateThread(slug string, requestBody *thread.RequestBody) (*thread.Thread, error) {
	u, err := t.UserRepository.GetUserProfile(requestBody.Author)
	if err != nil || u == nil {
		return nil, customErrors.ThreadUserNotFound
	}
	requestBody.Author = u.Nickname

	f, err := t.ForumRepository.GetForumDetails(slug)
	if err != nil || f == nil {
		return nil, customErrors.ThreadForumNotFound
	}
	slug = f.Slug

	if len(requestBody.Title) == 0 || len(requestBody.Author) == 0 || len(requestBody.Message) == 0 {
		return nil, customErrors.IncorrectInputData
	}

	return t.ThreadRepository.CreateThread(slug, requestBody)
}
