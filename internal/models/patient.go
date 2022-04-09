package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"pmain2/internal/apperror"
	"pmain2/pkg/utils"
	"strings"
)

type patientModel struct {
	DB *sql.DB
}

func createPatient(db *sql.DB) *patientModel {
	return &patientModel{DB: db}
}

func (m *patientModel) Get(id int) (*Patient, error) {
	data := Patient{}
	sql := fmt.Sprintf(`select patient_id, lname, fname, sname, bday, bl_group, sex, job, adres
from general g, find_adres(g.patient_id,0) where patient_id=%v`, id)
	INFO.Println(sql)
	row := m.DB.QueryRow(sql)

	err := row.Scan(
		&data.Id, &data.Lname, &data.Fname, &data.Sname,
		&data.Bday, &data.Visibility, &data.Sex, &data.Snils, &data.Address)
	if err != nil {
		ERROR.Println(err.Error())
		return nil, err
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
	INFO.Println(string(res))

	return &data, nil
}

func (m *patientModel) FindByFIO(lname, fname, sname string) (*[]Patient, error) {
	var data = []Patient{}
	sql := fmt.Sprintf(
		`select patient_id, lname, fname, sname, bday, bl_group, sex, job from general
				where lname like ?
				and fname like ?
				and sname like ?
				order by lname, fname, sname`)
	INFO.Println(sql)
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

func (m *patientModel) GetAddress(id int) (string, error) {
	var address string
	sql := fmt.Sprintf("select adres from find_adres(%v,0)", id)
	row := m.DB.QueryRow(sql)
	err := row.Scan(&address)
	if err != nil {
		return "", err
	}
	return address, nil
}

type FindDispS struct {
	VisitId       int    `json:"visitId"`
	Date          string `json:"date"`
	DockName      string `json:"dockName"`
	Diag          string `json:"diag"`
	DiagS         string `json:"diagS"`
	Reason        string `json:"reason"`
	TypeVisit     string `json:"typeVisit"`
	TypeVisitCode int    `json:"typeVisitCode"`
	Where         int    `json:"where"`
}

func (m *patientModel) FindDisp(id int) (*[]FindDispS, error) {
	var data []FindDispS
	sql := fmt.Sprintf(`select v_n, dat, dokf, diag, diag_t, prich, sost, m1, m2 from find_disp(%v)
     order by dat DESC`, id)
	rows, err := m.DB.Query(sql)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		row := FindDispS{}
		err := rows.Scan(&row.VisitId, &row.Date, &row.DockName, &row.Diag, &row.DiagS, &row.Reason, &row.TypeVisit,
			&row.TypeVisitCode, &row.Where)
		if err != nil {
			return nil, err
		}
		data = append(data, row)
	}
	return &data, nil
}

type FindUchetS struct {
	Id           int            `json:"id"`
	Date         string         `json:"date"`
	Category     int            `json:"category"`
	CategoryS    string         `json:"categoryS"`
	Reason       string         `json:"reason"`
	ReasonS      string         `json:"reasonS"`
	Diagnose     string         `json:"diagnose"`
	DiagnoseS    string         `json:"diagnoseS"`
	DockId       int            `json:"dockId"`
	DockNameNull sql.NullString `json:"-"`
	DockName     string         `json:"dockName"`
	Section      int            `json:"section"`
}

func (m *patientModel) FindUchet(patientId int) (*[]FindUchetS, error) {
	sql := fmt.Sprintf(`select nz, datp, m.categ, kat, rr, prich, diagt, diagts, trim(reg_doct), dock, uch from find_uchet_m(%v) m
order by datp desc,  nz DESC`, patientId)
	INFO.Println(sql)
	rows, err := m.DB.Query(sql)
	if err != nil {
		return nil, err
	}
	var data []FindUchetS
	for rows.Next() {
		r := FindUchetS{}
		err = rows.Scan(&r.Id, &r.Date, &r.Category, &r.CategoryS, &r.Reason, &r.ReasonS, &r.Diagnose, &r.DiagnoseS, &r.DockId, &r.DockNameNull, &r.Section)
		if err != nil {
			return nil, err
		}
		r.CategoryS, err = utils.ToUTF8(r.CategoryS)
		if err != nil {
			return nil, err
		}
		r.ReasonS, err = utils.ToUTF8(strings.Trim(r.ReasonS, " "))
		if err != nil {
			return nil, err
		}
		r.Reason = strings.Trim(r.Reason, " ")
		r.DiagnoseS, err = utils.ToUTF8(strings.Trim(r.DiagnoseS, " "))
		if err != nil {
			return nil, err
		}
		r.DockName = r.DockNameNull.String
		r.DockName, err = utils.ToUTF8(strings.Trim(r.DockName, " "))
		if err != nil {
			return nil, err
		}

		data = append(data, r)
	}

	return &data, nil
}
