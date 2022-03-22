package controller

import (
	"pmain2/internal/database"
	"pmain2/internal/models"
)

type userInterface interface {
	isAuth()
}

type user struct{}

func initUserController() *user {
	return &user{}
}

func (u *user) IsAuth(login, password string) (bool, error) {
	conn, err := database.Connect()
	if err != nil {
		ERROR.Println(err.Error())
		return false, err
	}
	model := models.SprDoctModel{Db: conn.DB}
	ok, err := model.UserAuth(login, password)
	if err != nil {
		ERROR.Println(err.Error())
		return false, err
	}
	return ok, nil
}
