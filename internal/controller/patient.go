package controller

import (
	"time"

	"pmain2/internal/database"
	"pmain2/internal/models"
	"pmain2/pkg/cache"
	"pmain2/pkg/utils"
)

var (
	cachePat = cache.CreateCache(time.Minute, time.Minute*5)
)

type patient struct{}

func initPatientController() *patient {
	return &patient{}
}

func (p *patient) FindByFio(lname, fname, sname string) (*[]models.Patient, error) {
	cacheName := lname + " " + fname + " " + sname

	item, ok := cachePat.Get(cacheName)
	if ok {
		res := item.(*[]models.Patient)
		return res, nil
	}

	conn, err := database.Connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	model := models.CreatePatient(conn.DB)
	lname, _ = utils.ToWin1251(lname)
	fname, _ = utils.ToWin1251(fname)
	sname, _ = utils.ToWin1251(sname)
	data, err := model.FindByFIO(lname, fname, sname)
	if err != nil {
		return nil, err
	}

	cachePat.Set(cacheName, data, 0)
	return data, nil
}

func (p *patient) FindById(id int) (*models.Patient, error) {

	conn, err := database.Connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	model := models.CreatePatient(conn.DB)
	data, err := model.Get(id)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (p *patient) FindUchet(id int) (*[]models.FindUchetS, error) {

	conn, err := database.Connect()
	if err != nil {
		ERROR.Println(err.Error())
		return nil, err
	}
	model := models.CreatePatient(conn.DB)
	data, err := model.FindUchet(id)
	if err != nil {
		ERROR.Println(err.Error())
		return nil, err
	}
	return data, nil
}
