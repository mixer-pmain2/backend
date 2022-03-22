package api

import (
	"encoding/json"
	"pmain2/pkg/logger"
)

var (
	INFO, _  = logger.New("api", logger.INFO)
	ERROR, _ = logger.New("api", logger.ERROR)

	AnswerOk   = Success{Success: true}
	AnswerFail = Success{Success: false}
)

type Api struct {
	User    *userApi
	Patient *patientApi
	Spr     *sprApi
}

func Init() *Api {
	return &Api{
		User:    userApiInit(),
		Patient: patientApiInit(),
		Spr:     sprApiInit(),
	}
}

type Success struct {
	Success bool `json:"success"`
}

func (s Success) Marshal() string {
	marshal, err := json.Marshal(s)
	if err != nil {
		return ""
	}
	return string(marshal)
}
