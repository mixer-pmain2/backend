package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pmain2/internal/controller"
	session2 "pmain2/internal/session"
	"pmain2/internal/types"
	"pmain2/pkg/utils"
	"strconv"

	"github.com/gorilla/mux"

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

	err, tx := models.Model.CreateTx()
	if err != nil {
		ERROR.Println(err)
		return err
	}
	defer tx.Rollback()

	model := models.Model.User
	data, err := model.Get(id, tx)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		ERROR.Println(err)
		return err
	}

	res, err := json.Marshal(data)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, string(res))
	return nil
}

func (u *userApi) Signin(w http.ResponseWriter, r *http.Request) error {
	username, _, ok := r.BasicAuth()
	if !ok {
		fmt.Fprintf(w, `{"success": false}`)
	}

	err, tx := models.Model.CreateTx()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	model := models.Model.User
	id, err := strconv.Atoi(username)
	if err != nil {
		return err
	}
	data, err := model.Get(id, tx)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		ERROR.Println(err)
		return err
	}

	res, err := json.Marshal(data)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, string(res))
	return nil
}

func (u *userApi) Login(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return nil
	}

	login := types.Login{}
	getParams(r, &login)

	c := controller.Init()
	var err error
	user, _ := utils.ToWin1251(login.Username)
	pass, _ := utils.ToWin1251(login.Password) // utils.ToASCII(password)
	isAuth, err := c.User.IsAuth(user, pass)
	if err != nil {
		ERROR.Println(err)
		return err
	}
	session, _ := session2.Store.Get(r, "user")
	if isAuth {
		session.Values["isAuth"] = true

		err, tx := models.Model.CreateTx()
		if err != nil {
			ERROR.Println(err)
			return err
		}
		defer tx.Rollback()
		model := models.Model.User
		id, err := strconv.Atoi(login.Username)
		if err != nil {
			return err
		}
		data, err := model.Get(id, tx)
		if err != nil {
			return err
		}
		err = tx.Commit()
		if err != nil {
			ERROR.Println(err)
			return err
		}
		session.Values["lname"] = data.Lname
		session.Values["fname"] = data.Fname
		session.Values["sname"] = data.Sname
		session.Values["id"] = data.Id

		err = session.Save(r, w)
		if err != nil {
			return err
		}
		res := types.HttpResponse{Success: true, Error: 0}
		res.Data = data
		mRes, err := json.Marshal(res)
		fmt.Fprintf(w, string(mRes))
		return nil
	}
	session.Values["isAuth"] = false
	session.Save(r, w)

	fmt.Fprintf(w, string(resSuccess(1)))
	return nil
}

func (u *userApi) GetPrava(w http.ResponseWriter, r *http.Request) error {

	params := getParams(r, nil)

	contr := controller.Init()
	data, err := contr.User.GetPrava(params.id, params.isCache)
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

func (u *userApi) ChangePassword(w http.ResponseWriter, r *http.Request) error {

	data := types.ChangePassword{}
	params := getParams(r, &data)
	data.UserId = int64(params.id)

	contr := controller.Init()
	fmt.Println(data)
	val, err := contr.User.ChangePassword(data)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, string(resSuccess(val)))
	return nil
}
