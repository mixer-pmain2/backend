package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"pmain2/internal/controller"
	"pmain2/internal/types"
	"strconv"
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

func (s *sprApi) GetSprInvalidKind(w http.ResponseWriter, r *http.Request) error {
	c := controller.Init()
	data, err := c.Spr.GetSprInvalidKind()
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

func (s *sprApi) GetSprInvalidChildAnomaly(w http.ResponseWriter, r *http.Request) error {
	c := controller.Init()
	data, err := c.Spr.GetSprInvalidChildAnomaly()
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

func (s *sprApi) GetSprInvalidChildLimit(w http.ResponseWriter, r *http.Request) error {
	c := controller.Init()
	data, err := c.Spr.GetSprInvalidChildLimit()
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

func (s *sprApi) GetSprInvalidReason(w http.ResponseWriter, r *http.Request) error {
	c := controller.Init()
	data, err := c.Spr.GetSprInvalidReason()
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

func (s *sprApi) GetSprCustodyWho(w http.ResponseWriter, r *http.Request) error {
	c := controller.Init()
	data, err := c.Spr.GetSprCustodyWho()
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

func (s *sprApi) FindRepublic(w http.ResponseWriter, r *http.Request) error {
	var err error
	find := &types.Find{}
	find.Name, err = url.QueryUnescape(r.URL.Query().Get("name"))
	if err != nil {
		return err
	}

	c := controller.Init()
	data, err := c.Spr.FindRepublic(find)
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

func (s *sprApi) FindRegion(w http.ResponseWriter, r *http.Request) error {
	var err error
	find := &types.Find{}
	find.Name, err = url.QueryUnescape(r.URL.Query().Get("name"))
	if err != nil {
		return err
	}

	c := controller.Init()
	data, err := c.Spr.FindRegion(find)
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

func (s *sprApi) FindDistrict(w http.ResponseWriter, r *http.Request) error {
	var err error
	find := &types.Find{}
	find.Name, err = url.QueryUnescape(r.URL.Query().Get("name"))
	if err != nil {
		return err
	}

	c := controller.Init()
	data, err := c.Spr.FindDistrict(find)
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

func (s *sprApi) FindArea(w http.ResponseWriter, r *http.Request) error {
	var err error
	find := &types.Find{}
	find.Name, err = url.QueryUnescape(r.URL.Query().Get("name"))
	if err != nil {
		return err
	}

	c := controller.Init()
	data, err := c.Spr.FindArea(find)
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

func (s *sprApi) FindStreet(w http.ResponseWriter, r *http.Request) error {
	var err error
	find := &types.Find{}
	find.Name, err = url.QueryUnescape(r.URL.Query().Get("name"))
	if err != nil {
		return err
	}

	c := controller.Init()
	data, err := c.Spr.FindStreet(find)
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

func (s *sprApi) FindSections(w http.ResponseWriter, r *http.Request) error {
	var err error
	params := getParams(r, nil)
	find := &types.FindI{}
	podr, err := url.QueryUnescape(r.URL.Query().Get("unit"))
	if err != nil {
		return err
	}
	podrI, _ := strconv.Atoi(podr)
	find.Name = int64(podrI)

	c := controller.Init()
	data, err := c.Spr.FindSections(find, params.isCache)
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

func (s *sprApi) FindSectionDoctor(w http.ResponseWriter, r *http.Request) error {
	var err error
	params := getParams(r, nil)
	find := &types.FindI{}
	podr, err := url.QueryUnescape(r.URL.Query().Get("unit"))
	if err != nil {
		return err
	}
	podrI, _ := strconv.Atoi(podr)
	find.Name = int64(podrI)

	c := controller.Init()
	data, err := c.Spr.FindSectionDoctor(find, params.isCache)
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

func (s *sprApi) FindSectionLead(w http.ResponseWriter, r *http.Request) error {
	var err error
	params := getParams(r, nil)
	find := &types.FindDoctorLead{}
	unit, err := url.QueryUnescape(r.URL.Query().Get("unit"))
	month, err := url.QueryUnescape(r.URL.Query().Get("month"))
	year, err := url.QueryUnescape(r.URL.Query().Get("year"))
	if err != nil {
		return err
	}
	find.Unit, _ = strconv.Atoi(unit)
	find.Month, _ = strconv.Atoi(month)
	find.Year, _ = strconv.Atoi(year)

	c := controller.Init()
	data, err := c.Spr.FindSectionLead(find, params.isCache)
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

func (s *sprApi) GetDoctors(w http.ResponseWriter, r *http.Request) error {
	var err error
	params := getParams(r, nil)
	find := &types.FindI{}
	podr, err := url.QueryUnescape(r.URL.Query().Get("unit"))
	if err != nil {
		return err
	}
	podrI, _ := strconv.Atoi(podr)
	find.Name = int64(podrI)

	c := controller.Init()
	data, err := c.Spr.GetDoctors(find, params.isCache)
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
