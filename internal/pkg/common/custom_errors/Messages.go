package custom_errors

import (
	"errors"
)

var (
	DatabaseError			= errors.New("something wrong with database")

	ForumAlreadyExists 		= errors.New("forum already exists")
	ForumSlugNotFound 		= errors.New("can't find slug")
	ForumUserNotFound 		= errors.New("can't find user")

	IncorrectInputData      = errors.New("incorrect input data")

	PostNotFound			= errors.New("can't find post")

	ThreadAlreadyExists		= errors.New("thread already exists")
	ThreadForumNotFound		= errors.New("can't find forum")
	ThreadParentConflict	= errors.New("parent post was created in another thread")
	ThreadParentNotFound	= errors.New("can't find parent post")
	ThreadSlugNotFound		= errors.New("can't find thread")
	ThreadUserNotFound 		= errors.New("can't find user")

	UserAlreadyExists       = errors.New("user already exists")
	UserNotFound            = errors.New("can't find user")
	UserProfileDataConflict = errors.New("profile data conflict")
)
