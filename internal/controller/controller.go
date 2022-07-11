package controller

import (
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
