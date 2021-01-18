package thread

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"technopark-dbms-forum/internal/pkg/common/consts"
	"technopark-dbms-forum/internal/pkg/thread"
	"technopark-dbms-forum/internal/pkg/thread/delivery"
	"technopark-dbms-forum/internal/pkg/thread/repository"
	"technopark-dbms-forum/internal/pkg/thread/usecase"
)

type Service struct {
	Delivery *delivery.ThreadDelivery
	Router *mux.Router
	UseCase thread.UseCase
	Repository thread.Repository
}

func Run(connectionDB *sql.DB) *Service {
	threadRepository := repository.NewThreadRepository(connectionDB)
	threadUseCase := usecase.NewThreadUseCase(threadRepository)
	threadDelivery := delivery.NewThreadDelivery(threadUseCase)

	router := mux.NewRouter()
	router.HandleFunc(fmt.Sprintf("/api/thread/{%s:.+}/create", consts.ThreadSlugPath), threadDelivery.CreateThreadPosts)
	router.HandleFunc(fmt.Sprintf("/api/thread/{%s:.+}/details", consts.ThreadSlugPath), threadDelivery.GetThread)
	router.HandleFunc(fmt.Sprintf("/api/thread/{%s:.+}/vote", consts.ThreadSlugPath), threadDelivery.ThreadVote)

	return &Service {
		Delivery: threadDelivery,
		Repository: threadRepository,
		UseCase: threadUseCase,
		Router: router,
	}
}
