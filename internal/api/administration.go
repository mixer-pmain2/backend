package api

import (
	"encoding/json"
	"net/http"
	"pmain2/internal/consts"
	"pmain2/internal/controller"
	"pmain2/internal/types"
)

type administrationApi struct{}

func administrationApiInit() *administrationApi {
	return &administrationApi{}
}

func (a *administrationApi) DoctorLocation(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return nil
	}
	var newLocation types.NewDoctorLocation
	err := json.NewDecoder(r.Body).Decode(&newLocation)
	if err != nil {
		return err
	}

	c := controller.Init()
	val, err := c.Administration.DoctorLocation(&newLocation)
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
