package user

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
	"net/http"
	"technopark-dbms-forum/internal/pkg/common/consts"
	"technopark-dbms-forum/internal/pkg/user"
	"technopark-dbms-forum/internal/pkg/user/delivery"
	"technopark-dbms-forum/internal/pkg/user/repository"
	"technopark-dbms-forum/internal/pkg/user/usecase"
)

type Service struct {
	Delivery *delivery.UserDelivery
	Router *mux.Router
	UseCase user.UseCase
	Repository user.Repository
}

func Run(connectionDB *pgx.ConnPool) *Service {
	userRepository := repository.NewUserRepository(connectionDB)
	userUseCase := usecase.NewUserUseCase(userRepository)
	userDelivery := delivery.NewUserDelivery(userUseCase)

	router := mux.NewRouter()
	router.HandleFunc(fmt.Sprintf("/api/user/{%s:.+}/create", consts.UserNickNamePath), userDelivery.CreateUser)
	router.HandleFunc(fmt.Sprintf("/api/user/{%s:.+}/profile", consts.UserNickNamePath), userDelivery.GetUserProfile).Methods(http.MethodGet)
	router.HandleFunc(fmt.Sprintf("/api/user/{%s:.+}/profile", consts.UserNickNamePath), userDelivery.UpdateUserProfile).Methods(http.MethodPost)

	return &Service {
		Delivery: userDelivery,
		Repository: userRepository,
		UseCase: userUseCase,
		Router: router,
	}
}
