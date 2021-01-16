package forum

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"technopark-dbms-forum/internal/pkg/common/consts"
	"technopark-dbms-forum/internal/pkg/forum"
	"technopark-dbms-forum/internal/pkg/forum/delivery"
	"technopark-dbms-forum/internal/pkg/forum/repository"
	"technopark-dbms-forum/internal/pkg/forum/usecase"
	threadRep "technopark-dbms-forum/internal/pkg/thread/repository"
	threadUC "technopark-dbms-forum/internal/pkg/thread/usecase"
	userRep "technopark-dbms-forum/internal/pkg/user/repository"
)

type Service struct {
	Delivery *delivery.ForumDelivery
	Router *mux.Router
	UseCase forum.UseCase
	Repository forum.Repository
}

func Run(connectionDB *sql.DB) *Service {
	threadRepository := threadRep.NewThreadRepository(connectionDB)
	userRepository := userRep.NewUserRepository(connectionDB)
	forumRepository := repository.NewForumRepository(connectionDB)
	forumUseCase := usecase.NewForumUseCase(forumRepository, userRepository)
	threadUseCase := threadUC.NewForumUseCase(threadRepository, forumRepository, userRepository)
	forumDelivery := delivery.NewForumDelivery(forumUseCase, threadUseCase)

	router := mux.NewRouter()
	router.HandleFunc(fmt.Sprintf("/api/forum/create"), forumDelivery.CreateForum)
	router.HandleFunc(fmt.Sprintf("/api/forum/{%s:.+}/details", consts.ForumSlugPath), forumDelivery.GetForumDetails).Methods(http.MethodGet)
	router.HandleFunc(fmt.Sprintf("/api/forum/{%s:.+}/create", consts.ForumSlugPath), forumDelivery.CreateForumThread)

	return &Service {
		Delivery: forumDelivery,
		Repository: forumRepository,
		UseCase: forumUseCase,
		Router: router,
	}
}
