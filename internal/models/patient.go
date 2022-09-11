package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"pmain2/internal/apperror"
	"pmain2/internal/consts"
	"pmain2/internal/types"
	"pmain2/pkg/utils"
	"strings"
	"time"
)

type patientModel struct {
}

func createPatient() *patientModel {
	return &patientModel{}
}

func (m *patientModel) GetMaxPatientId(tx *sql.Tx) (int64, error) {
	sqlQuery := "select max(patient_id) from general where patient_id < 1000000"
	row := tx.QueryRow(sqlQuery)
	INFO.Println(sqlQuery)
	var id int64
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *patientModel) New(p *types.NewPatient, tx *sql.Tx) (sql.Result, error) {
	sqlQuery := fmt.Sprintf(`insert into general(patient_id, LName, FName, SName, BDay, Sex, date_insert, who_insert)
values(?, ?, ?, ?, ?, ?, ?, ?)`)
	stmt, err := tx.Prepare(sqlQuery)
	defer stmt.Close()
	if err != nil {
		return nil, err
	}
	INFO.Println(sqlQuery)
	return stmt.Exec(p.PatientId, p.Lname, p.Fname, p.Sname, p.Bday, p.Sex, "NOW", p.UserId)
}

func (m *patientModel) Get(id int64, tx *sql.Tx) (*types.Patient, error) {
	data := types.Patient{}
	sql := fmt.Sprintf(`select patient_id, lname, fname, sname,
bday, bl_group, sex, job, adres,
passp_ser, passp_num, residence,
republic, region, district, pop_area, street, trim(house), trim(building), trim(flat), domicile
from general g, find_adres(g.patient_id,0) where patient_id=%v`, id)
	INFO.Println(sql)
	row := tx.QueryRow(sql)

	err := row.Scan(
		&data.Id, &data.Lname, &data.Fname, &data.Sname,
		&data.Bday, &data.Visibility, &data.Sex, &data.Snils, &data.Address,
		&data.PassportSeries, &data.PassportNumber, &data.Works,
		&data.Republic, &data.Region, &data.District, &data.Area, &data.Street, &data.House, &data.Build, &data.Flat, &data.Domicile)
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

func (m *patientModel) FindByFIO(lname, fname, sname string, tx *sql.Tx) (*[]types.Patient, error) {
	var data = make([]types.Patient, 0)
	sql := fmt.Sprintf(
		`select patient_id, lname, fname, sname, bday, bl_group, sex, job from general
				where lname like ?
				and fname like ?
				and sname like ?
				order by lname, fname, sname`)
	INFO.Println(sql)
	stmt, err := tx.Prepare(sql)
	defer stmt.Close()
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(lname+"%", fname+"%", sname+"%")
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	p := types.Patient{}
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

func (m *patientModel) FindByAddress(address types.Patient, tx *sql.Tx) (*[]types.Patient, error) {
	var data = make([]types.Patient, 0)
	sql := fmt.Sprintf(`select patient_id, lname, fname, sname, bday, bl_group, sex, job 
from general where republic = %v`, address.Republic)
	if address.Region > 0 {
		sql += fmt.Sprintf(` and region = %v`, address.Region)
	}
	if address.District > 0 {
		sql += fmt.Sprintf(` and district = %v`, address.District)
	}
	if address.Area > 0 {
		sql += fmt.Sprintf(` and pop_area = %v`, address.Area)
	}
	if address.Street > 0 {
		sql += fmt.Sprintf(` and street = %v`, address.Street)
	}
	if address.House != "" {
		sql += fmt.Sprintf(` and house = '%s'`, address.House)
	}
	if address.Build != "" {
		sql += fmt.Sprintf(` and building = '%s'`, address.Build)
	}
	sql += " order by lname, fname, sname"
	INFO.Println(sql)
	rows, err := tx.Query(sql)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	p := types.Patient{}
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

func (m *patientModel) GetAddress(id int64, tx *sql.Tx) (string, error) {
	var address string
	sql := fmt.Sprintf("select adres from find_adres(%v,0)", id)
	row := tx.QueryRow(sql)
	err := row.Scan(&address)
	if err != nil {
		return "", err
	}
	address, err = utils.ToUTF8(strings.Trim(address, " "))
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

func (m *patientModel) FindUchet(patientId int64, first, skip int, tx *sql.Tx) (*[]FindUchetS, error) {
	sql := fmt.Sprintf(`select FIRST %v SKIP %v nz, datp, m.categ, kat, rr, prich, diagt, diagts, trim(reg_doct), dock, uch from find_uchet_m(%v) m
order by datp desc,  nz DESC`, first, skip, patientId)
	INFO.Println(sql)
	rows, err := tx.Query(sql)
	defer rows.Close()
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

func (m *patientModel) FindLastUchet(patientId int64, tx *sql.Tx) (*FindUchetS, error) {
	data, err := m.FindUchet(patientId, 1, 0, tx)
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
      upd_who, upd_date, BDAY,maska4, uch_pid_i)
      values(%v, '%s', '%v', 
      %v, %v, %v, 
      '%s', '%v', %v, 
      %v, '%s', '%s', %v, %v)`,
		visit.PatientId, visit.Date, visit.DockId,
		visit.Visit, visit.Unit, 0,
		visit.Diagnose, visit.Uch, visit.Uch,
		visit.DockId, time.Now().Format("2006-01-02"), visit.PatientBDay, 0, visit.Uch,
	)
	INFO.Println(sql)
	return tx.Exec(sql)
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

func (m *patientModel) IsVisited(visit *types.NewVisit, tx *sql.Tx) (bool, error) {
	sql := fmt.Sprintf(
		`SELECT kol from KONTR_VISIT(%v, %v, %v, '%s')`, visit.DockId, visit.Uch, visit.PatientId, visit.Date)
	INFO.Println(sql)
	row := tx.QueryRow(sql)
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

func (m *patientModel) GetCountJudgment(patientId int64, tx *sql.Tx) (int, error) {
	sql := fmt.Sprintf(`select count(*) from prinud_m
where patient_id = %v
and nom_z = (select max(nom_z) from prinud_m where patient_id = %v)
and ((exit_date is null)or(exit_date < '01.01.1950'))`, patientId, patientId)
	INFO.Println(sql)
	row := tx.QueryRow(sql)
	count := 0
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (m *patientModel) IsInHospital(patientId int64, tx *sql.Tx) (bool, error) {
	sqlQuery := fmt.Sprintf(`select count(*) from kart_m
where patient_id = %v 
and ((exit_date is null)or(exit_date = '30.12.1899'))`, patientId)
	INFO.Println(sqlQuery)
	row := tx.QueryRow(sqlQuery)
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

func (m *patientModel) GetCountRegDataInDate(patientId int64, section int, date time.Time, tx *sql.Tx) (int, error) {
	d1 := utils.ToDate(date)
	d2 := utils.ToDate(date.AddDate(0, 0, 1))
	sqlQuery := fmt.Sprintf(`select count(*) from registrat_m 
where patient_id = %v 
and reg_date between '%s' and '%s'  
and RTRIM(reg_reas) not in ('S011','P001','001') 
and RTRIM(sec_tion) = %v`, patientId, d1, d2, section)
	INFO.Println(sqlQuery)
	row := tx.QueryRow(sqlQuery)
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

func (m *patientModel) UpdPatientVisible(patientId int64, typeVisible int, tx *sql.Tx) (sql.Result, error) {
	sqlQuery := fmt.Sprintf(`update general set bl_group = %v where patient_id = %v`, typeVisible, patientId)

	return tx.Exec(sqlQuery)
}

func (m *patientModel) NewSindrom(sindrom types.Sindrom, tx *sql.Tx) (sql.Result, error) {
	sql := fmt.Sprintf(`INSERT INTO SINDROM (patient_id, diag, ins_date, ins_dock) 
VALUES (%v, '%s', '%s', %v)`,
		sindrom.PatientId, sindrom.Diagnose,
		time.Now().Format("2006-01-02"), sindrom.DoctId,
	)
	INFO.Println(sql)
	result, err := tx.Exec(sql)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (m *patientModel) RemoveSindrom(sindrom types.Sindrom, tx *sql.Tx) (sql.Result, error) {
	sql := fmt.Sprintf(`delete from SINDROM where nom_z = %v`,
		sindrom.Id,
	)
	INFO.Println(sql)
	result, err := tx.Exec(sql)
	if err != nil {
		return nil, err
	}

	return result, nil
}

type FindInvalid struct {
	DateBegin  string `json:"dateBegin"`
	KindS      string `json:"kindS"`
	ReasonS    string `json:"reasonS"`
	DateChange string `json:"dateChange"`
	DateEnd    string `json:"dateEnd"`
	Id         int    `json:"id"`
}

func (m *patientModel) FindInvalid(patientId int64, tx *sql.Tx) (*[]FindInvalid, error) {
	sqlQuery := fmt.Sprintf(`select a.inv_begin, (select na_me from spr_visit_n where kod2 = a.inv_kind and kod1 = 6) as inv_kind,
(select na_me from spr_visit_n where kod2 = a.inv_reas and kod1 = 4) as inv_reas, a.inv_change, a.inv_end, a.nom_z
from invalid a where a.patient_id = %v
order by inv_begin DESC`, patientId)
	INFO.Println(sqlQuery)
	rows, err := tx.Query(sqlQuery)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	var data []FindInvalid
	for rows.Next() {
		r := FindInvalid{}
		var reason sql.NullString
		err = rows.Scan(&r.DateBegin, &r.KindS, &reason, &r.DateChange, &r.DateEnd, &r.Id)
		if err != nil {
			return nil, err
		}
		r.ReasonS = reason.String
		r.KindS, _ = utils.ToUTF8(strings.Trim(r.KindS, " "))
		r.ReasonS, _ = utils.ToUTF8(strings.Trim(r.ReasonS, " "))
		date, _ := time.Parse(time.RFC3339, r.DateBegin)
		r.DateBegin = utils.ToDate(date)
		date, _ = time.Parse(time.RFC3339, r.DateChange)
		r.DateChange = utils.ToDate(date)
		date, _ = time.Parse(time.RFC3339, r.DateEnd)
		r.DateEnd = utils.ToDate(date)
		data = append(data, r)
	}

	return &data, nil
}

func (m *patientModel) NewInvalid(newInvalid *types.NewInvalid, tx *sql.Tx) (sql.Result, error) {
	sqlQuery := fmt.Sprintf(`insert into invalid(
Patient_id, INV_BEGIN, INV_END,
INV_REAS, INV_KIND, INS_WHO, INS_DATE)
values(%v, '%s', '%s',
'%s', '%s', %v, '%s')`,
		newInvalid.PatientId, newInvalid.DateStart, newInvalid.DateEnd,
		newInvalid.Reason, newInvalid.Kind, newInvalid.DoctId, time.Now().Format("2006-01-02"))
	return tx.Exec(sqlQuery)
}

func (m *patientModel) NewChildInvalid(newInvalid *types.NewInvalid, tx *sql.Tx) (sql.Result, error) {
	sqlQuery := fmt.Sprintf(`insert into child_inv(patient_id, GLAVN_NAR, VED_OGR)
values (%v, '%s', '%s')`,
		newInvalid.PatientId, newInvalid.Anomaly, newInvalid.Limit)
	return tx.Exec(sqlQuery)
}

func (m *patientModel) UpdInvalid(newInvalid *types.NewInvalid, tx *sql.Tx) (sql.Result, error) {
	sqlQuery := fmt.Sprintf(`update invalid
set inv_change = '%s'
where patient_id = %v
and inv_end = (select max(inv_end) from invalid where patient_id = %v)`,
		newInvalid.DateDocument, newInvalid.PatientId, newInvalid.PatientId)
	return tx.Exec(sqlQuery)
}

func (m *patientModel) FindCustody(id int64, tx *sql.Tx) (*[]types.FindCustody, error) {
	sqlQuery := fmt.Sprintf(`select KTO, DB, DE from  find_opeka(%v)
order by db desc`, id)

	rows, err := tx.Query(sqlQuery)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	var data = make([]types.FindCustody, 0)
	for rows.Next() {
		r := types.FindCustody{}
		var dateStart sql.NullString
		var dateEnd sql.NullString
		rows.Scan(&r.Who, &dateStart, &dateEnd)

		r.DateStart = dateStart.String
		if r.DateStart != "" {
			date, _ := time.Parse(time.RFC3339, r.DateStart)
			r.DateStart = utils.ToDate(date)
		}

		r.DateEnd = dateEnd.String
		if r.DateEnd != "" {
			date, _ := time.Parse(time.RFC3339, r.DateEnd)
			r.DateEnd = utils.ToDate(date)
		}

		r.Who, _ = utils.ToUTF8(strings.Trim(r.Who, " "))
		data = append(data, r)
	}

	return &data, nil
}

func (m *patientModel) NewCustody(custody *types.NewCustody, tx *sql.Tx) (sql.Result, error) {
	sqlQuery := fmt.Sprintf(`insert into opekun(patient_id, op_begin, who, ins_who, ins_date)
values(%v, '%s', '%s', %v, '%s')`,
		custody.PatientId, custody.DateStart, custody.Custody, custody.DoctId, time.Now().Format("2006-01-02"))
	return tx.Exec(sqlQuery)
}

func (m *patientModel) UpdCustody(custody *types.NewCustody, tx *sql.Tx) (sql.Result, error) {
	sqlQuery := fmt.Sprintf(`update opekun set op_end = '%s' where patient_id = %v and op_begin = '%s'`,
		custody.DateEnd, custody.PatientId, custody.DateStart)
	return tx.Exec(sqlQuery)
}

func (m *patientModel) FindVaccination(id int64, tx *sql.Tx) (*[]types.FindVaccination, error) {
	sqlQuery := fmt.Sprintf(`SELECT dat, priv, nomer, seria, resul , med_otvod from find_priv(%v)
order by dat desc`, id)

	rows, err := tx.Query(sqlQuery)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	var data = make([]types.FindVaccination, 0)
	for rows.Next() {
		r := types.FindVaccination{}
		rows.Scan(&r.Date, &r.Vaccination, &r.Number, &r.Series, &r.Result, &r.Detached)
		r.Date, _ = utils.FormatToDate(r.Date)
		r.Vaccination, _ = utils.ToUTF8(r.Vaccination)
		r.Number, _ = utils.ToUTF8(r.Number)
		r.Series, _ = utils.ToUTF8(r.Series)
		r.Result, _ = utils.ToUTF8(r.Result)
		r.Detached, _ = utils.ToUTF8(r.Detached)

		data = append(data, r)
	}

	return &data, nil
}

func (m *patientModel) FindInfection(id int64, tx *sql.Tx) (*[]types.FindInfection, error) {
	sqlQuery := fmt.Sprintf(`select datp, diag from find_infec(%v)`, id)

	rows, err := tx.Query(sqlQuery)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	var data = make([]types.FindInfection, 0)
	for rows.Next() {
		r := types.FindInfection{}
		rows.Scan(&r.Date, &r.Diagnose)
		r.Date, _ = utils.FormatToDate(r.Date)
		r.Diagnose, _ = utils.ToUTF8(r.Diagnose)

		data = append(data, r)
	}

	return &data, nil
}

func (m *patientModel) UpdPassport(passport *types.Patient, tx *sql.Tx) (sql.Result, error) {
	sqlQuery := fmt.Sprintf(`update general
set passp_ser = '%s', passp_num = %v, job = '%s', residence = %v where patient_id = %v`,
		passport.PassportSeries, passport.PassportNumber, passport.Snils, passport.Works, passport.Id)
	INFO.Println(sqlQuery)
	return tx.Exec(sqlQuery)
}

func (m *patientModel) UpdAddress(address *types.Patient, tx *sql.Tx) (sql.Result, error) {
	sqlQuery := fmt.Sprintf(`update general
set republic = ?, region = ?, district = ?, pop_area = ?, street = ?, house = ?, building = ?, flat = ?, domicile = ? where patient_id = ?`)
	INFO.Println(sqlQuery)
	stmt, err := tx.Prepare(sqlQuery)
	defer stmt.Close()
	if err != nil {
		return nil, err
	}
	return stmt.Exec(address.Republic, address.Region, address.District, address.Area, address.Street, address.House, address.Build, address.Flat, address.Domicile, address.Id)
}

func (m *patientModel) GetSection22(id int64, tx *sql.Tx) (*[]types.ST22, error) {
	sqlQuery := fmt.Sprintf(`SELECT NOM_Z, PATIENT_ID, DAT_BEG, DAT_END, ST, CHAST, INS_WHO, INS_DAT FROM st22 WHERE PATIENT_ID = %v`, id)

	rows, err := tx.Query(sqlQuery)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	var data = make([]types.ST22, 0)
	for rows.Next() {
		r := types.ST22{}
		rows.Scan(&r.Id, &r.PatientId, &r.DateStart, &r.DateEnd, &r.Section, &r.Part, &r.InsWho, &r.InsDate)
		r.DateStart, _ = utils.FormatToDate(r.DateStart)
		r.DateEnd, _ = utils.FormatToDate(r.DateEnd)
		r.InsDate, _ = utils.FormatToDate(r.InsDate)

		data = append(data, r)
	}

	return &data, nil
}

func (m *patientModel) NewSection22(section *types.ST22, tx *sql.Tx) (sql.Result, error) {
	sqlQuery := fmt.Sprintf(`insert into st22(PATIENT_ID, DAT_BEG, DAT_END, ST, CHAST, INS_WHO)
values (%v, '%s', '%s', %v, %v, %v)`,
		section.PatientId, section.DateStart, section.DateEnd, section.Section, section.Part, section.InsWho)
	INFO.Println(sqlQuery)
	return tx.Exec(sqlQuery)
}

func (m *patientModel) SOD(id int64, tx *sql.Tx) (*[]types.SOD, error) {
	sqlQuery := fmt.Sprintf(`select sod_date,statia,chast    from  SOD  where  patient_id = %v
order  by   sod_date desc`, id)

	rows, err := tx.Query(sqlQuery)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	var data = make([]types.SOD, 0)
	for rows.Next() {
		r := types.SOD{}
		rows.Scan(&r.Date, &r.Section, &r.Part)

		data = append(data, r)
	}

	return &data, nil
}

func (m *patientModel) OODLast(id int64, tx *sql.Tx) (*types.OOD, error) {
	sqlQuery := fmt.Sprintf(`select maska1, maska2, maska3, maska4 from profil_ood p where p.patient_id = %v
and nom_z = ( select max(nom_z) from profil_ood p1 where p1.patient_id = p.patient_id)`, id)

	rows, err := tx.Query(sqlQuery)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	var r = types.OOD{}
	for rows.Next() {
		rows.Scan(&r.Danger, &r.Syndrome, &r.Difficulties, &r.Attitude)
	}

	return &r, nil
}

func (m *patientModel) FindSection29(id int64, tx *sql.Tx) (*[]types.FindSection29, error) {
	sqlQuery := fmt.Sprintf(`select f.dat_p, f.diag_p, f.dat_e, trim(f.rezult_p)||rtrim(f.rezult_k)  from  find_post_29(%v)   f
where  ((bin_and( f.priem_post_reason,32) = 32))
order by  dat_p DESC`, id)

	rows, err := tx.Query(sqlQuery)
	defer rows.Close()
	data := make([]types.FindSection29, 0)
	if err != nil {
		return nil, err
	}
	var r = types.FindSection29{}
	for rows.Next() {
		rows.Scan(&r.DateStart, &r.Diagnose, &r.DateEnd, &r.Section)
		r.Section, _ = utils.ToUTF8(r.Section)
		data = append(data, r)
	}

	return &data, nil
}

func (m *patientModel) NewOOD(ood *types.OOD, tx *sql.Tx) (sql.Result, error) {
	sqlQuery := fmt.Sprintf(` insert into profil_ood( patient_id,vid,maska1,maska2,maska3,maska4,
    beg_date,end_date ,upd_who, upd_date )
 values( %v, %v, %v, %v, %v, %v,
          '%s', '%s' , %v, '%s' )`,
		ood.PatientId, 0, ood.Danger, ood.Syndrome, ood.Difficulties, ood.Attitude,
		"01.01.1900", "01.01.2222", ood.UserId, time.Now().Format("2006-01-02"))
	INFO.Println(sqlQuery)
	return tx.Exec(sqlQuery)
}

func (m *patientModel) NewSOD(sod *types.SOD, tx *sql.Tx) (sql.Result, error) {
	sqlQuery := fmt.Sprintf(`insert into SOD (patient_id,sod_date,statia, chast,prinud)
values (%v, '%s', %v, %v, %v)`, sod.PatientId, sod.Date, sod.Section, sod.Part, 0)
	INFO.Println(sqlQuery)
	return tx.Exec(sqlQuery)
}

func (m *patientModel) GetDoctorsVisitByPatient(id int, date time.Time, tx *sql.Tx) (*[]types.Doctor, error) {
	sqlQuery := fmt.Sprintf(`select distinct b.KOD_DOCK_I , RTRIM(b.fio), rtrim(b.im), rtrim(b.ot)
 from visit a, spr_doct b
where a.patient_id = ?
      and b.kod_dock_i = RTRIM(a.name_doct)
      and a.v_date > ?
order by 2`)

	rows, err := tx.Query(sqlQuery, id, date)
	defer rows.Close()
	data := make([]types.Doctor, 0)
	if err != nil {
		return nil, err
	}
	var r = types.Doctor{}
	for rows.Next() {
		rows.Scan(&r.Id, &r.Lname, &r.Fname, &r.Sname)
		r.Lname, _ = utils.ToUTF8(r.Lname)
		r.Fname, _ = utils.ToUTF8(r.Fname)
		r.Sname, _ = utils.ToUTF8(r.Sname)
		data = append(data, r)
	}

	return &data, nil
}

func (m *patientModel) GetLastUKLByVisitPatient(id int, tx *sql.Tx) (*types.UKLData, error) {
	sqlQuery := fmt.Sprintf(`select nom_z,
p1_1, p1_2, p1_3, p1_4, p1_5, p1_6, p1_7, p1_8, p1_9, p1_10, p1_11, p1_12, p1_13, p1_14, p1_15, p1_16, p1_17, p1_18, p1_19, p1_20, p1_21, p1_22, p1_23, p1_24, p1_25, p1_26, p1_27, p1_28, p1_29, p1_30, p1_31, p1_32, p1_33, p1_34, p1_35,
p2_1, p2_2, p2_3, p2_4, p2_5, p2_6, p2_7, p2_8, p2_9, p2_10, p2_11, p2_12, p2_13, p2_14, p2_15, p2_16, p2_17, p2_18, p2_19, p2_20, p2_21, p2_22, p2_23, p2_24, p2_25, p2_26, p2_27, p2_28, p2_29, p2_30, p2_31, p2_32, p2_33, p2_34, p2_35,
p3_1, p3_2, p3_3, p3_4, p3_5, p3_6, p3_7, p3_8, p3_9, p3_10, p3_11, p3_12, p3_13, p3_14, p3_15, p3_16, p3_17, p3_18, p3_19, p3_20, p3_21, p3_22, p3_23, p3_24, p3_25, p3_26, p3_27, p3_28, p3_29, p3_30, p3_31, p3_32, p3_33, p3_34, p3_35,
NZ_REGISTRAT, p1_user, p2_user, p3_user, p1_date, p2_date, p3_date, dock
from ukl u
where patient_id = ? and nom_z = (select max(nom_z) from ukl
where patient_id = u.PATIENT_ID  and nz_Registrat <> 0)`)

	row := tx.QueryRow(sqlQuery, id)
	data := types.UKLData{}
	row.Scan(
		&data.Id,
		&data.P1_1, &data.P1_2, &data.P1_3, &data.P1_4, &data.P1_5, &data.P1_6, &data.P1_7, &data.P1_8, &data.P1_9, &data.P1_10, &data.P1_11, &data.P1_12, &data.P1_13, &data.P1_14, &data.P1_15, &data.P1_16, &data.P1_17, &data.P1_18, &data.P1_19, &data.P1_20, &data.P1_21, &data.P1_22, &data.P1_23, &data.P1_24, &data.P1_25, &data.P1_26, &data.P1_27, &data.P1_28, &data.P1_29, &data.P1_30, &data.P1_31, &data.P1_32, &data.P1_33, &data.P1_34, &data.P1_35,
		&data.P2_1, &data.P2_2, &data.P2_3, &data.P2_4, &data.P2_5, &data.P2_6, &data.P2_7, &data.P2_8, &data.P2_9, &data.P2_10, &data.P2_11, &data.P2_12, &data.P2_13, &data.P2_14, &data.P2_15, &data.P2_16, &data.P2_17, &data.P2_18, &data.P2_19, &data.P2_20, &data.P2_21, &data.P2_22, &data.P2_23, &data.P2_24, &data.P2_25, &data.P2_26, &data.P2_27, &data.P2_28, &data.P2_29, &data.P2_30, &data.P2_31, &data.P2_32, &data.P2_33, &data.P2_34, &data.P2_35,
		&data.P3_1, &data.P3_2, &data.P3_3, &data.P3_4, &data.P3_5, &data.P3_6, &data.P3_7, &data.P3_8, &data.P3_9, &data.P3_10, &data.P3_11, &data.P3_12, &data.P3_13, &data.P3_14, &data.P3_15, &data.P3_16, &data.P3_17, &data.P3_18, &data.P3_19, &data.P3_20, &data.P3_21, &data.P3_22, &data.P3_23, &data.P3_24, &data.P3_25, &data.P3_26, &data.P3_27, &data.P3_28, &data.P3_29, &data.P3_30, &data.P3_31, &data.P3_32, &data.P3_33, &data.P3_34, &data.P3_35,
		&data.RegistratId, &data.User1, &data.User2, &data.User3, &data.Date1, &data.Date2, &data.Date3, &data.Doctor,
	)

	return &data, nil
}

func (m *patientModel) NewUKLByVisitPatientLvl1(data *types.NewUKL, regId int, tx *sql.Tx) (sql.Result, error) {
	sqlQuery := fmt.Sprintf(`insert into ukl( PATIENT_ID,
    P1_1, P1_2, P1_3,
    P1_4, P1_5, P1_6,
    P1_7, P1_8, P1_9,
    P1_10, P1_11, P1_12,
    P1_13, P1_14, P1_15,
    P1_16, 
    P1_USER, P1_DATE, NZ_REGISTRAT, dock, ins_dat)
values( %v, 
    %v, %v, %v,
    %v, %v, %v,
    %v, %v, %v,
    %v, %v, %v,
    %v, %v, %v,
    %v,
    %v, '%s', %v, %v, '%s')`,
		data.PatientId,
		data.Evaluations[0], data.Evaluations[1], data.Evaluations[2],
		data.Evaluations[3], data.Evaluations[4], data.Evaluations[5],
		data.Evaluations[6], data.Evaluations[7], data.Evaluations[8],
		data.Evaluations[9], data.Evaluations[10], data.Evaluations[11],
		data.Evaluations[12], data.Evaluations[13], data.Evaluations[14],
		data.Evaluations[15],
		data.UserId, data.Date, regId, data.DoctorId, time.Now().Format(consts.DATE_FORMAT_DB),
	)
	INFO.Println(sqlQuery)
	return tx.Exec(sqlQuery)
}

func (m *patientModel) NewUKLByVisitPatientLvl2(data *types.NewUKL, tx *sql.Tx) (sql.Result, error) {
	sqlQuery := fmt.Sprintf(`update ukl
set P2_1 = %v,
       P2_2 = %v, 
       P2_3 = %v,
       P2_4 = %v, 
       P2_5 = %v, 
       P2_6 = %v,
       P2_7 = %v, 
       P2_8 = %v, 
       P2_9 = %v,
       P2_10 = %v, 
       P2_11 = %v, 
       P2_12 = %v,
       P2_13 = %v, 
       P2_14 = %v, 
       P2_15 = %v,
       P2_16 = %v, 
       P2_USER = %v, 
       P2_DATE = '%s'
where nom_z = %v`,
		data.Evaluations[0], data.Evaluations[1], data.Evaluations[2],
		data.Evaluations[3], data.Evaluations[4], data.Evaluations[5],
		data.Evaluations[6], data.Evaluations[7], data.Evaluations[8],
		data.Evaluations[9], data.Evaluations[10], data.Evaluations[11],
		data.Evaluations[12], data.Evaluations[13], data.Evaluations[14],
		data.Evaluations[15],
		data.UserId, data.Date, data.Id,
	)
	INFO.Println(sqlQuery)
	return tx.Exec(sqlQuery)
}

func (m *patientModel) NewUKLByVisitPatientLvl3(data *types.NewUKL, tx *sql.Tx) (sql.Result, error) {
	sqlQuery := fmt.Sprintf(`update ukl
set P3_1 = %v,
       P3_2 = %v, 
       P3_3 = %v,
       P3_4 = %v, 
       P3_5 = %v, 
       P3_6 = %v,
       P3_7 = %v, 
       P3_8 = %v, 
       P3_9 = %v,
       P3_10 = %v, 
       P3_11 = %v, 
       P3_12 = %v,
       P3_13 = %v, 
       P3_14 = %v, 
       P3_15 = %v,
       P3_16 = %v, 
       P3_USER = %v, 
       P3_DATE = '%s'
where nom_z = %v`,
		data.Evaluations[0], data.Evaluations[1], data.Evaluations[2],
		data.Evaluations[3], data.Evaluations[4], data.Evaluations[5],
		data.Evaluations[6], data.Evaluations[7], data.Evaluations[8],
		data.Evaluations[9], data.Evaluations[10], data.Evaluations[11],
		data.Evaluations[12], data.Evaluations[13], data.Evaluations[14],
		data.Evaluations[15],
		data.UserId, data.Date, data.Id,
	)
	INFO.Println(sqlQuery)
	return tx.Exec(sqlQuery)
}

func (m *patientModel) CheckUKLLastVisit(patientId int, unit int, tx *sql.Tx) (*types.Visit, error) {
	sqlQuery := `select V_NUM, V_DATE, NAME_DOCT, DIAGNOSE from visit
where patient_id = ?
and bin_and(maska2, ?) = ?
and v_num = (select max(v_num) from visit where patient_id = ? and bin_and(maska2, ?) = ?)`
	row := tx.QueryRow(sqlQuery, patientId, unit, unit, patientId, unit, unit)
	data := types.Visit{}
	err := row.Scan(&data.Id, &data.Date, &data.DockName, &data.Diag)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &data, nil
}

func (m *patientModel) GetLastUKLBySuicide(id int, tx *sql.Tx) (*types.UKLData, error) {
	sqlQuery := fmt.Sprintf(`select nom_z,
p1_1, p1_2, p1_3, p1_4, p1_5, p1_6, p1_7, p1_8, p1_9, p1_10, p1_11, p1_12, p1_13, p1_14, p1_15, p1_16, p1_17, p1_18, p1_19, p1_20, p1_21, p1_22, p1_23, p1_24, p1_25, p1_26, p1_27, p1_28, p1_29, p1_30, p1_31, p1_32, p1_33, p1_34, p1_35,
p2_1, p2_2, p2_3, p2_4, p2_5, p2_6, p2_7, p2_8, p2_9, p2_10, p2_11, p2_12, p2_13, p2_14, p2_15, p2_16, p2_17, p2_18, p2_19, p2_20, p2_21, p2_22, p2_23, p2_24, p2_25, p2_26, p2_27, p2_28, p2_29, p2_30, p2_31, p2_32, p2_33, p2_34, p2_35,
p3_1, p3_2, p3_3, p3_4, p3_5, p3_6, p3_7, p3_8, p3_9, p3_10, p3_11, p3_12, p3_13, p3_14, p3_15, p3_16, p3_17, p3_18, p3_19, p3_20, p3_21, p3_22, p3_23, p3_24, p3_25, p3_26, p3_27, p3_28, p3_29, p3_30, p3_31, p3_32, p3_33, p3_34, p3_35,
NZ_REGISTRAT, p1_user, p2_user, p3_user, p1_date, p2_date, p3_date, dock, nz_visit
from ukl u
where patient_id = ? and nom_z = (select max(nom_z) from ukl
where patient_id = u.PATIENT_ID  and nz_visit > 0)`)

	row := tx.QueryRow(sqlQuery, id)
	data := types.UKLData{}
	row.Scan(
		&data.Id,
		&data.P1_1, &data.P1_2, &data.P1_3, &data.P1_4, &data.P1_5, &data.P1_6, &data.P1_7, &data.P1_8, &data.P1_9, &data.P1_10, &data.P1_11, &data.P1_12, &data.P1_13, &data.P1_14, &data.P1_15, &data.P1_16, &data.P1_17, &data.P1_18, &data.P1_19, &data.P1_20, &data.P1_21, &data.P1_22, &data.P1_23, &data.P1_24, &data.P1_25, &data.P1_26, &data.P1_27, &data.P1_28, &data.P1_29, &data.P1_30, &data.P1_31, &data.P1_32, &data.P1_33, &data.P1_34, &data.P1_35,
		&data.P2_1, &data.P2_2, &data.P2_3, &data.P2_4, &data.P2_5, &data.P2_6, &data.P2_7, &data.P2_8, &data.P2_9, &data.P2_10, &data.P2_11, &data.P2_12, &data.P2_13, &data.P2_14, &data.P2_15, &data.P2_16, &data.P2_17, &data.P2_18, &data.P2_19, &data.P2_20, &data.P2_21, &data.P2_22, &data.P2_23, &data.P2_24, &data.P2_25, &data.P2_26, &data.P2_27, &data.P2_28, &data.P2_29, &data.P2_30, &data.P2_31, &data.P2_32, &data.P2_33, &data.P2_34, &data.P2_35,
		&data.P3_1, &data.P3_2, &data.P3_3, &data.P3_4, &data.P3_5, &data.P3_6, &data.P3_7, &data.P3_8, &data.P3_9, &data.P3_10, &data.P3_11, &data.P3_12, &data.P3_13, &data.P3_14, &data.P3_15, &data.P3_16, &data.P3_17, &data.P3_18, &data.P3_19, &data.P3_20, &data.P3_21, &data.P3_22, &data.P3_23, &data.P3_24, &data.P3_25, &data.P3_26, &data.P3_27, &data.P3_28, &data.P3_29, &data.P3_30, &data.P3_31, &data.P3_32, &data.P3_33, &data.P3_34, &data.P3_35,
		&data.RegistratId, &data.User1, &data.User2, &data.User3, &data.Date1, &data.Date2, &data.Date3, &data.Doctor, &data.VisitId,
	)

	return &data, nil
}

func (m *patientModel) NewUKLBySuicide1(data *types.NewUKL, visitId int, tx *sql.Tx) (sql.Result, error) {
	sqlQuery := fmt.Sprintf(`insert into ukl( PATIENT_ID,
    P1_1, P1_2, P1_3,
    P1_4, P1_5, P1_6,
    P1_7, P1_8, P1_9,
    P1_10, P1_11, P1_12,
    P1_USER, P1_DATE, NZ_VISIT, dock, ins_dat)
values( %v, 
    %v, %v, %v,
    %v, %v, %v,
    %v, %v, %v,
    %v, %v, %v,
    %v, '%s', %v, %v, '%s')`,
		data.PatientId,
		data.Evaluations[0], data.Evaluations[1], data.Evaluations[2],
		data.Evaluations[3], data.Evaluations[4], data.Evaluations[5],
		data.Evaluations[6], data.Evaluations[7], data.Evaluations[8],
		data.Evaluations[9], data.Evaluations[10], data.Evaluations[11],
		data.UserId, data.Date, visitId, data.DoctorId, time.Now().Format(consts.DATE_FORMAT_DB),
	)
	INFO.Println(sqlQuery)
	return tx.Exec(sqlQuery)
}

func (m *patientModel) NewUKLBySuicide2(data *types.NewUKL, tx *sql.Tx) (sql.Result, error) {
	sqlQuery := fmt.Sprintf(`update ukl
set P2_1 = %v,
       P2_2 = %v, 
       P2_3 = %v,
       P2_4 = %v, 
       P2_5 = %v, 
       P2_6 = %v,
       P2_7 = %v, 
       P2_8 = %v, 
       P2_9 = %v,
       P2_10 = %v, 
       P2_11 = %v, 
       P2_12 = %v,
       P2_USER = %v, 
       P2_DATE = '%s'
where nom_z = %v`,
		data.Evaluations[0], data.Evaluations[1], data.Evaluations[2],
		data.Evaluations[3], data.Evaluations[4], data.Evaluations[5],
		data.Evaluations[6], data.Evaluations[7], data.Evaluations[8],
		data.Evaluations[9], data.Evaluations[10], data.Evaluations[11],
		data.UserId, data.Date, data.Id,
	)
	INFO.Println(sqlQuery)
	return tx.Exec(sqlQuery)
}

func (m *patientModel) NewUKLBySuicide3(data *types.NewUKL, tx *sql.Tx) (sql.Result, error) {
	sqlQuery := fmt.Sprintf(`update ukl
set P3_1 = %v,
       P3_2 = %v, 
       P3_3 = %v,
       P3_4 = %v, 
       P3_5 = %v, 
       P3_6 = %v,
       P3_7 = %v, 
       P3_8 = %v, 
       P3_9 = %v,
       P3_10 = %v, 
       P3_11 = %v, 
       P3_12 = %v,
       P3_USER = %v, 
       P3_DATE = '%s'
where nom_z = %v`,
		data.Evaluations[0], data.Evaluations[1], data.Evaluations[2],
		data.Evaluations[3], data.Evaluations[4], data.Evaluations[5],
		data.Evaluations[6], data.Evaluations[7], data.Evaluations[8],
		data.Evaluations[9], data.Evaluations[10], data.Evaluations[11],
		data.UserId, data.Date, data.Id,
	)
	INFO.Println(sqlQuery)
	return tx.Exec(sqlQuery)
}

func (m *patientModel) GetLastUKLByPsychotherapy(id int, tx *sql.Tx) (*types.UKLData, error) {
	sqlQuery := fmt.Sprintf(`select nom_z,
p1_1, p1_2, p1_3, p1_4, p1_5, p1_6, p1_7, p1_8, p1_9, p1_10, p1_11, p1_12, p1_13, p1_14, p1_15, p1_16, p1_17, p1_18, p1_19, p1_20, p1_21, p1_22, p1_23, p1_24, p1_25, p1_26, p1_27, p1_28, p1_29, p1_30, p1_31, p1_32, p1_33, p1_34, p1_35,
p2_1, p2_2, p2_3, p2_4, p2_5, p2_6, p2_7, p2_8, p2_9, p2_10, p2_11, p2_12, p2_13, p2_14, p2_15, p2_16, p2_17, p2_18, p2_19, p2_20, p2_21, p2_22, p2_23, p2_24, p2_25, p2_26, p2_27, p2_28, p2_29, p2_30, p2_31, p2_32, p2_33, p2_34, p2_35,
p3_1, p3_2, p3_3, p3_4, p3_5, p3_6, p3_7, p3_8, p3_9, p3_10, p3_11, p3_12, p3_13, p3_14, p3_15, p3_16, p3_17, p3_18, p3_19, p3_20, p3_21, p3_22, p3_23, p3_24, p3_25, p3_26, p3_27, p3_28, p3_29, p3_30, p3_31, p3_32, p3_33, p3_34, p3_35,
NZ_REGISTRAT, p1_user, p2_user, p3_user, p1_date, p2_date, p3_date, dock, nz_visit
from ukl u
where patient_id = ? and nom_z = (select max(nom_z) from ukl
where patient_id = u.PATIENT_ID  and nz_visit > 0)`)

	row := tx.QueryRow(sqlQuery, id)
	data := types.UKLData{}
	row.Scan(
		&data.Id,
		&data.P1_1, &data.P1_2, &data.P1_3, &data.P1_4, &data.P1_5, &data.P1_6, &data.P1_7, &data.P1_8, &data.P1_9, &data.P1_10, &data.P1_11, &data.P1_12, &data.P1_13, &data.P1_14, &data.P1_15, &data.P1_16, &data.P1_17, &data.P1_18, &data.P1_19, &data.P1_20, &data.P1_21, &data.P1_22, &data.P1_23, &data.P1_24, &data.P1_25, &data.P1_26, &data.P1_27, &data.P1_28, &data.P1_29, &data.P1_30, &data.P1_31, &data.P1_32, &data.P1_33, &data.P1_34, &data.P1_35,
		&data.P2_1, &data.P2_2, &data.P2_3, &data.P2_4, &data.P2_5, &data.P2_6, &data.P2_7, &data.P2_8, &data.P2_9, &data.P2_10, &data.P2_11, &data.P2_12, &data.P2_13, &data.P2_14, &data.P2_15, &data.P2_16, &data.P2_17, &data.P2_18, &data.P2_19, &data.P2_20, &data.P2_21, &data.P2_22, &data.P2_23, &data.P2_24, &data.P2_25, &data.P2_26, &data.P2_27, &data.P2_28, &data.P2_29, &data.P2_30, &data.P2_31, &data.P2_32, &data.P2_33, &data.P2_34, &data.P2_35,
		&data.P3_1, &data.P3_2, &data.P3_3, &data.P3_4, &data.P3_5, &data.P3_6, &data.P3_7, &data.P3_8, &data.P3_9, &data.P3_10, &data.P3_11, &data.P3_12, &data.P3_13, &data.P3_14, &data.P3_15, &data.P3_16, &data.P3_17, &data.P3_18, &data.P3_19, &data.P3_20, &data.P3_21, &data.P3_22, &data.P3_23, &data.P3_24, &data.P3_25, &data.P3_26, &data.P3_27, &data.P3_28, &data.P3_29, &data.P3_30, &data.P3_31, &data.P3_32, &data.P3_33, &data.P3_34, &data.P3_35,
		&data.RegistratId, &data.User1, &data.User2, &data.User3, &data.Date1, &data.Date2, &data.Date3, &data.Doctor, &data.VisitId,
	)

	return &data, nil
}

func (m *patientModel) NewUKLByPsychotherapy1(data *types.NewUKL, visitId int, tx *sql.Tx) (sql.Result, error) {
	sqlQuery := fmt.Sprintf(`insert into ukl( PATIENT_ID,
    P1_1, P1_2, P1_3,
    P1_4, P1_5, P1_6,
    P1_7, P1_8, P1_9,
    P1_10, P1_11, P1_12,
    P1_13, P1_14, P1_15,
    P1_16, 
    P1_USER, P1_DATE, NZ_VISIT, dock, ins_dat)
values( %v, 
    %v, %v, %v,
    %v, %v, %v,
    %v, %v, %v,
    %v, %v, %v,
	%v, %v, %v,
	%v,
    %v, '%s', %v, %v, '%s')`,
		data.PatientId,
		data.Evaluations[0], data.Evaluations[1], data.Evaluations[2],
		data.Evaluations[3], data.Evaluations[4], data.Evaluations[5],
		data.Evaluations[6], data.Evaluations[7], data.Evaluations[8],
		data.Evaluations[9], data.Evaluations[10], data.Evaluations[11],
		data.Evaluations[12], data.Evaluations[13], data.Evaluations[14],
		data.Evaluations[15],
		data.UserId, data.Date, visitId, data.DoctorId, time.Now().Format(consts.DATE_FORMAT_DB),
	)
	INFO.Println(sqlQuery)
	return tx.Exec(sqlQuery)
}

func (m *patientModel) NewUKLByPsychotherapy2(data *types.NewUKL, tx *sql.Tx) (sql.Result, error) {
	sqlQuery := fmt.Sprintf(`update ukl
set P2_1 = %v,
P2_2 = %v, 
P2_3 = %v,
P2_4 = %v, 
P2_5 = %v, 
P2_6 = %v,
P2_7 = %v, 
P2_8 = %v, 
P2_9 = %v,
P2_10 = %v, 
P2_11 = %v, 
P2_12 = %v,
P2_13 = %v, 
P2_14 = %v,
P2_15 = %v,
P2_16 = %v, 
P2_USER = %v, 
P2_DATE = '%s'
where nom_z = %v`,
		data.Evaluations[0], data.Evaluations[1], data.Evaluations[2],
		data.Evaluations[3], data.Evaluations[4], data.Evaluations[5],
		data.Evaluations[6], data.Evaluations[7], data.Evaluations[8],
		data.Evaluations[9], data.Evaluations[10], data.Evaluations[11],
		data.Evaluations[12], data.Evaluations[13], data.Evaluations[14],
		data.Evaluations[15],
		data.UserId, data.Date, data.Id,
	)
	INFO.Println(sqlQuery)
	return tx.Exec(sqlQuery)
}

func (m *patientModel) NewUKLByPsychotherapy3(data *types.NewUKL, tx *sql.Tx) (sql.Result, error) {
	sqlQuery := fmt.Sprintf(`update ukl
set P3_1 = %v,
P3_2 = %v, 
P3_3 = %v,
P3_4 = %v, 
P3_5 = %v, 
P3_6 = %v,
P3_7 = %v, 
P3_8 = %v, 
P3_9 = %v,
P3_10 = %v, 
P3_11 = %v, 
P3_12 = %v,
P3_13 = %v, 
P3_14 = %v,
P3_15 = %v,
P3_16 = %v, 
P3_USER = %v, 
P3_DATE = '%s'
where nom_z = %v`,
		data.Evaluations[0], data.Evaluations[1], data.Evaluations[2],
		data.Evaluations[3], data.Evaluations[4], data.Evaluations[5],
		data.Evaluations[6], data.Evaluations[7], data.Evaluations[8],
		data.Evaluations[9], data.Evaluations[10], data.Evaluations[11],
		data.Evaluations[12], data.Evaluations[13], data.Evaluations[14],
		data.Evaluations[15],
		data.UserId, data.Date, data.Id,
	)
	INFO.Println(sqlQuery)
	return tx.Exec(sqlQuery)
}

func (m *patientModel) GetListUKLByPatient(id int, isType int, tx *sql.Tx) (*[]types.UKLData, error) {
	sqlQuery := fmt.Sprintf(`select nom_z,
p1_1, p1_2, p1_3, p1_4, p1_5, p1_6, p1_7, p1_8, p1_9, p1_10, p1_11, p1_12, p1_13, p1_14, p1_15, p1_16, p1_17, p1_18, p1_19, p1_20, p1_21, p1_22, p1_23, p1_24, p1_25, p1_26, p1_27, p1_28, p1_29, p1_30, p1_31, p1_32, p1_33, p1_34, p1_35,
p2_1, p2_2, p2_3, p2_4, p2_5, p2_6, p2_7, p2_8, p2_9, p2_10, p2_11, p2_12, p2_13, p2_14, p2_15, p2_16, p2_17, p2_18, p2_19, p2_20, p2_21, p2_22, p2_23, p2_24, p2_25, p2_26, p2_27, p2_28, p2_29, p2_30, p2_31, p2_32, p2_33, p2_34, p2_35,
p3_1, p3_2, p3_3, p3_4, p3_5, p3_6, p3_7, p3_8, p3_9, p3_10, p3_11, p3_12, p3_13, p3_14, p3_15, p3_16, p3_17, p3_18, p3_19, p3_20, p3_21, p3_22, p3_23, p3_24, p3_25, p3_26, p3_27, p3_28, p3_29, p3_30, p3_31, p3_32, p3_33, p3_34, p3_35,
NZ_REGISTRAT, p1_user, p2_user, p3_user, p1_date, p2_date, p3_date, dock, nz_visit
from ukl where patient_id = ?`)
	if isType == consts.TYPE_UKL_VISIT {
		sqlQuery += " and nz_visit > 0"
	} else {
		sqlQuery += " and nz_registrat > 0"
	}
	INFO.Println(sqlQuery)
	rows, err := tx.Query(sqlQuery, id)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	list := make([]types.UKLData, 0)
	for rows.Next() {
		data := types.UKLData{}
		rows.Scan(
			&data.Id,
			&data.P1_1, &data.P1_2, &data.P1_3, &data.P1_4, &data.P1_5, &data.P1_6, &data.P1_7, &data.P1_8, &data.P1_9, &data.P1_10, &data.P1_11, &data.P1_12, &data.P1_13, &data.P1_14, &data.P1_15, &data.P1_16, &data.P1_17, &data.P1_18, &data.P1_19, &data.P1_20, &data.P1_21, &data.P1_22, &data.P1_23, &data.P1_24, &data.P1_25, &data.P1_26, &data.P1_27, &data.P1_28, &data.P1_29, &data.P1_30, &data.P1_31, &data.P1_32, &data.P1_33, &data.P1_34, &data.P1_35,
			&data.P2_1, &data.P2_2, &data.P2_3, &data.P2_4, &data.P2_5, &data.P2_6, &data.P2_7, &data.P2_8, &data.P2_9, &data.P2_10, &data.P2_11, &data.P2_12, &data.P2_13, &data.P2_14, &data.P2_15, &data.P2_16, &data.P2_17, &data.P2_18, &data.P2_19, &data.P2_20, &data.P2_21, &data.P2_22, &data.P2_23, &data.P2_24, &data.P2_25, &data.P2_26, &data.P2_27, &data.P2_28, &data.P2_29, &data.P2_30, &data.P2_31, &data.P2_32, &data.P2_33, &data.P2_34, &data.P2_35,
			&data.P3_1, &data.P3_2, &data.P3_3, &data.P3_4, &data.P3_5, &data.P3_6, &data.P3_7, &data.P3_8, &data.P3_9, &data.P3_10, &data.P3_11, &data.P3_12, &data.P3_13, &data.P3_14, &data.P3_15, &data.P3_16, &data.P3_17, &data.P3_18, &data.P3_19, &data.P3_20, &data.P3_21, &data.P3_22, &data.P3_23, &data.P3_24, &data.P3_25, &data.P3_26, &data.P3_27, &data.P3_28, &data.P3_29, &data.P3_30, &data.P3_31, &data.P3_32, &data.P3_33, &data.P3_34, &data.P3_35,
			&data.RegistratId, &data.User1, &data.User2, &data.User3, &data.Date1, &data.Date2, &data.Date3, &data.Doctor, &data.VisitId,
		)
		list = append(list, data)
	}

	return &list, nil
}

func (m *patientModel) GetForcedByPatient(id int, tx *sql.Tx) (*[]types.ForcedM, error) {
	sqlQuery := fmt.Sprintf(`select nz, date_b, date_e, nabl, meh, stat, num_p from load_prinudka(%v)`, id)

	rows, err := tx.Query(sqlQuery)
	defer rows.Close()
	data := make([]types.ForcedM, 0)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var r = types.ForcedM{}
		dateStart := sql.NullString{}
		dateEnd := sql.NullString{}
		rows.Scan(&r.Id, &dateStart, &dateEnd, &r.Watch, &r.Mechanism, &r.State, &r.Number)
		r.DateEnd = dateEnd.String
		r.DateStart = dateStart.String
		r.Watch, _ = utils.ToUTF8(strings.Trim(r.Watch, " "))
		r.Mechanism, _ = utils.ToUTF8(strings.Trim(r.Mechanism, " "))
		r.State, _ = utils.ToUTF8(strings.Trim(r.State, " "))
		data = append(data, r)
	}

	return &data, nil
}

func (m *patientModel) GetForcedNumberByPatient(patientId int, number int, tx *sql.Tx) (*types.ForcedM, error) {
	sqlQuery := fmt.Sprintf(`select nz, date_b, date_e, nabl, meh, stat, num_p from load_prinudka(%v) where num_p = %v`, patientId, number)

	row := tx.QueryRow(sqlQuery)
	var r = types.ForcedM{}
	dateStart := sql.NullString{}
	dateEnd := sql.NullString{}
	row.Scan(&r.Id, &dateStart, &dateEnd, &r.Watch, &r.Mechanism, &r.State, &r.Number)
	r.DateEnd = dateEnd.String
	r.DateStart = dateStart.String
	r.Watch, _ = utils.ToUTF8(strings.Trim(r.Watch, " "))
	r.Mechanism, _ = utils.ToUTF8(strings.Trim(r.Mechanism, " "))
	r.State, _ = utils.ToUTF8(strings.Trim(r.State, " "))

	return &r, nil
}

func (m *patientModel) GetViewed(id int, number int, tx *sql.Tx) (*[]types.ViewedM, error) {
	sqlQuery := fmt.Sprintf(`select 
nz, osm_date, fio, fio1,
na_me, n_akt, akt_date, na_me1,
op_date, pol_date, vid_zapis, exit_date, sud 
from load_prin_osmotr(%v, %v)
order by osm_date desc`, id, number)

	INFO.Println(sqlQuery)

	rows, err := tx.Query(sqlQuery)
	defer rows.Close()
	data := make([]types.ViewedM, 0)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var r = types.ViewedM{}
		dateEnd := sql.NullString{}
		doctorName1 := sql.NullString{}
		doctorName2 := sql.NullString{}
		conclusion := sql.NullString{}
		courtDate := sql.NullString{}
		courtConclusionDate := sql.NullString{}
		courtName := sql.NullString{}
		err = rows.Scan(&r.Id, &r.ViewDate, &doctorName1, &doctorName2,
			&conclusion, &r.ActNumber, &r.ActDate, &r.View,
			&courtDate, &courtConclusionDate, &r.Type, &dateEnd, &courtName)
		if err != nil {
			ERROR.Println(err)
			return nil, err
		}
		r.DateEnd = dateEnd.String
		r.DoctorName1, _ = utils.ToUTF8(doctorName1.String)
		r.DoctorName2, _ = utils.ToUTF8(doctorName2.String)
		r.Conclusion, _ = utils.ToUTF8(conclusion.String)
		r.CourtDate = courtDate.String
		r.CourtConclusionDate = courtConclusionDate.String
		r.CourtName, _ = utils.ToUTF8(courtName.String)
		r.View, _ = utils.ToUTF8(r.View)
		//r.DateStart = dateStart.String
		//r.Watch, _ = utils.ToUTF8(strings.Trim(r.Watch, " "))
		//r.Mechanism, _ = utils.ToUTF8(strings.Trim(r.Mechanism, " "))
		//r.State, _ = utils.ToUTF8(strings.Trim(r.State, " "))
		data = append(data, r)
	}

	return &data, nil
}

func (m *patientModel) GetForced(id int, tx *sql.Tx) (*types.Forced, error) {
	sqlQuery := fmt.Sprintf(`select nom_z, PATIENT_ID, num_pr, 
op_date, pol_date, sud,
st, T_NABL, P_PRIN,
TRUD, MEH, dock1,
dock2, osm_date, zakl,
n_akt, akt_date, exit_date,
VID_ZAPIS
from prinud_m
where nom_z = %v`, id)
	INFO.Println(sqlQuery)
	row := tx.QueryRow(sqlQuery)
	var r = types.Forced{}
	actDate := sql.NullString{}
	err := row.Scan(&r.Id, &r.PatientId, &r.Number,
		&r.CourtDate, &r.CourtConclusionDate, &r.CourtId,
		&r.TypeCrimeId, &r.ViewId, &r.ForcedP,
		&r.Sick, &r.Mechanism, &r.DoctorId1,
		&r.DoctorId2, &r.DateView, &r.ConclusionId,
		&r.ActNumber, &actDate, &r.DateEnd,
		&r.TypeId,
	)
	if err != nil {
		ERROR.Println(err)
		return nil, err
	}
	r.ActDate = actDate.String

	return &r, nil
}

func (m *patientModel) GetForcedLastByPatient(patientId int, tx *sql.Tx) (*types.Forced, error) {
	sqlQuery := fmt.Sprintf(`select nom_z, PATIENT_ID, num_pr, 
op_date, pol_date, sud,
st, T_NABL, P_PRIN,
TRUD, MEH, dock1,
dock2, osm_date, zakl,
n_akt, akt_date, exit_date,
VID_ZAPIS
from prinud_m
where patient_id = %v
and osm_date = (select max(osm_date) from prinud_m
where patient_id = %v)`, patientId, patientId)
	INFO.Println(sqlQuery)
	row := tx.QueryRow(sqlQuery)
	var r = types.Forced{}
	actDate := sql.NullString{}
	err := row.Scan(&r.Id, &r.PatientId, &r.Number,
		&r.CourtDate, &r.CourtConclusionDate, &r.CourtId,
		&r.TypeCrimeId, &r.ViewId, &r.ForcedP,
		&r.Sick, &r.Mechanism, &r.DoctorId1,
		&r.DoctorId2, &r.DateView, &r.ConclusionId,
		&r.ActNumber, &actDate, &r.DateEnd,
		&r.TypeId,
	)
	if err != nil {
		ERROR.Println(err)
		return nil, err
	}
	r.ActDate = actDate.String

	return &r, nil
}

func (m *patientModel) PostForcedByPatient(forced *types.Forced, tx *sql.Tx) (sql.Result, error) {
	sql := fmt.Sprintf(`insert into prinud_m(
patient_id, num_pr, op_date,
pol_date, sud, st,
t_nabl, P_PRIN, trud,
MEH, DOCK1, DOCK2,
OSM_DATE, ZAKL, N_AKT, 
AKT_DATE, EXIT_DATE, VID_ZAPIS,
ins_who, ins_date, UPD_WHO, UPD_DATE)
values(%v, %v, '%s',
'%s', %v, %v, 
%v, %v, %v,
%v, %v, %v,
'%s', %v, %v,
'%s', '%s', %v,
%v, '%s', %v, '%s')`,
		forced.PatientId, forced.Number, forced.CourtDate,
		forced.CourtConclusionDate, forced.CourtId, forced.TypeCrimeId,
		forced.ViewId, forced.ForcedP, forced.Sick,
		forced.Mechanism, forced.DoctorId1, forced.DoctorId2,
		forced.DateView, forced.ConclusionId, forced.ActNumber,
		forced.ActDate, forced.DateEnd, forced.TypeId,
		forced.UserId, time.Now().Format(consts.DATE_FORMAT_DB), forced.UserId, time.Now().Format(consts.DATE_FORMAT_DB),
	)
	INFO.Println(sql)
	return tx.Exec(sql)
}

func (m *patientModel) UpdForcedByPatient(forced *types.Forced, tx *sql.Tx) (sql.Result, error) {
	sql := fmt.Sprintf(`update prinud_m 
set num_pr = %v,
op_date = '%s',
pol_date = '%s',
sud = %v,
st = %v,
t_nabl = %v,
P_PRIN = %v,
trud = %v,
MEH = %v,
DOCK1 = %v,
DOCK2 = %v,
OSM_DATE = '%s',
ZAKL = %v,
N_AKT = %v, 
AKT_DATE = '%s',
EXIT_DATE = '%s',
VID_ZAPIS = %v,
UPD_WHO = %v,
UPD_DATE = '%s'
where nom_z = %v`,
		forced.Number, forced.CourtDate,
		forced.CourtConclusionDate, forced.CourtId, forced.TypeCrimeId,
		forced.ViewId, forced.ForcedP, forced.Sick,
		forced.Mechanism, forced.DoctorId1, forced.DoctorId2,
		forced.DateView, forced.ConclusionId, forced.ActNumber,
		forced.ActDate, forced.DateEnd, forced.TypeId,
		forced.UserId, time.Now().Format(consts.DATE_FORMAT_DB), forced.Id,
	)
	INFO.Println(sql)
	return tx.Exec(sql)
}

func (m *patientModel) DeleteForcedByViewDate(forced *types.Forced, tx *sql.Tx) (sql.Result, error) {
	sql := fmt.Sprintf(`delete from prinud_m where OSM_DATE = '%s' and patient_id = %v`,
		forced.DateView, forced.PatientId,
	)
	INFO.Println(sql)
	return tx.Exec(sql)
}

func (m *patientModel) GetNumForcedByPatient(patientId int, tx *sql.Tx) (int, error) {
	sql := fmt.Sprintf(`select num from load_num_pr(%v)`, patientId)
	INFO.Println(sql)
	row := tx.QueryRow(sql)
	var num int
	err := row.Scan(&num)
	if err != nil {
		return 0, nil
	}
	return num, nil
}
