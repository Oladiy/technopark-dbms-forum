package delivery

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"technopark-dbms-forum/internal/pkg/common/consts"
	customErrors "technopark-dbms-forum/internal/pkg/common/custom_errors"
	"technopark-dbms-forum/internal/pkg/common/utils"
	"technopark-dbms-forum/internal/pkg/user"
	"technopark-dbms-forum/internal/pkg/user/models"
)

type UserDelivery struct {
	UserUseCase user.UseCase
}

func NewUserDelivery(userUseCase user.UseCase) *UserDelivery {
	return &UserDelivery {
		UserUseCase: userUseCase,
	}
}

func (t *UserDelivery) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	nickname := vars[consts.UserNickNamePath]
	requestBody := new(models.RequestBody)
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		utils.MakeErrorResponse(w, err)
		return
	}

	response, err := t.UserUseCase.CreateUser(nickname, requestBody)
	if err != nil {
		if response != nil {
			w.WriteHeader(customErrors.StatusCodes[err])
			output, _ := json.Marshal(response)
			_, _ = w.Write(output)
		} else {
			utils.MakeErrorResponse(w, err)
		}

		return
	}

	w.WriteHeader(http.StatusCreated)
	output, _ := json.Marshal((*response)[0])
	_, _ = w.Write(output)
}

func (t *UserDelivery) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	nickname := vars[consts.UserNickNamePath]

	response, err := t.UserUseCase.GetUserProfile(nickname)
	if err != nil {
		w.WriteHeader(customErrors.StatusCodes[err])
		utils.MakeErrorResponse(w, err)

		return
	}

	output, _ := json.Marshal(response)
	_, _ = w.Write(output)
}

func (t *UserDelivery) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	nickname := vars[consts.UserNickNamePath]
	requestBody := new(models.RequestBody)
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		utils.MakeErrorResponse(w, err)
		return
	}

	response, err := t.UserUseCase.UpdateUserProfile(nickname, requestBody)
	if err != nil {
		w.WriteHeader(customErrors.StatusCodes[err])
		utils.MakeErrorResponse(w, err)
		return
	}

	output, _ := json.Marshal(response)
	_, _ = w.Write(output)
}
