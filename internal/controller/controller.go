package controller

import (
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
