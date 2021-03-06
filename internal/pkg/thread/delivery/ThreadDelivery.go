package delivery

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"technopark-dbms-forum/internal/pkg/common/consts"
	customErrors "technopark-dbms-forum/internal/pkg/common/custom_errors"
	"technopark-dbms-forum/internal/pkg/common/utils"
	postModel "technopark-dbms-forum/internal/pkg/post/models"
	"technopark-dbms-forum/internal/pkg/thread"
	threadModel "technopark-dbms-forum/internal/pkg/thread/models"
	voteModel "technopark-dbms-forum/internal/pkg/vote/models"
)

type ThreadDelivery struct {
	ThreadUseCase thread.UseCase
}

func NewThreadDelivery(threadUseCase thread.UseCase) *ThreadDelivery {
	return &ThreadDelivery {
		ThreadUseCase: threadUseCase,
	}
}

func (t *ThreadDelivery) CreateThreadPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)

	var slug string
	var id int
	slugOrId := vars[consts.ThreadSlugPath]

	id, err := strconv.Atoi(slugOrId)
	if err != nil {
		id = 0
		slug = slugOrId
	}

	requestBody := make([]postModel.RequestBody, 0)
	if err = json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		utils.MakeErrorResponse(w, err)
		return
	}

	response, err := t.ThreadUseCase.CreateThreadPosts(slug, id, &requestBody)
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
	output, _ := json.Marshal(*response)
	_, _ = w.Write(output)
}

func (t *ThreadDelivery) GetThread(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)

	var slug string
	var id int
	slugOrId := vars[consts.ThreadSlugPath]

	id, err := strconv.Atoi(slugOrId)
	if err != nil {
		slug = slugOrId
	}

	response, err := t.ThreadUseCase.GetThread("", id, slug)
	if err != nil {
		w.WriteHeader(customErrors.StatusCodes[err])
		utils.MakeErrorResponse(w, err)

		return
	}

	output, _ := json.Marshal(response)
	_, _ = w.Write(output)
}

func (t *ThreadDelivery) ThreadVote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)

	var slug string
	var id int
	slugOrId := vars[consts.ThreadSlugPath]

	id, err := strconv.Atoi(slugOrId)
	if err != nil {
		slug = slugOrId
	}

	v := new(	voteModel.Vote)
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		utils.MakeErrorResponse(w, err)
		return
	}

	response, err := t.ThreadUseCase.ThreadVote(id, slug, v)
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

	w.WriteHeader(http.StatusOK)
	output, _ := json.Marshal(*response)
	_, _ = w.Write(output)
}

func (t *ThreadDelivery) GetThreadPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var desc bool
	var limit int
	var since int
	var sort string
	var err error

	parameterLimit := r.URL.Query().Get("limit")
	limit, err = strconv.Atoi(parameterLimit)
	if err != nil {
		limit = 100
	}

	parameterSince := r.URL.Query().Get("since")
	since, err = strconv.Atoi(parameterSince)
	if err != nil {
		since = -1
	}

	parameterSort := r.URL.Query().Get("sort")
	if len(parameterSort) > 0 {
		if parameterSort == "tree" {
			sort = "tree"
		}
		if parameterSort == "parent_tree" {
			sort = "parent_tree"
		}
	}

	parameterDesc := r.URL.Query().Get("desc")
	if parameterDesc == "true" {
		desc = true
	}

	vars := mux.Vars(r)

	var slug string
	var id int
	slugOrId := vars[consts.ThreadSlugPath]

	id, err = strconv.Atoi(slugOrId)
	if err != nil {
		slug = slugOrId
	}

	response, err := t.ThreadUseCase.GetThreadPosts(id, slug, limit, since, sort, desc)
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

	output, _ := json.Marshal(response)
	_, _ = w.Write(output)
}

func (t *ThreadDelivery) UpdateThread(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)

	var slug string
	var id int
	slugOrId := vars[consts.ThreadSlugPath]

	id, err := strconv.Atoi(slugOrId)
	if err != nil {
		slug = slugOrId
	}

	threadToUpdate := new(threadModel.RequestBody)
	if err := json.NewDecoder(r.Body).Decode(&threadToUpdate); err != nil {
		utils.MakeErrorResponse(w, err)
		return
	}

	response, err := t.ThreadUseCase.UpdateThread(id, slug, threadToUpdate)
	if err != nil {
		w.WriteHeader(customErrors.StatusCodes[err])
		utils.MakeErrorResponse(w, err)

		return
	}

	output, _ := json.Marshal(response)
	_, _ = w.Write(output)
}
