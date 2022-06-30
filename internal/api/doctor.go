package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"pmain2/internal/apperror"
	"pmain2/internal/controller"
	"pmain2/internal/types"
	"strconv"
)

type doctorApi struct{}

func doctorApiInit() *doctorApi {
	return &doctorApi{}
}

func (d *doctorApi) GetRate(w http.ResponseWriter, r *http.Request) error {
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

	data := types.DoctorFindParams{DoctorId: id}
	data.Month, _ = strconv.Atoi(r.URL.Query().Get("month"))
	data.Year, _ = strconv.Atoi(r.URL.Query().Get("year"))
	data.Unit, _ = strconv.Atoi(r.URL.Query().Get("unit"))

	c := controller.Init()
	pData, err := c.Doctor.GetRate(data, isCache)
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

func (d *doctorApi) VisitCountPlan(w http.ResponseWriter, r *http.Request) error {
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

	data := types.DoctorFindParams{DoctorId: id}
	data.Month, _ = strconv.Atoi(r.URL.Query().Get("month"))
	data.Year, _ = strconv.Atoi(r.URL.Query().Get("year"))

	c := controller.Init()
	pData, err := c.Doctor.VisitCountPlan(data, isCache)
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

func (d *doctorApi) GetUnits(w http.ResponseWriter, r *http.Request) error {
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

	data := types.DoctorFindParams{DoctorId: id}
	data.Unit, _ = strconv.Atoi(r.URL.Query().Get("unit"))

	c := controller.Init()
	pData, err := c.Doctor.GetUnits(data, isCache)
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

func (d *doctorApi) UpdRate(w http.ResponseWriter, r *http.Request) error {
	data := types.DoctorQueryUpdRate{}
	params := getParams(r, &data)
	data.DoctorId = params.id

	c := controller.Init()
	val, err := c.Doctor.UpdRate(data)
	if err != nil && val < 0 {
		return err
	}

	w.Write(resSuccess(val))

	return nil
}

func (d *doctorApi) DelRate(w http.ResponseWriter, r *http.Request) error {
	data := types.DoctorQueryUpdRate{}
	params := getParams(r, &data)
	data.DoctorId = params.id

	c := controller.Init()
	val, err := c.Doctor.DelRate(data)
	if err != nil && val < 0 {
		return err
	}

	w.Write(resSuccess(val))

	return nil
}
