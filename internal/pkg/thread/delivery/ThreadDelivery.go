package delivery

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"technopark-dbms-forum/internal/pkg/common/consts"
	customErrors "technopark-dbms-forum/internal/pkg/common/custom_errors"
	"technopark-dbms-forum/internal/pkg/common/utils"
	"technopark-dbms-forum/internal/pkg/post"
	"technopark-dbms-forum/internal/pkg/thread"
	"technopark-dbms-forum/internal/pkg/vote"
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
		slug = slugOrId
	}

	requestBody := make([]post.RequestBody, 0)
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		utils.MakeErrorResponse(w, err)
		return
	}

	response, err := t.ThreadUseCase.CreateThreadPosts(slug, id, &requestBody)
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

	v := new(vote.Vote)
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		utils.MakeErrorResponse(w, err)
		return
	}

	response, err := t.ThreadUseCase.ThreadVote(id, slug, v)
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

	w.WriteHeader(http.StatusOK)
	output, _ := json.Marshal(*response)
	_, _ = w.Write(output)
}
