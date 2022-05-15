package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"pmain2/internal/apperror"
	"pmain2/internal/types"
	"pmain2/pkg/utils"
	"strings"
	"time"
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

type FindUchetS struct {
	Id            int            `json:"id"`
	Date          string         `json:"date"`
	Category      int            `json:"category"`
	CategoryS     string         `json:"categoryS"`
	Reason        string         `json:"reason"`
	ReasonS       string         `json:"reasonS"`
	Diagnose      string         `json:"diagnose"`
	DiagnoseS     string         `json:"diagnoseS"`
	DiagnoseSNull sql.NullString `json:"-"`
	DockId        int            `json:"dockId"`
	DockNameNull  sql.NullString `json:"-"`
	DockName      string         `json:"dockName"`
	Section       int            `json:"section"`
}

func (m *patientModel) FindUchet(patientId, first, skip int) (*[]FindUchetS, error) {
	sql := fmt.Sprintf(`select FIRST %v SKIP %v nz, datp, m.categ, kat, rr, prich, diagt, diagts, trim(reg_doct), dock, uch from find_uchet_m(%v) m
order by datp desc,  nz DESC`, first, skip, patientId)
	INFO.Println(sql)
	rows, err := m.DB.Query(sql)
	if err != nil {
		return nil, err
	}
	var data []FindUchetS
	for rows.Next() {
		r := FindUchetS{}
		err = rows.Scan(&r.Id, &r.Date, &r.Category, &r.CategoryS, &r.Reason, &r.ReasonS, &r.Diagnose, &r.DiagnoseSNull, &r.DockId, &r.DockNameNull, &r.Section)
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
		r.DiagnoseS = r.DiagnoseSNull.String
		r.DiagnoseS, err = utils.ToUTF8(strings.Trim(r.DiagnoseS, " "))
		if err != nil {
			return nil, err
		}
		r.DockName = r.DockNameNull.String
		r.DockName, err = utils.ToUTF8(strings.Trim(r.DockName, " "))
		if err != nil {
			return nil, err
		}

		dateUchet, _ := time.Parse(time.RFC3339, r.Date)
		r.Date = utils.ToDate(dateUchet)

		data = append(data, r)
	}

	return &data, nil
}

func (m *patientModel) FindLastUchet(patientId int) (*FindUchetS, error) {
	data, err := m.FindUchet(patientId, 1, 0)
	if err != nil {
		return nil, err
	}
	res := *data
	var r FindUchetS
	if len(res) > 0 {
		r = res[0]
	}

	return &r, nil
}

func (m *patientModel) NewVisit(visit types.NewVisit, tx *sql.Tx) (sql.Result, error) {
	sql := fmt.Sprintf(`insert into visit(
      PATIENT_ID, V_DATE, name_doct, 
      MASKA1, MASKA2, MASKA3, 
      DIAGNose, UCH_PID, UCH_Dock, 
      upd_who, upd_date, BDAY,maska4)
      values(%v, '%s', '%v', 
      %v, %v, %v, 
      '%s', '%v', %v, 
      %v, '%s', '%s', %v)`,
		visit.PatientId, visit.Date, visit.DockId,
		visit.Visit, visit.Unit, 0,
		visit.Diagnose, visit.Uch, visit.Uch,
		visit.DockId, time.Now().Format("2006-01-02"), visit.PatientBDay, 0,
	)
	INFO.Println(sql)
	result, err := tx.Exec(sql)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (m *patientModel) NewSRC(spc *types.NewSRC, tx *sql.Tx) (sql.Result, error) {
	sql := fmt.Sprintf(`insert into detstvo_src(patient_id, date_add, id_dock, podr, zakl)
values(%v, '%s', %v, %v, %v)`, spc.PatientId, spc.DateAdd, spc.DockId, spc.Unit, spc.Zakl)
	INFO.Println(sql)
	result, err := tx.Exec(sql)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (m *patientModel) IsVisited(visit *types.NewVisit) (bool, error) {
	sql := fmt.Sprintf(
		`SELECT kol from KONTR_VISIT(%v, %v, %v, '%s')`, visit.DockId, visit.Uch, visit.PatientId, visit.Date)
	INFO.Println(sql)
	row := m.DB.QueryRow(sql)
	count := false
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}

	return count, nil

}

func (m *patientModel) NewProf(visit types.NewProf, tx *sql.Tx) (sql.Result, error) {
	sql := fmt.Sprintf(`insert into visit(
      PATIENT_ID, V_DATE, name_doct, 
      MASKA1, MASKA2, MASKA3, 
      DIAGNose, UCH_PID, UCH_Dock, 
      upd_who, upd_date, BDAY,maska4)
      values(%v, '%s', '%v', 
      %v, %v, %v, 
      '%s', '%v', %v, 
      %v, '%s', '%s', %v)`,
		306258, visit.Date, visit.DockId,
		1024, visit.Unit, 0,
		"Z", visit.Uch, visit.Uch,
		visit.DockId, time.Now().Format("2006-01-02"), "25.10.1917", 0,
	)
	INFO.Println(sql)
	result, err := tx.Exec(sql)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (m *patientModel) GetCountJudgment(patientId int) (int, error) {
	sql := fmt.Sprintf(`select count(*) from prinud_m
where patient_id = %v
and nom_z = (select max(nom_z) from prinud_m where patient_id = %v)
and ((exit_date is null)or(exit_date < '01.01.1950'))`, patientId, patientId)
	INFO.Println(sql)
	row := m.DB.QueryRow(sql)
	count := 0
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (m *patientModel) IsInHospital(patientId int) (bool, error) {
	sqlQuery := fmt.Sprintf(`select count(*) from kart_m
where patient_id = %v 
and ((exit_date is null)or(exit_date = '30.12.1899'))`, patientId)
	INFO.Println(sqlQuery)
	row := m.DB.QueryRow(sqlQuery)
	count := 0
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}

	return false, nil
}

func (m *patientModel) GetCountRegDataInDate(patientId, section int, date time.Time) (int, error) {
	d1 := utils.ToDate(date)
	d2 := utils.ToDate(date.AddDate(0, 0, 1))
	sqlQuery := fmt.Sprintf(`select count(*) from registrat_m 
where patient_id = %v 
and reg_date between '%s' and '%s'  
and RTRIM(reg_reas) not in ('S011','P001','001') 
and RTRIM(sec_tion) = %v`, patientId, d1, d2, section)
	INFO.Println(sqlQuery)
	row := m.DB.QueryRow(sqlQuery)
	count := 0
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (m *patientModel) InsertReg(reg types.NewRegister, tx *sql.Tx) (sql.Result, error) {
	sqlQuery := fmt.Sprintf(`insert into registrat_m(
PATIENT_ID, SEC_TION, REG_DATE,
REG_DOCT, REG_REAS, CATEG_UCH,
DIAGNOS, INS_WHO)
values(%v, %v, '%s',
		%v, '%s', %v,
		'%s', %v)`,
		reg.PatientId, reg.Section, reg.Date,
		reg.DockId, reg.Reason, reg.Category,
		reg.Diagnose, reg.DockId)

	return tx.Exec(sqlQuery)
}

func (m *patientModel) UpdPatientVisible(patientId, typeVisible int, tx *sql.Tx) (sql.Result, error) {
	sqlQuery := fmt.Sprintf(`update general set bl_group = %v where patient_id = %v`, typeVisible, patientId)

	return tx.Exec(sqlQuery)
}
