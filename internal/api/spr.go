package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pmain2/internal/controller"
)

type sprApi struct{}

func sprApiInit() *sprApi {
	return &sprApi{}
}

func (s *sprApi) GetPodr(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return nil
	}
	contr := controller.Init()
	data, err := contr.Spr.GetPodr()
	res, err := json.Marshal(data)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, string(res))
	return nil

}

func (s *sprApi) GetPrava(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return nil
	}
	contr := controller.Init()
	data, err := contr.Spr.GetPrava()
	res, err := json.Marshal(data)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, string(res))
	return nil

}

func (s *sprApi) GetSprVisit(w http.ResponseWriter, r *http.Request) error {
	c := controller.Init()
	data, err := c.Spr.GetSprVisit()
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

func (s *sprApi) GetSprDiags(w http.ResponseWriter, r *http.Request) error {
	query := r.URL.Query()
	diag := query.Get("diag")

	c := controller.Init()
	data, err := c.Spr.GetDiags(diag)
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

func (s *sprApi) GetParams(w http.ResponseWriter, r *http.Request) error {
	c := controller.Init()
	data, err := c.Spr.GetParams()
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

func (s *sprApi) GetSprReasons(w http.ResponseWriter, r *http.Request) error {
	c := controller.Init()
	data, err := c.Spr.GetSprReason()
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
