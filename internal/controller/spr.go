package controller

import (
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

	model := models.Model.Spr
	err, tx := models.Model.CreateTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	data, err := model.GetPodr(tx)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		ERROR.Println(err)
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

	model := models.Model.Spr
	err, tx := models.Model.CreateTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	data, err := model.GetPrava(tx)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		ERROR.Println(err)
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

	model := models.Model.Spr
	err, tx := models.Model.CreateTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	data, err := model.GetSprVisit(tx)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		ERROR.Println(err)
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

	model := models.Model.Spr
	err, tx := models.Model.CreateTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	data, err := model.GetDiags(diag, tx)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		ERROR.Println(err)
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

	model := models.Model.Spr
	err, tx := models.Model.CreateTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	data, err := model.GetParams(tx)
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
	err = tx.Commit()
	if err != nil {
		ERROR.Println(err)
		return nil, err
	}

	cache.AppCache.Set(cacheName, data, 0)
	return data, nil
}

func (s *spr) GetSprReason() (*map[string]string, error) {
	cacheName := "spr_reason"

	item, ok := cache.AppCache.Get(cacheName)
	if ok {
		res := item.(*map[string]string)
		return res, nil
	}

	model := models.Model.Spr
	err, tx := models.Model.CreateTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	data, err := model.GetSprReason(tx)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		ERROR.Println(err)
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

	model := models.Model.Spr
	err, tx := models.Model.CreateTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	data, err := model.GetSprInvalidKind(tx)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		ERROR.Println(err)
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

	model := models.Model.Spr
	err, tx := models.Model.CreateTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	data, err := model.GetSprInvalidChildAnomaly(tx)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		ERROR.Println(err)
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

	model := models.Model.Spr
	err, tx := models.Model.CreateTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	data, err := model.GetSprInvalidChildLimit(tx)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		ERROR.Println(err)
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

	model := models.Model.Spr
	err, tx := models.Model.CreateTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	data, err := model.GetSprInvalidReason(tx)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		ERROR.Println(err)
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

	model := models.Model.Spr
	err, tx := models.Model.CreateTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	data, err := model.GetSprCustodyWho(tx)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		ERROR.Println(err)
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

	err, tx := models.Model.CreateTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	find.Name, err = utils.ToWin1251(find.Name)
	if err != nil {
		return nil, err
	}

	model := models.Model.Spr
	data, err := model.FindRepublic(find, tx)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		ERROR.Println(err)
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

	err, tx := models.Model.CreateTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	find.Name, err = utils.ToWin1251(find.Name)
	if err != nil {
		return nil, err
	}

	model := models.Model.Spr
	data, err := model.FindRegion(find, tx)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		ERROR.Println(err)
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

	err, tx := models.Model.CreateTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	find.Name, err = utils.ToWin1251(find.Name)
	if err != nil {
		return nil, err
	}

	model := models.Model.Spr
	data, err := model.FindDistrict(find, tx)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		ERROR.Println(err)
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

	err, tx := models.Model.CreateTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	find.Name, err = utils.ToWin1251(find.Name)
	if err != nil {
		return nil, err
	}

	model := models.Model.Spr
	data, err := model.FindArea(find, tx)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		ERROR.Println(err)
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

	err, tx := models.Model.CreateTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	find.Name, err = utils.ToWin1251(find.Name)
	if err != nil {
		return nil, err
	}

	model := models.Model.Spr
	data, err := model.FindStreet(find, tx)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		ERROR.Println(err)
		return nil, err
	}

	cache.AppCache.Set(cacheName, data, time.Hour)
	return data, nil
}

func (s *spr) FindSections(find *types.FindI, isCache bool) (*[]types.SprUchN, error) {
	cacheName := "spr_section_" + string(find.Name)

	item, ok := cache.AppCache.Get(cacheName)
	if ok && isCache {
		res := item.(*[]types.SprUchN)
		return res, nil
	}

	model := models.Model.Spr
	err, tx := models.Model.CreateTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	data, err := model.FindSections(find, tx)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		ERROR.Println(err)
		return nil, err
	}

	cache.AppCache.Set(cacheName, data, time.Hour)
	return data, nil
}

func (s *spr) FindSectionDoctor(find *types.FindI, isCache bool) (*[]types.LocationDoctor, error) {
	cacheName := "find_section_doctor_by_unit_" + string(find.Name)

	item, ok := cache.AppCache.Get(cacheName)
	if ok && isCache {
		res := item.(*[]types.LocationDoctor)
		return res, nil
	}

	model := models.Model.Spr
	err, tx := models.Model.CreateTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	data, err := model.FindSectionDoctor(find, tx)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		ERROR.Println(err)
		return nil, err
	}

	cache.AppCache.Set(cacheName, data, time.Hour)
	return data, nil
}

func (s *spr) FindSectionLead(find *types.FindDoctorLead, isCache bool) (*[]types.LocationDoctor, error) {
	cacheName := "find_section_lead_doctor_by_unit_" + string(find.Unit) + string(find.Year) + string(find.Month)

	item, ok := cache.AppCache.Get(cacheName)
	if ok && isCache {
		res := item.(*[]types.LocationDoctor)
		return res, nil
	}

	model := models.Model.Spr
	err, tx := models.Model.CreateTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	data, err := model.FindSectionLead(find, tx)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		ERROR.Println(err)
		return nil, err
	}

	cache.AppCache.Set(cacheName, data, time.Hour)
	return data, nil
}

func (s *spr) GetDoctors(find *types.FindI, isCache bool) (*[]types.Doctor, error) {
	cacheName := "get_doctors_by_unit_" + string(find.Name)

	item, ok := cache.AppCache.Get(cacheName)
	if ok && isCache {
		res := item.(*[]types.Doctor)
		return res, nil
	}

	model := models.Model.Spr
	err, tx := models.Model.CreateTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	data, err := model.GetDoctors(find, tx)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		ERROR.Println(err)
		return nil, err
	}

	cache.AppCache.Set(cacheName, data, time.Hour)
	return data, nil
}
