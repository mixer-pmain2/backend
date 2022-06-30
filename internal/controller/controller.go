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
	User           *user
	Patient        *patient
	Doctor         *doctor
	Spr            *spr
	Administration *administration
}

func Init() *Controller {
	return &Controller{
		User:           initUserController(),
		Patient:        initPatientController(),
		Doctor:         initDoctorController(),
		Spr:            initSprController(),
		Administration: initAdministrationController(),
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
