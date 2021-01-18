package custom_errors

import (
	"errors"
)

var (
	ForumAlreadyExists 		= errors.New("forum already exists")
	ForumSlugNotFound 		= errors.New("can't find slug")
	ForumUserNotFound 		= errors.New("can't find user")

	IncorrectInputData      = errors.New("incorrect input data")

	ThreadAlreadyExists		= errors.New("thread already exists")
	ThreadForumNotFound		= errors.New("can't find forum")
	ThreadParentNotFound	= errors.New("can't find parent post")
	ThreadSlugNotFound		= errors.New("can't find thread")
	ThreadUserNotFound 		= errors.New("can't find user")

	UserAlreadyExists       = errors.New("user already exists")
	UserNotFound            = errors.New("can't find user")
	UserProfileDataConflict = errors.New("profile data conflict")
)
