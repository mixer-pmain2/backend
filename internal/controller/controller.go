package controller

import (
	"database/sql"
	"pmain2/internal/database"
	"pmain2/pkg/logger"
)

var (
	INFO, _  = logger.New("controller", logger.INFO)
	ERROR, _ = logger.New("controller", logger.ERROR)
)

type Controller struct {
	User    *user
	Patient *patient
	Spr     *spr
}

func Init() *Controller {
	return &Controller{
		User:    initUserController(),
		Patient: initPatientController(),
		Spr:     initSprController(),
	}
}

func CreateTx() (error, *sql.Tx) {

	conn, err := database.Connect()
	if err != nil {
		return err, nil
	}
	tx, err := conn.DB.Begin()
	if err != nil {
		return err, nil
	}
	return nil, tx
}
