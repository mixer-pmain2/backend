package controller

import (
	"pmain2/internal/models"
	"pmain2/internal/types"
	"time"
)

type administration struct{}

func initAdministrationController() *administration {
	return &administration{}
}

func (a *administration) DoctorLocation(location *types.NewDoctorLocation) (int, error) {
	model := models.Model.Administration

	err, tx := models.Model.CreateTx()
	if err != nil {
		return 21, err
	}
	defer tx.Rollback()

	_, err = model.DisableSections(location.Unit, tx)
	if err != nil {
		tx.Rollback()
		return 502, err
	}
	date, err := time.Parse("2006-01-02", location.Date)
	if err != nil {
		return 503, err
	}
	_, err = model.DeleteSectionsByDate(location.Unit, date, tx)
	if err != nil {
		tx.Rollback()
		return 504, err
	}

	if len(location.Data) == 0 {
		return 505, err
	}

	for _, data := range location.Data {
		if data.DoctorId > 0 {
			_, err = model.DoctorLocation(location.Unit, date, data, tx)
			if err != nil {
				tx.Rollback()
				return 501, err
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		return 22, err
	}

	return 0, nil
}

func (a *administration) DoctorLeadSection(location *types.NewDoctorLeadSection) (int, error) {
	model := models.Model.Administration

	err, tx := models.Model.CreateTx()
	if err != nil {
		return 21, err
	}
	defer tx.Rollback()

	d1 := time.Date(location.Year, time.Month(location.Month), 1, 0, 0, 0, 0, time.UTC)
	d2 := d1.AddDate(0, 1, -1)
	if err != nil {
		return 503, err
	}
	_, err = model.DeleteLeadSectionsByDate(location.Unit, d1, d2, tx)
	if err != nil {
		tx.Rollback()
		return 504, err
	}

	if len(location.Data) == 0 {
		return 505, err
	}

	for _, data := range location.Data {
		if data.DoctorId > 0 {
			_, err = model.DoctorLeadSection(location.Unit, d1, d2, data, tx)
			if err != nil {
				tx.Rollback()
				return 501, err
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		return 22, err
	}

	return 0, nil
}
