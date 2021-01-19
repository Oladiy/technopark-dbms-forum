package delivery

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
	"technopark-dbms-forum/internal/pkg/common/consts"
	customErrors "technopark-dbms-forum/internal/pkg/common/custom_errors"
	"technopark-dbms-forum/internal/pkg/common/utils"
	"technopark-dbms-forum/internal/pkg/post"
	"technopark-dbms-forum/internal/pkg/thread/models"
)

type PostDelivery struct {
	PostUseCase post.UseCase
}

func NewPostDelivery(postUseCase post.UseCase) *PostDelivery {
	return &PostDelivery {
		PostUseCase: postUseCase,
	}
}

func (t *PostDelivery) GetPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	parameterId := vars[consts.PostIdPath]
	id, err := strconv.Atoi(parameterId)

	var isUserRelated bool
	var isThreadRelated bool
	var isForumRelated bool

	related := r.URL.Query().Get(consts.RelatedPath)
	for _, value := range strings.Split(related, ",") {
		if value == "user" {
			isUserRelated = true
		} else if value == "thread" {
			isThreadRelated = true
		} else if value == "forum" {
			isForumRelated = true
		}
	}

	response, err := t.PostUseCase.GetPost(id, isUserRelated, isThreadRelated, isForumRelated)
	if err != nil {
		w.WriteHeader(customErrors.StatusCodes[err])
		utils.MakeErrorResponse(w, err)

		return
	}

	output, _ := json.Marshal(response)
	_, _ = w.Write(output)
}


func (t *PostDelivery) UpdatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	parameterId := vars[consts.PostIdPath]
	id, err := strconv.Atoi(parameterId)

	p := new(models.RequestBody)
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		utils.MakeErrorResponse(w, err)
		return
	}
	message := p.Message

	response, err := t.PostUseCase.UpdatePost(id, message)
	if err != nil {
		w.WriteHeader(customErrors.StatusCodes[err])
		utils.MakeErrorResponse(w, err)

		return
	}

	output, _ := json.Marshal(response)
	_, _ = w.Write(output)
}
