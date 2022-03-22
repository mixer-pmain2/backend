package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"pmain2/internal/apperror"
)

type PatientModel struct {
	DB *sql.DB
}

func CreatePatient(db *sql.DB) *PatientModel {
	return &PatientModel{DB: db}
}

func (m *PatientModel) Get(id int) (*Patient, error) {
	data := Patient{}
	sql := fmt.Sprintf("select patient_id, lname, fname, sname, bday, bl_group, sex, job from general where patient_id=%v", id)
	INFO.Println(sql)
	row := m.DB.QueryRow(sql)

	err := row.Scan(&data.Id, &data.Lname, &data.Fname, &data.Sname, &data.Bday, &data.Visibility, &data.Sex, &data.Snils)
	if err != nil {
		ERROR.Println(err.Error())
		return nil, apperror.ErrDataNotFound
	}
	err = data.Serialize()
	if err != nil {
		ERROR.Println(err.Error())
		return nil, err
	}

	if data.Id == 0 {
		return nil, apperror.ErrDataNotFound
	}

	res, err := json.Marshal(&data)
	if err != nil {
		ERROR.Println(err.Error())
		return nil, err
	}
	INFO.Println(res)

	return &data, nil
}

func (m *PatientModel) FindByFIO(lname, fname, sname string) (*[]Patient, error) {
	var data = []Patient{}
	sql := fmt.Sprintf(
		`select patient_id, lname, fname, sname, bday, bl_group, sex, job from general
				where lname like ?
				and fname like ?
				and sname like ?
				order by lname, fname, sname`)
	stmt, err := m.DB.Prepare(sql)
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(lname+"%", fname+"%", sname+"%")
	if err != nil {
		return nil, err
	}
	p := Patient{}
	for rows.Next() {
		err = rows.Scan(&p.Id, &p.Lname, &p.Fname, &p.Sname, &p.Bday, &p.Visibility, &p.Sex, &p.Snils)
		if err != nil {
			return nil, err
		}
		p.Serialize()
		data = append(data, p)
	}

	return &data, nil
}
