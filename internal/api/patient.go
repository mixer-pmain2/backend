package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"net/url"
	"pmain2/internal/controller"
	"strconv"

	"pmain2/internal/database"
	"pmain2/internal/models"
)

type patientApi struct{}

func patientApiInit() *patientApi {
	return &patientApi{}
}

func (p *patientApi) Get(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
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
	defer conn.Close()

	model := models.CreatePatient(conn.DB)
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

func (p *patientApi) Find(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		return nil
	}
	lname, err := url.QueryUnescape(r.URL.Query().Get("lname"))
	if err != nil {
		return err
	}
	fname, err := url.QueryUnescape(r.URL.Query().Get("fname"))
	if err != nil {
		return err
	}
	sname, err := url.QueryUnescape(r.URL.Query().Get("sname"))
	if err != nil {
		return err
	}

	c := controller.Init()
	data, err := c.Patient.FindByFio(lname, fname, sname)

	if len(*data) == 0 {
		fmt.Fprintf(w, "[]")
		return nil
	}

	marshal, err := json.Marshal(data)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, string(marshal))
	return nil
}
