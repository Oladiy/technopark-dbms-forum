package main

import (
	"github.com/jackc/pgx"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"strconv"
	databaseService "technopark-dbms-forum/internal/app/database_service"
	"technopark-dbms-forum/internal/app/forum"
	"technopark-dbms-forum/internal/app/post"
	"technopark-dbms-forum/internal/app/thread"
	"technopark-dbms-forum/internal/app/user"
	"time"
)

type ServerConfig struct {
	DatabaseService *databaseService.Service
	ForumService	*forum.Service
	PostService		*post.Service
	ThreadService 	*thread.Service
	UserService   	*user.Service
}

type ServiceConfig struct {
	Domain           string
	Port             int
	DatabaseUser     string
	DatabasePassword string
	DatabaseName     string
	DatabaseDomain 	 string
	DatabasePort   	 int
}

func CreateDBConnection(config *ServiceConfig) (*pgx.ConnPool, error) {
	pool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			User:     config.DatabaseUser,
			Password: config.DatabasePassword,
			Port:     uint16(config.DatabasePort),
			Database: config.DatabaseName,
			Host: config.DatabaseDomain,
		},
		MaxConnections: 100,
	})
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func configureMainRouter(application *ServerConfig) http.Handler{
	handler := http.NewServeMux()

	handler.Handle("/api/forum/", application.ForumService.Router)
	handler.Handle("/api/post/", application.PostService.Router)
	handler.Handle("/api/service/", application.DatabaseService.Router)
	handler.Handle("/api/thread/", application.ThreadService.Router)
	handler.Handle("/api/user/", application.UserService.Router)

	return handler
}

func InitService(connectionDB *pgx.ConnPool) *ServerConfig{
	forumService 	:= forum.Run(connectionDB)
	postService 	:= post.Run(connectionDB)
	dbService		:= databaseService.Run(connectionDB)
	threadService 	:= thread.Run(connectionDB)
	userService 	:= user.Run(connectionDB)

	return &ServerConfig{
		DatabaseService: 	dbService,
		ForumService:		forumService,
		PostService: 		postService,
		ThreadService:		threadService,
		UserService:   		userService,
	}
}

func main() {
	var config = new(ServiceConfig)

	config.Domain = ""
	config.Port = 5000
	config.DatabaseDomain = "localhost"
	config.DatabasePort = 5432
	config.DatabaseName = "Forum"
	config.DatabaseUser = "docker"
	config.DatabasePassword = "docker"

	conn, err := CreateDBConnection(config)
	if err != nil{
		log.Fatalln(err)
	}

	conf := InitService(conn)
	URLHandler := configureMainRouter(conf)
	httpServer := &http.Server{
		Addr: config.Domain + ":" + strconv.Itoa(config.Port),
		Handler: URLHandler,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("STARTING SERVER AT PORT: ", config.Port)
	serverErr := httpServer.ListenAndServe()
	if serverErr != nil{
		log.Fatalln(serverErr)
	}

	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()
}
