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
republic, region, district, pop_area, street, house, building, flat, domicile
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
set republic = %v, region = %v, district = %v, pop_area = %v, street = %v, house = '%s', building = '%s', flat = '%s', domicile = %v where patient_id = %v`,
		address.Republic, address.Region, address.District, address.Area, address.Street, address.House, address.Build, address.Flat, address.Domicile, address.Id)
	INFO.Println(sqlQuery)
	return tx.Exec(sqlQuery)
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
