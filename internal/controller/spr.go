package controller

import (
	"pmain2/internal/database"
	"pmain2/internal/models"
	"pmain2/internal/types"
	"pmain2/pkg/cache"
	"pmain2/pkg/utils"
	"time"
)

type spr struct{}

func initSprController() *spr {
	return &spr{}
}

func (m *spr) GetPodr() (*map[int]string, error) {
	cacheName := "spr_podr"

	item, ok := cache.AppCache.Get(cacheName)
	if ok {
		res := item.(*map[int]string)
		return res, nil
	}

	conn, err := database.Connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	model := models.Init(conn.DB).Spr
	data, err := model.GetPodr()
	if err != nil {
		return nil, err
	}

	cache.AppCache.Set(cacheName, data, time.Hour)
	return data, nil
}

func (s *spr) GetPrava() (*[]models.PravaDict, error) {
	cacheName := "spr_prava"

	item, ok := cache.AppCache.Get(cacheName)
	if ok {
		res := item.(*[]models.PravaDict)
		return res, nil
	}

	conn, err := database.Connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	model := models.Init(conn.DB).Spr
	data, err := model.GetPrava()
	if err != nil {
		return nil, err
	}

	cache.AppCache.Set(cacheName, data, time.Hour)
	return data, nil
}

func (s *spr) GetSprVisit() (*map[int]string, error) {
	cacheName := "spr_visit"

	item, ok := cache.AppCache.Get(cacheName)
	if ok {
		res := item.(*map[int]string)
		return res, nil
	}

	conn, err := database.Connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	model := models.Init(conn.DB).Spr
	data, err := model.GetSprVisit()
	if err != nil {
		return nil, err
	}

	cache.AppCache.Set(cacheName, data, time.Hour)
	return data, nil
}

func (s *spr) GetDiags(diag string) (*[]models.DiagM, error) {
	cacheName := "spr_diag_" + diag

	item, ok := cache.AppCache.Get(cacheName)
	if ok {
		res := item.(*[]models.DiagM)
		return res, nil
	}

	conn, err := database.Connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	model := models.Init(conn.DB).Spr
	data, err := model.GetDiags(diag)
	if err != nil {
		return nil, err
	}

	cache.AppCache.Set(cacheName, data, time.Hour)
	return data, nil
}

func (s *spr) GetParams() (*[]models.ServiceM, error) {
	cacheName := "service_params"

	item, ok := cache.AppCache.Get(cacheName)
	if ok {
		res := item.(*[]models.ServiceM)
		return res, nil
	}

	conn, err := database.Connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	model := models.Init(conn.DB).Spr
	data, err := model.GetParams()
	if err != nil {
		return nil, err
	}

	*data = append(*data, []models.ServiceM{{
		Param:  "current_date",
		ParamS: utils.ToDate(time.Now()),
	}, {
		Param:  "registrat_min_date",
		ParamS: utils.ToDate(time.Date(time.Now().Year(), 1, 1, 0, 0, 0, 0, time.Local)),
	},
	}...)

	cache.AppCache.Set(cacheName, data, time.Minute*10)
	return data, nil
}

func (s *spr) GetSprReason() (*map[string]string, error) {
	cacheName := "spr_reason"

	item, ok := cache.AppCache.Get(cacheName)
	if ok {
		res := item.(*map[string]string)
		return res, nil
	}

	conn, err := database.Connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	model := models.Init(conn.DB).Spr
	data, err := model.GetSprReason()
	if err != nil {
		return nil, err
	}

	cache.AppCache.Set(cacheName, data, time.Hour)
	return data, nil
}

func (s *spr) GetSprInvalidKind() (*map[string]string, error) {
	cacheName := "spr_invalid_kind"

	item, ok := cache.AppCache.Get(cacheName)
	if ok {
		res := item.(*map[string]string)
		return res, nil
	}

	conn, err := database.Connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	model := models.Init(conn.DB).Spr
	data, err := model.GetSprInvalidKind()
	if err != nil {
		return nil, err
	}

	cache.AppCache.Set(cacheName, data, time.Hour)
	return data, nil
}

func (s *spr) GetSprInvalidChildAnomaly() (*map[string]string, error) {
	cacheName := "spr_invalid_child_anomaly"

	item, ok := cache.AppCache.Get(cacheName)
	if ok {
		res := item.(*map[string]string)
		return res, nil
	}

	conn, err := database.Connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	model := models.Init(conn.DB).Spr
	data, err := model.GetSprInvalidChildAnomaly()
	if err != nil {
		return nil, err
	}

	cache.AppCache.Set(cacheName, data, time.Hour)
	return data, nil
}

