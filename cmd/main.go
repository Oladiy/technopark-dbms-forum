package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"strconv"
	"technopark-dbms-forum/internal/app/forum"
	"technopark-dbms-forum/internal/app/post"
	"technopark-dbms-forum/internal/app/thread"
	"technopark-dbms-forum/internal/app/user"
	"time"
)

type ServerConfig struct {
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


func CreateDBConnection(config *ServiceConfig) (*sql.DB, error) {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.DatabaseDomain, config.DatabasePort, config.DatabaseUser,
		config.DatabasePassword, config.DatabaseName)

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	log.Println("Connected to database")
	return db, nil
}

func configureMainRouter(application *ServerConfig) http.Handler{
	handler := http.NewServeMux()

	handler.Handle("/api/forum/", application.ForumService.Router)
	handler.Handle("/api/post/", application.PostService.Router)
	handler.Handle("/api/thread/", application.ThreadService.Router)
	handler.Handle("/api/user/", application.UserService.Router)

	return handler
}

func InitService(connectionDB *sql.DB) *ServerConfig{
	forumService 	:= forum.Run(connectionDB)
	postService 	:= post.Run(connectionDB)
	threadService 	:= thread.Run(connectionDB)
	userService 	:= user.Run(connectionDB)

	return &ServerConfig{
		ForumService:	forumService,
		PostService: 	postService,
		ThreadService:	threadService,
		UserService:   	userService,
	}
}

func main() {
	var config = new(ServiceConfig)

	config.Domain = ""
	config.Port = 5000
	config.DatabaseDomain = "postgresql"
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
