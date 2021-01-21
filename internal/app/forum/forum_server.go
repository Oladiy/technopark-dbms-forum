package forum

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
	"net/http"
	"technopark-dbms-forum/internal/pkg/common/consts"
	"technopark-dbms-forum/internal/pkg/forum"
	"technopark-dbms-forum/internal/pkg/forum/delivery"
	"technopark-dbms-forum/internal/pkg/forum/repository"
	"technopark-dbms-forum/internal/pkg/forum/usecase"
	threadRep "technopark-dbms-forum/internal/pkg/thread/repository"
	userRep "technopark-dbms-forum/internal/pkg/user/repository"
)

type Service struct {
	Delivery *delivery.ForumDelivery
	Router *mux.Router
	UseCase forum.UseCase
	Repository forum.Repository
}

func Run(connectionDB *pgx.ConnPool) *Service {
	threadRepository := threadRep.NewThreadRepository(connectionDB)
	userRepository := userRep.NewUserRepository(connectionDB)
	forumRepository := repository.NewForumRepository(connectionDB)
	forumUseCase := usecase.NewForumUseCase(forumRepository, threadRepository, userRepository)
	forumDelivery := delivery.NewForumDelivery(forumUseCase)

	router := mux.NewRouter()
	router.HandleFunc(fmt.Sprintf("/api/forum/create"), forumDelivery.CreateForum)
	router.HandleFunc(fmt.Sprintf("/api/forum/{%s:.+}/details", consts.ForumSlugPath), forumDelivery.GetForumDetails).Methods(http.MethodGet)
	router.HandleFunc(fmt.Sprintf("/api/forum/{%s:.+}/create", consts.ForumSlugPath), forumDelivery.CreateForumThread)
	router.HandleFunc(fmt.Sprintf("/api/forum/{%s:.+}/threads", consts.ForumSlugPath), forumDelivery.GetForumThreads)
	router.HandleFunc(fmt.Sprintf("/api/forum/{%s:.+}/users", consts.ForumSlugPath), forumDelivery.GetForumUsers)

	return &Service {
		Delivery: forumDelivery,
		Repository: forumRepository,
		UseCase: forumUseCase,
		Router: router,
	}
}
