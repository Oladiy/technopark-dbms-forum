package thread

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
	"net/http"
	"technopark-dbms-forum/internal/pkg/common/consts"
	"technopark-dbms-forum/internal/pkg/thread"
	"technopark-dbms-forum/internal/pkg/thread/delivery"
	"technopark-dbms-forum/internal/pkg/thread/repository"
	"technopark-dbms-forum/internal/pkg/thread/usecase"
	userRep "technopark-dbms-forum/internal/pkg/user/repository"
)

type Service struct {
	Delivery *delivery.ThreadDelivery
	Router *mux.Router
	UseCase thread.UseCase
	Repository thread.Repository
}

func Run(connectionDB *pgx.ConnPool) *Service {
	userRepository := userRep.NewUserRepository(connectionDB)
	threadRepository := repository.NewThreadRepository(connectionDB)
	threadUseCase := usecase.NewThreadUseCase(threadRepository, userRepository)
	threadDelivery := delivery.NewThreadDelivery(threadUseCase)

	router := mux.NewRouter()
	router.HandleFunc(fmt.Sprintf("/api/thread/{%s:.+}/create", consts.ThreadSlugPath), threadDelivery.CreateThreadPosts)
	router.HandleFunc(fmt.Sprintf("/api/thread/{%s:.+}/details", consts.ThreadSlugPath), threadDelivery.GetThread).Methods(http.MethodGet)
	router.HandleFunc(fmt.Sprintf("/api/thread/{%s:.+}/details", consts.ThreadSlugPath), threadDelivery.UpdateThread).Methods(http.MethodPost)
	router.HandleFunc(fmt.Sprintf("/api/thread/{%s:.+}/vote", consts.ThreadSlugPath), threadDelivery.ThreadVote)
	router.HandleFunc(fmt.Sprintf("/api/thread/{%s:.+}/posts", consts.ThreadSlugPath), threadDelivery.GetThreadPosts)

	return &Service {
		Delivery: threadDelivery,
		Repository: threadRepository,
		UseCase: threadUseCase,
		Router: router,
	}
}