func (s *spr) GetSprInvalidChildLimit() (*map[string]string, error) {
	cacheName := "spr_invalid_child_limit"

	item, ok := cache.AppCache.Get(cacheName)
	if ok {
		res := item.(*map[string]string)
		return res, nil
	}

	conn, err := database.Connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	model := models.Init(conn.DB).Spr
	data, err := model.GetSprInvalidChildLimit()
	if err != nil {
		return nil, err
	}

	cache.AppCache.Set(cacheName, data, time.Hour)
	return data, nil
}

func (s *spr) GetSprInvalidReason() (*map[string]string, error) {
	cacheName := "spr_invalid_reason"

	item, ok := cache.AppCache.Get(cacheName)
	if ok {
		res := item.(*map[string]string)
		return res, nil
	}

	conn, err := database.Connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	model := models.Init(conn.DB).Spr
	data, err := model.GetSprInvalidReason()
	if err != nil {
		return nil, err
	}

	cache.AppCache.Set(cacheName, data, time.Hour)
	return data, nil
}

func (s *spr) GetSprCustodyWho() (*map[string]string, error) {
	cacheName := "spr_custody_who"

	item, ok := cache.AppCache.Get(cacheName)
	if ok {
		res := item.(*map[string]string)
		return res, nil
	}

	conn, err := database.Connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	model := models.Init(conn.DB).Spr
	data, err := model.GetSprCustodyWho()
	if err != nil {
		return nil, err
	}

	cache.AppCache.Set(cacheName, data, time.Hour)
	return data, nil
}

func (s *spr) FindRepublic(find *types.Find) (*[]types.Spr, error) {
	cacheName := "spr_republic_" + find.Name

	item, ok := cache.AppCache.Get(cacheName)
	if ok {
		res := item.(*[]types.Spr)
		return res, nil
	}

	conn, err := database.Connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	find.Name, err = utils.ToWin1251(find.Name)
	if err != nil {
		return nil, err
	}

	model := models.Init(conn.DB).Spr
	data, err := model.FindRepublic(find)
	if err != nil {
		return nil, err
	}

	cache.AppCache.Set(cacheName, data, time.Hour)
	return data, nil
}

func (s *spr) FindRegion(find *types.Find) (*[]types.Spr, error) {
	cacheName := "spr_region_" + find.Name

	item, ok := cache.AppCache.Get(cacheName)
	if ok {
		res := item.(*[]types.Spr)
		return res, nil
	}

	conn, err := database.Connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	find.Name, err = utils.ToWin1251(find.Name)
	if err != nil {
		return nil, err
	}

	model := models.Init(conn.DB).Spr
	data, err := model.FindRegion(find)
	if err != nil {
		return nil, err
	}

	cache.AppCache.Set(cacheName, data, time.Hour)
	return data, nil
}

func (s *spr) FindDistrict(find *types.Find) (*[]types.Spr, error) {
	cacheName := "spr_district_" + find.Name

	item, ok := cache.AppCache.Get(cacheName)
	if ok {
		res := item.(*[]types.Spr)
		return res, nil
	}

	conn, err := database.Connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	find.Name, err = utils.ToWin1251(find.Name)
	if err != nil {
		return nil, err
	}

	model := models.Init(conn.DB).Spr
	data, err := model.FindDistrict(find)
	if err != nil {
		return nil, err
	}

	cache.AppCache.Set(cacheName, data, time.Hour)
	return data, nil
}

func (s *spr) FindArea(find *types.Find) (*[]types.Spr, error) {
	cacheName := "spr_area_" + find.Name

	item, ok := cache.AppCache.Get(cacheName)
	if ok {
		res := item.(*[]types.Spr)
		return res, nil
	}

	conn, err := database.Connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	find.Name, err = utils.ToWin1251(find.Name)
	if err != nil {
		return nil, err
	}

	model := models.Init(conn.DB).Spr
	data, err := model.FindArea(find)
	if err != nil {
		return nil, err
	}

	cache.AppCache.Set(cacheName, data, time.Hour)
	return data, nil
}

func (s *spr) FindStreet(find *types.Find) (*[]types.Spr, error) {
	cacheName := "spr_street_" + find.Name

	item, ok := cache.AppCache.Get(cacheName)
	if ok {
		res := item.(*[]types.Spr)
		return res, nil
	}

	conn, err := database.Connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	find.Name, err = utils.ToWin1251(find.Name)
	if err != nil {
		return nil, err
	}

	model := models.Init(conn.DB).Spr
	data, err := model.FindStreet(find)
	if err != nil {
		return nil, err
	}

	cache.AppCache.Set(cacheName, data, time.Hour)
	return data, nil
}
