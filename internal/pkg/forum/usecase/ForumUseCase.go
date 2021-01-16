package usecase

import (
	"regexp"
	customErrors "technopark-dbms-forum/internal/pkg/common/custom_errors"
	"technopark-dbms-forum/internal/pkg/forum"
	"technopark-dbms-forum/internal/pkg/thread"
	"technopark-dbms-forum/internal/pkg/user"
)

type ForumUseCase struct {
	ForumRepository forum.Repository
	ThreadRepository thread.Repository
	UserRepository 	user.Repository
}

func NewForumUseCase(forumRepository forum.Repository, threadRepository thread.Repository, userRepository user.Repository) *ForumUseCase {
	return &ForumUseCase {
		ForumRepository: forumRepository,
		ThreadRepository: threadRepository,
		UserRepository:  userRepository,
	}
}

func (t *ForumUseCase) CreateForum(requestBody *forum.RequestBody) (*forum.Forum, error) {
	u, err := t.UserRepository.GetUserProfile(requestBody.User)

	if err != nil || u == nil {
		return nil, customErrors.ForumUserNotFound
	}

	requestBody.User = u.Nickname
	regExp := `^(\d|\w|-|_)*(\w|-|_)(\d|\w|-|_)*$`
	match, _ := regexp.MatchString(regExp, requestBody.Slug)
	if len(requestBody.Title) == 0 || !match {
		return nil, customErrors.IncorrectInputData
	}

	return t.ForumRepository.CreateForum(requestBody)
}

func (t *ForumUseCase) GetForumDetails(slug string) (*forum.Forum, error) {
	if len(slug) == 0 {
		return nil, customErrors.IncorrectInputData
	}

	return t.ForumRepository.GetForumDetails(slug)
}

func (t *ForumUseCase) GetForumThreadList(slug string, limit int, since string, desc bool) (*[]thread.Thread, error) {
	if len(slug) == 0 {
		return nil, customErrors.IncorrectInputData
	}

	_, err := t.ForumRepository.GetForumDetails(slug)
	if err == customErrors.ForumSlugNotFound {
		return nil, customErrors.ForumSlugNotFound
	}

	return t.ThreadRepository.GetThreadList(slug, limit, since, desc)
}
