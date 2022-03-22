package controller

import (
	"pmain2/internal/database"
	"pmain2/internal/models"
	"pmain2/pkg/utils"
)

type patient struct{}

func initPatientController() *patient {
	return &patient{}
}

func (p *patient) FindByFio(lname, fname, sname string) (*[]models.Patient, error) {

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

	return data, nil
}
