package custom_errors

import (
	"net/http"
)

var (
	StatusCodes = map[error]int {
		ForumAlreadyExists: 	 http.StatusConflict,
		ForumSlugNotFound: 		 http.StatusNotFound,
		ForumUserNotFound: 		 http.StatusNotFound,

		IncorrectInputData:      http.StatusBadRequest,

		ThreadAlreadyExists:	 http.StatusConflict,
		ThreadForumNotFound: 	 http.StatusNotFound,
		ThreadUserNotFound: 	 http.StatusNotFound,

		UserNotFound:            http.StatusNotFound,
		UserProfileDataConflict: http.StatusConflict,
		UserAlreadyExists:       http.StatusConflict,
	}
)
