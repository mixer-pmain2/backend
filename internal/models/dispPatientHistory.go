package models

import (
	"database/sql"
	"fmt"
	"pmain2/pkg/utils"
	"strings"
	"time"
)

type HistoryVisit struct {
	Id       int    `json:"id"`
	Date     string `json:"date"`
	DockName string `json:"dockName"`
	Diag     string `json:"diag"`
	DiagS    string `json:"diagS"`
	Reason   string `json:"reason"`
	Where    string `json:"where"`
	Type     int    `json:"typeVisit"`
	Unit     int    `json:"unit"`
}

func (m *patientModel) HistoryVisits(id int) (*[]HistoryVisit, error) {
	var data []HistoryVisit
	sql := fmt.Sprintf(
		`select first 1000 v_n, dat, dokf, diag, diag_t, prich, sost, m1, m2 from find_disp(%v)
     order by dat DESC`, id)
	INFO.Println(sql)
	rows, err := m.DB.Query(sql)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		p := HistoryVisit{}
		err = rows.Scan(&p.Id, &p.Date, &p.DockName, &p.Diag, &p.DiagS, &p.Reason, &p.Where, &p.Type, &p.Unit)
		dateVisit, _ := time.Parse(time.RFC3339, p.Date)
		p.Date = utils.ToDate(dateVisit)
		p.DockName, _ = utils.ToUTF8(p.DockName)
		p.Diag, _ = utils.ToUTF8(strings.Trim(p.Diag, " "))
		p.DiagS, _ = utils.ToUTF8(strings.Trim(p.DiagS, " "))
		p.Reason, _ = utils.ToUTF8(p.Reason)
		p.Where, _ = utils.ToUTF8(p.Where)
		if err != nil {
			return nil, err
		}
		data = append(data, p)
	}

	return &data, nil

}

type HistoryHospital struct {
	Id            int            `json:"id"`
	DateStart     string         `json:"dateStart"`
	DateEnd       string         `json:"dateEnd"`
	DateEndNull   sql.NullString `json:"-"`
	DiagStart     string         `json:"diagStart"`
	DiagEnd       string         `json:"diagEnd"`
	DiagStartS    string         `json:"diagStartS"`
	DiagEndS      string         `json:"diagEndS"`
	DiagEndSNull  sql.NullString `json:"-"`
	Otd           string         `json:"otd"`
	HistoryNumber int            `json:"historyNumber"`
	Where         string         `json:"where"`
}

func (m *patientModel) HistoryHospital(id int) (*[]HistoryHospital, error) {
	var data []HistoryHospital
	sql := fmt.Sprintf(`select datp, datv, dp, dv, diagp, diagv, otd, ni, we, nom_z from find_stac_otkaz(%v) order by datp DESC`, id)
	INFO.Println(sql)
	rows, err := m.DB.Query(sql)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		p := HistoryHospital{}
		err = rows.Scan(&p.DateStart, &p.DateEndNull, &p.DiagStart, &p.DiagEnd, &p.DiagStartS, &p.DiagEndSNull, &p.Otd, &p.HistoryNumber, &p.Where, &p.Id)
		dateStart, _ := time.Parse(time.RFC3339, p.DateStart)
		p.DateStart = utils.ToDate(dateStart)
		dateEnd, _ := time.Parse(time.RFC3339, p.DateEndNull.String)
		p.DateEnd = utils.ToDate(dateEnd)
		if p.DateEndNull.String == "" {
			p.DateEnd = ""
		}
		p.DiagStart = strings.Trim(p.DiagStart, " ")
		p.DiagEnd = strings.Trim(p.DiagEnd, " ")
		p.DiagStartS, _ = utils.ToUTF8(strings.Trim(p.DiagStartS, " "))
		p.DiagEndS = p.DiagEndSNull.String
		p.DiagEndS, _ = utils.ToUTF8(strings.Trim(p.DiagEndS, " "))
		p.Where, _ = utils.ToUTF8(strings.Trim(p.Where, " "))
		if err != nil {
			return nil, err
		}
		data = append(data, p)
	}

	return &data, nil

}

type HistorySPC struct {
	Date string `json:"date"`
	Res  string `json:"res"`
}

func (m *patientModel) HistorySPC(patientId int, podr int) (*[]HistorySPC, error) {
	var data []HistorySPC
	sql := fmt.Sprintf(`select cast(ds.date_add as date) as date_, case
                    when ds.zakl = 0 then 'Согласие'
                    when ds.zakl = 1 then 'Отказ'
                    else 'ошибка'
                    end as res_ 
from detstvo_src ds
where ds.patient_id = %v
and ds.podr = %v
order by ds.date_add`, patientId, podr)
	INFO.Println(sql)
	rows, err := m.DB.Query(sql)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		p := HistorySPC{}
		err = rows.Scan(&p.Date, &p.Res)
		date, _ := time.Parse(time.RFC3339, p.Date)
		p.Date = utils.ToDate(date)
		if err != nil {
			return nil, err
		}
		data = append(data, p)
	}

	return &data, nil

}
