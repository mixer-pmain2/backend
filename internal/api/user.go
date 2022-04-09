package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pmain2/internal/controller"
	"strconv"

	"github.com/gorilla/mux"

	"pmain2/internal/database"
	"pmain2/internal/models"
)

type userApi struct{}

func userApiInit() *userApi {
	return &userApi{}
}

func (u *userApi) GetUser(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		return nil
	}
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		return err
	}

	conn, err := database.Connect()
	if err != nil {
		return err
	}

	model := models.Init(conn.DB).User
	data, err := model.Get(id)
	if err != nil {
		return err
	}

	res, err := json.Marshal(data)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, string(res))
	return nil
}

type SignInBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u *userApi) Signin(w http.ResponseWriter, r *http.Request) error {
	username, _, ok := r.BasicAuth()
	if !ok {
		fmt.Fprintf(w, `{"success": false}`)
	}

	conn, err := database.Connect()
	if err != nil {
		return err
	}
	model := models.Init(conn.DB).User
	id, err := strconv.Atoi(username)
	if err != nil {
		return err
	}
	data, err := model.Get(id)
	if err != nil {
		return err
	}

	res, err := json.Marshal(data)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, string(res))
	return nil
}

func (u *userApi) GetPrava(w http.ResponseWriter, r *http.Request) error {

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		return err
	}

	conn, err := database.Connect()
	if err != nil {
		return err
	}
	model := models.Init(conn.DB).User
	data, err := model.GetPrava(id)
	if err != nil {
		return err
	}

	res, err := json.Marshal(data)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, string(res))
	return nil
}

func (u *userApi) GetUch(w http.ResponseWriter, r *http.Request) error {

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		return err
	}

	contr := controller.Init()
	data, err := contr.User.GetUch(id)

	res, err := json.Marshal(data)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, string(res))
	return nil
}
