package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"net/url"
	"pmain2/internal/apperror"
	"pmain2/internal/consts"
	"pmain2/internal/controller"
	"pmain2/internal/types"
	"strconv"
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
	c := controller.Init()
	pData, err := c.Patient.FindById(id)
	if pData != nil {
		res, err := json.Marshal(pData)
		if err != nil {
			return err
		}
		fmt.Fprintf(w, string(res))
		return nil
	}

	return apperror.ErrDataNotFound
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

func (p *patientApi) FindUchet(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return nil
	}
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		return err
	}

	c := controller.Init()
	data, err := c.Patient.FindUchet(id)
	if err != nil {
		return err
	}

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

func (p *patientApi) HistoryVisits(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return nil
	}
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		return err
	}
	isCache, err := strconv.ParseBool(r.URL.Query().Get("cache"))
	if err != nil {
		isCache = true
	}

	c := controller.Init()
	data, err := c.Patient.HistoryVisits(id, isCache)
	if err != nil {
		return err
	}

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

func (p *patientApi) HistoryHospital(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return nil
	}
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		return err
	}

	c := controller.Init()
	data, err := c.Patient.HistoryHospital(id)
	if err != nil {
		return err
	}

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

func (p *patientApi) NewVisit(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return nil
	}
	var newVisit types.NewVisit
	err := json.NewDecoder(r.Body).Decode(&newVisit)
	if err != nil {
		return err
	}

	c := controller.Init()
	val, err := c.Patient.NewVisit(&newVisit)
	if err != nil {
		return err
	}

	res := types.HttpResponse{Success: true, Error: 0}

	if val > 0 {
		res.Success = false
		res.Error = val
		res.Message = consts.ArrErrors[val]
	}
	resMarshal, _ := json.Marshal(res)
	w.Write(resMarshal)
	return nil
}
