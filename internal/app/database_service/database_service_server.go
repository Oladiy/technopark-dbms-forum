package database_service

import (
	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
	databaseService "technopark-dbms-forum/internal/pkg/database_service"
	"technopark-dbms-forum/internal/pkg/database_service/delivery"
	"technopark-dbms-forum/internal/pkg/database_service/repository"
)

type Service struct {
	Delivery *delivery.DatabaseServiceDelivery
	Repository databaseService.Repository
	Router *mux.Router
}

func Run(connectionDB *pgx.ConnPool) *Service{
	databaseServiceRepository := repository.NewDatabaseServiceRepository(connectionDB)
	databaseServiceDelivery := delivery.NewDatabaseServiceDelivery(databaseServiceRepository)

	router := mux.NewRouter()

	router.HandleFunc("/api/service/clear", databaseServiceDelivery.ClearDatabase)
	router.HandleFunc("/api/service/status", databaseServiceDelivery.GetStatus)

	return &Service{
		Delivery: databaseServiceDelivery,
		Repository: databaseServiceRepository,
		Router: router,
	}
}