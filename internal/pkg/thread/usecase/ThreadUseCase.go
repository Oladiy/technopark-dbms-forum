package usecase

import (
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
