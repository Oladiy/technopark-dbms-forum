package post

type Post struct {
	Id int `json:"id"`
	Parent int `json:"parent"`
	Author string `json:"author"`
	Message string `json:"message"`
	IsEdited bool `json:"isEdited"`
	Forum string `json:"forum"`
	Thread int `json:"thread"`
	Created string `json:"created"`
}

type RequestBody struct {
	Parent int `json:"parent"`
	Author string `json:"author"`
	Message string `json:"message"`
}