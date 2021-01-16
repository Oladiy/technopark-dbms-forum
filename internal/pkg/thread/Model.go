package thread

type Thread struct {
	Id 		int 	`json:"id"`
	Title 	string 	`json:"title"`
	Author 	string 	`json:"author"`
	Forum 	string 	`json:"forum"`
	Message string 	`json:"message"`
	Votes 	int 	`json:"votes"`
	Slug 	string 	`json:"slug"`
	Created string 	`json:"created"`
}

type RequestBody struct {
	Author 	string 	`json:"author"`
	Created string 	`json:"created"`
	Forum 	string 	`json:"forum"`
	Message string 	`json:"message"`
	Slug 	string 	`json:"slug"`
	Title 	string 	`json:"title"`
}
