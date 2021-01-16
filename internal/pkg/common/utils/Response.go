package utils

import (
	"encoding/json"
	"net/http"
	customErrors "technopark-dbms-forum/internal/pkg/common/custom_errors"
)

func MakeErrorResponse(w http.ResponseWriter, err error) {
	output, _ := json.Marshal(customErrors.Response {
		Message: err.Error(),
	})
	_, _ = w.Write(output)
}
