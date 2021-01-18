package delivery

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"technopark-dbms-forum/internal/pkg/common/consts"
	customErrors "technopark-dbms-forum/internal/pkg/common/custom_errors"
	"technopark-dbms-forum/internal/pkg/common/utils"
	"technopark-dbms-forum/internal/pkg/forum"
	"technopark-dbms-forum/internal/pkg/thread"
)

type ForumDelivery struct {
	ForumUseCase forum.UseCase
	ThreadUseCase thread.UseCase
}

func NewForumDelivery(forumUseCase forum.UseCase) *ForumDelivery {
	return &ForumDelivery {
		ForumUseCase: forumUseCase,
	}
}

func (t *ForumDelivery) CreateForum(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	requestBody := new(forum.RequestBody)
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		utils.MakeErrorResponse(w, err)
		return
	}

	response, err := t.ForumUseCase.CreateForum(requestBody)
	if err != nil {
		w.WriteHeader(customErrors.StatusCodes[err])
		if response != nil {
			output, _ := json.Marshal(response)
			_, _ = w.Write(output)
		} else {
			utils.MakeErrorResponse(w, err)
		}

		return
	}

	w.WriteHeader(http.StatusCreated)
	output, _ := json.Marshal(response)
	_, _ = w.Write(output)
}

func (t *ForumDelivery) GetForumDetails(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	slug := vars[consts.ForumSlugPath]

	response, err := t.ForumUseCase.GetForumDetails(slug)
	if err != nil {
		w.WriteHeader(customErrors.StatusCodes[err])
		utils.MakeErrorResponse(w, err)

		return
	}

	output, _ := json.Marshal(response)
	_, _ = w.Write(output)
}

func (t *ForumDelivery) CreateForumThread(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	slug := vars[consts.ForumSlugPath]
	requestBody := new(thread.RequestBody)
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		utils.MakeErrorResponse(w, err)
		return
	}

	response, err := t.ForumUseCase.CreateForumThread(slug, requestBody)
	if err != nil {
		w.WriteHeader(customErrors.StatusCodes[err])
		if response != nil {
			output, _ := json.Marshal(response)
			_, _ = w.Write(output)
		} else {
			utils.MakeErrorResponse(w, err)
		}

		return
	}

	w.WriteHeader(http.StatusCreated)
	output, _ := json.Marshal(response)
	_, _ = w.Write(output)
}

func (t *ForumDelivery) GetForumThreads(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var limit int
	var since string
	var desc bool
	var err error

	vars := mux.Vars(r)
	slug := vars[consts.ForumSlugPath]

	parameterLimit, ok := r.URL.Query()["limit"]
	if ok && len(parameterLimit[0]) > 0 {
		limit, err = strconv.Atoi(parameterLimit[0])
		if err != nil {
			limit = 0
		}
	}
	parameterSince, ok := r.URL.Query()["since"]
	if ok && len(parameterSince[0]) > 0 {
		since = parameterSince[0]
	}
	parameterDesc, ok := r.URL.Query()["desc"]
	if ok && len(parameterDesc) > 0 {
		if parameterDesc[0] == "true" {
			desc = true
		}
	}

	response, err := t.ForumUseCase.GetForumThreads(slug, limit, since, desc)
	if err != nil {
		w.WriteHeader(customErrors.StatusCodes[err])
		utils.MakeErrorResponse(w, err)

		return
	}

	output, _ := json.Marshal(response)
	_, _ = w.Write(output)
}
