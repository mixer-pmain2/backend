package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pmain2/internal/config"
	"pmain2/internal/controller"
	"pmain2/internal/types"
	"pmain2/pkg/utils"
	"pmain2/pkg/utils/jwt"
	"strconv"
	"strings"

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
	//username, _, ok := r.BasicAuth()
	//if !ok {
	//	fmt.Fprintf(w, `{"success": false}`)
	//}

	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	reqToken = splitToken[1]

	jwtT := jwt.JWT(config.AppConfig.SecretKey)
	user := jwtT.GetBody(reqToken)

	err, tx := models.Model.CreateTx()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	model := models.Model.User
	data, err := model.Get(user.UserId, tx)
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

	if isAuth {
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
		jwtT := jwt.JWT(config.AppConfig.SecretKey)
		token := jwtT.GetToken(jwt.Body{
			UserId: id,
		})
		userData := struct {
			Id    int64  `json:"id"`
			Lname string `json:"lname"`
			Fname string `json:"fname"`
			Sname string `json:"sname"`
			Token string `json:"token"`
		}{data.Id, data.Lname, data.Fname, data.Sname, token}

		if err != nil {
			return err
		}
		res := types.HttpResponse{Success: true, Error: 0}
		res.Data = userData
		mRes, err := json.Marshal(res)
		fmt.Fprintf(w, string(mRes))
		return nil
	}

	w.Header().Set("WWW-Authenticate", `Bearer error="invalid_token, charset="UTF-8"`)
	http.Error(w, "Unauthorized", http.StatusUnauthorized)
	//fmt.Fprintf(w, string(resSuccess(1)))
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
	val, err := contr.User.ChangePassword(data)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, string(resSuccess(val)))
	return nil
}
