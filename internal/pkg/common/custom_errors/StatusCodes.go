package custom_errors

import (
	"net/http"
)

var (
	StatusCodes = map[error]int {
		DatabaseError: 			 http.StatusInsufficientStorage,

		ForumAlreadyExists: 	 http.StatusConflict,
		ForumSlugNotFound: 		 http.StatusNotFound,
		ForumUserNotFound: 		 http.StatusNotFound,

		IncorrectInputData:      http.StatusBadRequest,

		PostNotFound: 			 http.StatusNotFound,

		ThreadAlreadyExists:	 http.StatusConflict,
		ThreadForumNotFound: 	 http.StatusNotFound,
		ThreadParentConflict: 	 http.StatusConflict,
		ThreadParentNotFound: 	 http.StatusConflict,
		ThreadSlugNotFound:		 http.StatusNotFound,
		ThreadUserNotFound: 	 http.StatusNotFound,

		UserAlreadyExists:       http.StatusConflict,
		UserNotFound:            http.StatusNotFound,
		UserProfileDataConflict: http.StatusConflict,
	}
)
