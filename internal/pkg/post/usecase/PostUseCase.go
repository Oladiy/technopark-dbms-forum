package usecase

import (
	customErrors "technopark-dbms-forum/internal/pkg/common/custom_errors"
	"technopark-dbms-forum/internal/pkg/forum"
	"technopark-dbms-forum/internal/pkg/post"
	"technopark-dbms-forum/internal/pkg/post/models"
	"technopark-dbms-forum/internal/pkg/thread"
	"technopark-dbms-forum/internal/pkg/user"

)

type PostUseCase struct {
	ForumRepository forum.Repository
	PostRepository post.Repository
	ThreadRepository thread.Repository
	UserRepository 	user.Repository
}

func NewPostUseCase(forumRepository forum.Repository, postRepository post.Repository, threadRepository thread.Repository, userRepository user.Repository) *PostUseCase {
	return &PostUseCase {
		ForumRepository: forumRepository,
		PostRepository: postRepository,
		ThreadRepository: threadRepository,
		UserRepository:  userRepository,
	}
}

func (t *PostUseCase) GetPost(id int, isUserRelated bool, isThreadRelated bool, isForumRelated bool) (*models.Details, error) {
	details := new(models.Details)
	p, err := t.PostRepository.GetPost(id)
	if err != nil {
		return nil, customErrors.PostNotFound
	}
	details.Post = p

	if isUserRelated {
		details.User, _ = t.UserRepository.GetUserProfile(details.Post.Author)
	}

	if isForumRelated {
		details.Forum, _ = t.ForumRepository.GetForumDetails(details.Post.Forum)
	}

	if isThreadRelated {
		details.Thread, _ = t.ThreadRepository.GetThread("", details.Post.Thread, "")
	}

	return details, nil
}

func (t *PostUseCase) UpdatePost(id int, message string) (*models.Post, error) {
	p, err := t.PostRepository.GetPost(id)
	if p == nil || err != nil {
		return nil, customErrors.PostNotFound
	}

	if len(message) == 0 || message == p.Message {
		return p, nil
	}

	return t.PostRepository.UpdatePost(id, message)
}
