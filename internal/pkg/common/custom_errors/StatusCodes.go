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
		ThreadParentNotFound: 	 http.StatusConflict,
		ThreadSlugNotFound:		 http.StatusNotFound,
		ThreadUserNotFound: 	 http.StatusNotFound,

		UserAlreadyExists:       http.StatusConflict,
		UserNotFound:            http.StatusNotFound,
		UserProfileDataConflict: http.StatusConflict,
	}
)
