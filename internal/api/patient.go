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

	isCache, err := strconv.ParseBool(r.URL.Query().Get("cache"))
	if err != nil {
		isCache = true
	}

	c := controller.Init()
	pData, err := c.Patient.FindById(int64(id), isCache)
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

func (p *patientApi) New(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return nil
	}
	newPatient := types.NewPatient{}
	err := json.NewDecoder(r.Body).Decode(&newPatient)

	c := controller.Init()
	val, err, data := c.Patient.New(&newPatient)
	if err != nil {
		return err
	}

	type result struct {
		IsForced bool             `json:"isForced"`
		Data     *[]types.Patient `json:"data"`
		types.HttpResponse
	}
	_res := types.HttpResponse{Success: true, Error: 0}
	if val > 0 {
		_res.Success = false
		_res.Error = val
		_res.Message = consts.ArrErrors[val]
	}
	response := result{
		newPatient.IsForced,
		data,
		_res,
	}

	marshal, err := json.Marshal(response)
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

	isCache, err := strconv.ParseBool(r.URL.Query().Get("cache"))
	if err != nil {
		isCache = true
	}

	c := controller.Init()
	data, err := c.Patient.FindUchet(int64(id), isCache)
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

func (p *patientApi) NewReg(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return nil
	}
	var newReg types.NewRegister
	err := json.NewDecoder(r.Body).Decode(&newReg)
	if err != nil {
		return err
	}

	marshal, err := json.Marshal(newReg)
	if err != nil {
		return err
	}
	INFO.Println(string(marshal))

	c := controller.Init()
	val, err := c.Patient.NewReg(&newReg)
	if err != nil {
		return err
	}

	res := types.HttpResponse{Success: true, Error: 0}

	if val > 0 {
		res.Success = false
		res.Error = val
		res.Message = consts.ArrErrors[val]
		ERROR.Println(res)
	}

	marshal, err = json.Marshal(res)
	if err != nil {
		return err
	}

	INFO.Println(string(marshal))
	fmt.Fprintf(w, string(marshal))
	return nil
}

func (p *patientApi) NewRegTransfer(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return nil
	}
	var newReg types.NewRegisterTransfer
	err := json.NewDecoder(r.Body).Decode(&newReg)
	if err != nil {
		return err
	}

	marshal, err := json.Marshal(newReg)
	if err != nil {
		return err
	}
	INFO.Println(string(marshal))

	c := controller.Init()
	val, err := c.Patient.NewRegisterTransfer(&newReg)
	if err != nil {
		return err
	}

	res := types.HttpResponse{Success: true, Error: 0}

	if val > 0 {
		res.Success = false
		res.Error = val
		res.Message = consts.ArrErrors[val]
		ERROR.Println(res)
	}

	marshal, err = json.Marshal(res)
	if err != nil {
		return err
	}

	INFO.Println(string(marshal))
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

func (p *patientApi) GetAddress(w http.ResponseWriter, r *http.Request) error {
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
	data, err := c.Patient.GetAddress(int64(id), isCache)
	if err != nil {
		return err
	}

	res := struct {
		Address string `json:"address"`
	}{data}

	marshal, err := json.Marshal(res)
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

func (p *patientApi) NewProf(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return nil
	}
	var newProf types.NewProf
	err := json.NewDecoder(r.Body).Decode(&newProf)
	if err != nil {
		return err
	}

	c := controller.Init()
	val, err := c.Patient.NewProf(&newProf)
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

func (p *patientApi) GetSindrom(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
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
	data, err := c.Patient.HistorySindrom(id, isCache)
	if err != nil {
		return err
	}

	resMarshal, _ := json.Marshal(data)
	w.Write(resMarshal)
	return nil
}

func (p *patientApi) NewSindrom(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return nil
	}
	var newSindrom types.Sindrom
	err := json.NewDecoder(r.Body).Decode(&newSindrom)
	if err != nil {
		return err
	}

	c := controller.Init()
	val, err := c.Patient.NewSindrom(&newSindrom)
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

func (p *patientApi) RemoveSindrom(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return nil
	}
	var sindrom types.Sindrom
	err := json.NewDecoder(r.Body).Decode(&sindrom)
	if err != nil {
		return err
	}

	c := controller.Init()
	val, err := c.Patient.RemoveSindrom(&sindrom)
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

func (p *patientApi) FindInvalid(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
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
	data, err := c.Patient.FindInvalid(int64(id), isCache)
	if err != nil {
		return err
	}

	resMarshal, _ := json.Marshal(data)
	w.Write(resMarshal)
	return nil
}

func (p *patientApi) NewInvalid(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return nil
	}
	var newInvalid types.NewInvalid
	fmt.Println(r.Body)
	err := json.NewDecoder(r.Body).Decode(&newInvalid)
	if err != nil {
		return err
	}

	c := controller.Init()
	val, err := c.Patient.NewInvalid(&newInvalid)
	if err != nil && val < 0 {
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

func (p *patientApi) UpdInvalid(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return nil
	}
	var newInvalid types.NewInvalid
	fmt.Println(r.Body)
	err := json.NewDecoder(r.Body).Decode(&newInvalid)
	if err != nil {
		return err
	}

	c := controller.Init()
	val, err := c.Patient.UpdInvalid(&newInvalid)
	if err != nil && val < 0 {
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

func (p *patientApi) FindCustody(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
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
	data, err := c.Patient.FindCustody(int64(id), isCache)
	if err != nil {
		return err
	}

	resMarshal, _ := json.Marshal(data)
	w.Write(resMarshal)
	return nil
}
