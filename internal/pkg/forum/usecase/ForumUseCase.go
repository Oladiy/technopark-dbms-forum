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

func (t *ForumUseCase) GetForumThreads(slug string, limit int, since string, desc bool) (*[]thread.Thread, error) {
	if len(slug) == 0 {
		return nil, customErrors.IncorrectInputData
	}

	_, err := t.ForumRepository.GetForumDetails(slug)
	if err == customErrors.ForumSlugNotFound {
		return nil, customErrors.ForumSlugNotFound
	}

	return t.ForumRepository.GetForumThreads(slug, limit, since, desc)
}

func (t *ForumUseCase) CreateForumThread(slug string, requestBody *thread.RequestBody) (*thread.Thread, error) {
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
	requestBody.Forum = slug

	if len(requestBody.Title) == 0 || len(requestBody.Author) == 0 || len(requestBody.Message) == 0 {
		return nil, customErrors.IncorrectInputData
	}

	th, err := t.ThreadRepository.GetThread("", 0, requestBody.Slug)
	if th != nil {
		return th, customErrors.ThreadAlreadyExists
	}

	return t.ForumRepository.CreateForumThread(slug, requestBody)
}

func (t *ForumUseCase) GetForumUsers(slug string, limit int, since string, desc bool) (*[]user.User, error) {
	_, err := t.GetForumDetails(slug)
	if err != nil {
		return nil, customErrors.ForumSlugNotFound
	}

	return t.ForumRepository.GetForumUsers(slug, limit, since, desc)
}
