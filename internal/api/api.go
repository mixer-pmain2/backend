package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"pmain2/pkg/logger"
	"strconv"
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

type apiParams struct {
	id      int
	isCache bool
}

func getParams(r *http.Request, t interface{}) *apiParams {

	params := mux.Vars(r)
	var err error
	p := apiParams{}
	p.id, err = strconv.Atoi(params["id"])
	if err != nil {
		p.id = 0
	}

	p.isCache, err = strconv.ParseBool(r.URL.Query().Get("cache"))
	if err != nil {
		p.isCache = true
	}

	if t != nil {
		fmt.Println(r.Body)
		json.NewDecoder(r.Body).Decode(&t)
	}

	return &p
}
