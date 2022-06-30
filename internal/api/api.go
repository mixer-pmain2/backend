package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"pmain2/internal/consts"
	"pmain2/internal/types"
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
	User           *userApi
	Patient        *patientApi
	Doctor         *doctorApi
	Spr            *sprApi
	Administration *administrationApi
}

func Init() *Api {
	return &Api{
		User:           userApiInit(),
		Patient:        patientApiInit(),
		Doctor:         doctorApiInit(),
		Spr:            sprApiInit(),
		Administration: administrationApiInit(),
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

func resSuccess(val int) []byte {
	res := types.HttpResponse{Success: true, Error: 0}

	if val > 0 {
		res.Success = false
		res.Error = val
		res.Message = consts.ArrErrors[val]
	}
	resMarshal, _ := json.Marshal(res)
	return resMarshal
}
