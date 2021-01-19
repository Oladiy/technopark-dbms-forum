package models

import (
	models3 "technopark-dbms-forum/internal/pkg/forum/models"
	"technopark-dbms-forum/internal/pkg/thread/models"
	models2 "technopark-dbms-forum/internal/pkg/user/models"
)

type Post struct {
	Id 			int 	`json:"id"`
	Parent 		int 	`json:"parent"`
	Author 		string 	`json:"author"`
	Message 	string 	`json:"message"`
	IsEdited 	bool 	`json:"isEdited"`
	Forum 		string 	`json:"forum"`
	Thread 		int 	`json:"thread"`
	Created 	string 	`json:"created"`
}

type RequestBody struct {
	Parent 	int 	`json:"parent"`
	Author 	string 	`json:"author"`
	Message string 	`json:"message"`
}

type Details struct {
	Forum 	*models3.Forum  `json:"forum"`
	Post 	*Post            `json:"post"`
	Thread 	*models.Thread `json:"thread"`
	User 	*models2.User    `json:"author"`
}
