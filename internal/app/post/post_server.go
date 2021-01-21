package post

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
	"net/http"
	"technopark-dbms-forum/internal/pkg/common/consts"
	forumRep "technopark-dbms-forum/internal/pkg/forum/repository"
	"technopark-dbms-forum/internal/pkg/post"
	"technopark-dbms-forum/internal/pkg/post/delivery"
	"technopark-dbms-forum/internal/pkg/post/repository"
	"technopark-dbms-forum/internal/pkg/post/usecase"
	threadRep "technopark-dbms-forum/internal/pkg/thread/repository"
	userRep "technopark-dbms-forum/internal/pkg/user/repository"
)

type Service struct {
	Delivery *delivery.PostDelivery
	Router *mux.Router
	UseCase post.UseCase
	Repository post.Repository
}

func Run(connectionDB *pgx.ConnPool) *Service {
	forumRepository := forumRep.NewForumRepository(connectionDB)
	postRepository := repository.NewPostRepository(connectionDB)
	threadRepository := threadRep.NewThreadRepository(connectionDB)
	userRepository := userRep.NewUserRepository(connectionDB)

	postUseCase := usecase.NewPostUseCase(forumRepository, postRepository, threadRepository, userRepository)
	postDelivery := delivery.NewPostDelivery(postUseCase)

	router := mux.NewRouter()
	router.HandleFunc(fmt.Sprintf("/api/post/{%s:[0-9]+}/details", consts.PostIdPath), postDelivery.GetPost).Methods(http.MethodGet)
	router.HandleFunc(fmt.Sprintf("/api/post/{%s:[0-9]+}/details", consts.PostIdPath), postDelivery.UpdatePost).Methods(http.MethodPost)

	return &Service {
		Delivery: postDelivery,
		Repository: postRepository,
		UseCase: postUseCase,
		Router: router,
	}
}
