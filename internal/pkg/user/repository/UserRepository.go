package repository

import (
	"fmt"
	"github.com/jackc/pgx"
	customErrors "technopark-dbms-forum/internal/pkg/common/custom_errors"
	"technopark-dbms-forum/internal/pkg/user/models"
)

type UserRepository struct {
	connectionDB *pgx.ConnPool
}

func NewUserRepository(connectionDB *pgx.ConnPool) *UserRepository {
	return &UserRepository {
		connectionDB: connectionDB,
	}
}

func (t* UserRepository) CreateUser(nickname string, profile *models.RequestBody) (*[]models.User, error) {
	queryInsert := `INSERT INTO Users (nickname, fullname, about, email) 
					VALUES($1, $2, $3, $4);`
	querySelect := `SELECT nickname, fullname, about, email 
					FROM Users 
					WHERE nickname = $1 or email = $2;`
	selection := make([]models.User, 0)

	_, err := t.connectionDB.Exec(queryInsert, nickname, profile.FullName, profile.About, profile.Email)
	if err != nil {
		querySelectResult, err := t.connectionDB.Query(querySelect, nickname, profile.Email)
		if err != nil || querySelectResult.Err() != nil {
			return nil, customErrors.IncorrectInputData
		}

		for querySelectResult.Next() {
			u := new(models.User)
			_ = querySelectResult.Scan(&u.Nickname, &u.FullName, &u.About, &u.Email)
			selection = append(selection, *u)
		}

		querySelectResult.Close()
		return &selection, customErrors.UserAlreadyExists
	}

	selection = append(selection, models.User{
		Nickname: nickname,
		FullName: profile.FullName,
		About: profile.About,
		Email: profile.Email,
	})
	return &selection, nil
}

func (t* UserRepository) GetUserProfile(nickname string) (*models.User, error) {
	querySelect := `SELECT nickname, fullname, about, email 
					FROM Users 
					WHERE nickname = $1;`
	selection := t.connectionDB.QueryRow(querySelect, nickname)
	if selection == nil {
		return nil, customErrors.UserNotFound
	}
	u := new(models.User)
	if err := selection.Scan(&u.Nickname, &u.FullName, &u.About, &u.Email); err != nil {
		return nil, customErrors.UserNotFound
	}

	return u, nil
}

func (t* UserRepository) UpdateUserProfile(nickname string, profile *models.RequestBody) (*models.User, error) {
	queryUpdate := "UPDATE Users SET "
	fieldsToUpdate := make([]interface{}, 0)
	counter := 1
	isFirst := true

	if len(profile.FullName) != 0 {
		queryUpdate += fmt.Sprintf("fullname=$%d ", counter)
		isFirst = false
		fieldsToUpdate = append(fieldsToUpdate, profile.FullName)
		counter++
	}

	if len(profile.About) != 0 {
		if isFirst {
			queryUpdate += fmt.Sprintf("about=$%d ", counter)
			isFirst = false
		} else {
			queryUpdate += fmt.Sprintf(", about=$%d ", counter)
		}
		fieldsToUpdate = append(fieldsToUpdate, profile.About)
		counter++
	}

	if len(profile.Email) != 0 {
		if isFirst {
			queryUpdate += fmt.Sprintf("email=$%d ", counter)
		} else {
			queryUpdate += fmt.Sprintf(", email=$%d ", counter)
		}
		fieldsToUpdate = append(fieldsToUpdate, profile.Email)
		counter++
	}

	queryUpdate += fmt.Sprintf("WHERE nickname=$%d returning nickname, fullname, about, email", counter)
	fieldsToUpdate = append(fieldsToUpdate, nickname)

	selection := new(models.User)
	err := t.connectionDB.QueryRow(queryUpdate, fieldsToUpdate...).Scan(&selection.Nickname, &selection.FullName,
																		&selection.About, &selection.Email)
	if err != nil {
		return nil, customErrors.UserProfileDataConflict
	}

	return selection, nil
}
