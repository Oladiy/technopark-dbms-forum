package database_service

import (
	"database/sql"
	"github.com/gorilla/mux"
	databaseService "technopark-dbms-forum/internal/pkg/database_service"
	"technopark-dbms-forum/internal/pkg/database_service/delivery"
	"technopark-dbms-forum/internal/pkg/database_service/repository"
)

type Service struct {
	Delivery *delivery.DatabaseServiceDelivery
	Repository databaseService.Repository
	Router *mux.Router
}

func Run(connectionDB *sql.DB) *Service{
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