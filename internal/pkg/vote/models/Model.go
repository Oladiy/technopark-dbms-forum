package models

type Vote struct {
	Id			int		`json:"id"`
	Nickname 	string 	`json:"nickname"`
	Voice 		int 	`json:"voice"`
}

type RequestBody struct {
	Nickname 	string 	`json:"nickname"`
	Voice 		int 	`json:"voice"`
}
