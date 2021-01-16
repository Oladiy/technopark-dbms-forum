package usecase

import (
	customErrors "technopark-dbms-forum/internal/pkg/common/custom_errors"
	"technopark-dbms-forum/internal/pkg/user"
)

type UserUseCase struct {
	UserRepository user.Repository
}

func NewUserUseCase(userRepository user.Repository) *UserUseCase {
	return &UserUseCase {
		UserRepository: userRepository,
	}
}

func (t* UserUseCase) CreateUser(nickname string, profile *user.RequestBody) (*[]user.User, error) {
	if len(nickname) == 0 || len(profile.Email) == 0 {
		return nil, customErrors.IncorrectInputData
	}

	return t.UserRepository.CreateUser(nickname, profile)
}

func (t* UserUseCase) GetUserProfile(nickname string) (*user.User, error) {
	if len(nickname) == 0 {
		return nil, customErrors.IncorrectInputData
	}

	return t.UserRepository.GetUserProfile(nickname)
}

func (t* UserUseCase) UpdateUserProfile(nickname string, profile *user.RequestBody) (*user.User, error) {
	userRelevant, err := t.UserRepository.GetUserProfile(nickname)
	if err != nil {
		switch err {
			case customErrors.IncorrectInputData: return nil, customErrors.IncorrectInputData
			case customErrors.UserNotFound: return nil, customErrors.UserNotFound
		}
	}

	if len(profile.FullName) == 0 && len(profile.About) == 0 && len(profile.Email) == 0 {
		return userRelevant, nil
	}

	return t.UserRepository.UpdateUserProfile(nickname, profile)
}
