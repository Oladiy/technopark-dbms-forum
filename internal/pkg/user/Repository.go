package user

type Repository interface {
	CreateUser(nickname string, user *RequestBody) (*[]User, error)
	GetUserProfile(nickname string) (*User, error)
	UpdateUserProfile(nickname string, user *RequestBody) (*User, error)
}
