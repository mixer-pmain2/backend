package controller

import (
	"fmt"
	cacheName "pmain2/internal/cache"
	"pmain2/internal/models"
	"pmain2/internal/types"
	"pmain2/pkg/cache"
	"time"
)

var (
	cacheDoctor = cache.CreateCache(time.Minute, time.Minute)
)

type doctor struct{}

func initDoctorController() *doctor {
	return &doctor{}
}

func (p *doctor) GetRate(data types.DoctorFindParams, isCache bool) (*[]types.DoctorRate, error) {
	cacheName := cacheName.DoctorRate(data.DoctorId, data.Year, data.Month, data.Unit)

	item, ok := cache.AppCache.Get(cacheName)
	if ok && isCache {
		return item.(*[]types.DoctorRate), nil
	}

	model := models.Model.Doctor
	err, tx := models.Model.CreateTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	res, err := model.GetRate(data, tx)
	if err != nil {
		tx.Rollback()
		ERROR.Println(err)
		return nil, err
	}

	cache.AppCache.Set(cacheName, res, time.Hour)
	return res, nil
}

func (p *doctor) VisitCountPlan(data types.DoctorFindParams, isCache bool) (*[]types.DoctorVisitCountPlan, error) {
	cacheName := fmt.Sprintf("doctor_visit_count_plan_%v_%v_%v_%v", data.DoctorId, data.Year, data.Month, data.Unit)

	item, ok := cache.AppCache.Get(cacheName)
	if ok && isCache {
		return item.(*[]types.DoctorVisitCountPlan), nil
	}

	err, tx := models.Model.CreateTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	model := models.Model.Doctor
	res, err := model.VisitCountPlan(data, tx)
	if err != nil {
		tx.Rollback()
		ERROR.Println(err)
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	cache.AppCache.Set(cacheName, res, time.Hour)
	return res, nil
}

func (p *doctor) GetUnits(data types.DoctorFindParams, isCache bool) (*[]int, error) {
	cacheName := fmt.Sprintf("get_unit_%v_%v", data.DoctorId, data.Unit)

	item, ok := cache.AppCache.Get(cacheName)
	if ok && isCache {
		return item.(*[]int), nil
	}

	err, tx := models.Model.CreateTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	model := models.Model.Doctor
	res, err := model.GetUnits(data, tx)
	if err != nil {
		ERROR.Println(err)
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	cache.AppCache.Set(cacheName, res, time.Minute)
	return res, nil
}

func (p *doctor) UpdRate(data types.DoctorQueryUpdRate) (int, error) {

	err, tx := models.Model.CreateTx()
	if err != nil {
		return 20, err
	}
	defer tx.Rollback()

	model := models.Model.Doctor
	res, err := model.GetRate(types.DoctorFindParams{
		data.DoctorId,
		data.Month,
		data.Year,
		data.Unit,
	}, tx)
	if err != nil {
		tx.Rollback()
		ERROR.Println(err)
		return 21, err
	}
	if len(*res) > 0 {
		_, err = model.UpdRate(data, tx)
		if err != nil {
			tx.Rollback()
			ERROR.Println(err)
			return -1, err
		}
	} else {
		_, err = model.AddRate(data, tx)
		if err != nil {
			tx.Rollback()
			ERROR.Println(err)
			return -1, err
		}
	}
	err = tx.Commit()
	if err != nil {
		return 22, err
	}
	cache.AppCache.Delete(cacheName.DoctorRate(data.DoctorId, data.Year, data.Month, data.Unit))
	return 0, nil
}

func (p *doctor) DelRate(data types.DoctorQueryUpdRate) (int, error) {

	err, tx := models.Model.CreateTx()
	if err != nil {
		return 20, err
	}
	defer tx.Rollback()

	model := models.Model.Doctor
	_, err = model.DelRate(data, tx)
	if err != nil {
		tx.Rollback()
		ERROR.Println(err)
		return -1, err
	}
	err = tx.Commit()
	if err != nil {
		return 22, err
	}

	return 0, nil
}
