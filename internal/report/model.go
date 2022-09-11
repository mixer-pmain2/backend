package report

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"pmain2/internal/database"
	"pmain2/pkg/utils"
	"time"
)

var (
	db *sql.DB
)

func initDBConnect() error {
	conn, err := database.Connect()
	if err != nil {
		return err
	}
	db = conn.DB
	return nil
}

func CreateTx() (error, *sql.Tx) {

	tx, err := db.Begin()

	if err != nil {
		return err, nil
	}
	return nil, tx
}

func newJob(p *reportParams, tx *sql.Tx) (sql.Result, error) {
	bF, _ := json.Marshal(p.Filters)
	sqlQuery := fmt.Sprintf(`insert into report_job (user_id, code, filters) values (%v, '%s', '%s')`, p.UserId, p.Code, bF)

	return tx.Exec(sqlQuery)
}

func getJobs(userId int, tx *sql.Tx) (*[]reportParams, error) {
	sqlQuery := fmt.Sprintf(`select id, user_id, code, filters, status, ins_date from report_job where user_id = %v order by ins_date desc`, userId)
	rows, err := tx.Query(sqlQuery)
	if err != nil {
		ERROR.Println(err)
		return nil, err
	}
	data := make([]reportParams, 0)
	for rows.Next() {
		row := reportParams{}
		var filter string
		rows.Scan(&row.Id, &row.UserId, &row.Code, &filter, &row.Status, &row.Date)

		err = json.Unmarshal([]byte(filter), &row.Filters)
		if err != nil {
			ERROR.Println(err)
		}
		data = append(data, row)
	}

	return &data, nil
}

func getNewJobs(tx *sql.Tx) (*[]reportParams, error) {
	sqlQuery := fmt.Sprintf(`select id, user_id, code, filters from report_job where status in ('NEW', 'PROGRESS')`)
	rows, err := tx.Query(sqlQuery)
	defer rows.Close()
	if err != nil {
		ERROR.Println(err)
		return nil, err
	}
	data := make([]reportParams, 0)
	defer rows.Close()
	for rows.Next() {
		row := reportParams{}
		var filter string
		rows.Scan(&row.Id, &row.UserId, &row.Code, &filter)

		err = json.Unmarshal([]byte(filter), &row.Filters)
		if err != nil {
			ERROR.Println(err)
			setStatusByJob(row, statusType.error, tx)
			continue
		}
		data = append(data, row)
	}

	return &data, nil
}

func setStatusByJob(p reportParams, status string, tx *sql.Tx) (sql.Result, error) {
	sqlQuery := fmt.Sprintf(`update report_job set status = '%s' where id = %v`, status, p.Id)
	INFO.Println(sqlQuery)
	return tx.Exec(sqlQuery)
}

func saveReport(p reportParams, buf *bytes.Buffer, tx *sql.Tx) (sql.Result, error) {
	sqlQuery := fmt.Sprintf(`update report_job set status = ?, report = ? where id = ?`)
	INFO.Println(sqlQuery)
	return tx.Exec(sqlQuery, "DONE", buf.String(), p.Id)
}

func deleteOlderReport(day int, tx *sql.Tx) (sql.Result, error) {
	sqlQuery := fmt.Sprintf(`DELETE FROM REPORT_JOB rj WHERE rj.INS_DATE <= dateadd(-%v DAY TO timestamp 'NOW')`, day)
	INFO.Println(sqlQuery)
	return tx.Exec(sqlQuery)
}

func getJob(id int, tx *sql.Tx) (*[]byte, error) {
	sqlQuery := fmt.Sprintf(`select report from report_job where id = %v`, id)
	row := tx.QueryRow(sqlQuery)
	var data []byte
	err := row.Scan(&data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func getVisitById(id int, tx *sql.Tx) (int, error) {
	sqlQuery := fmt.Sprintf(`SELECT maska2 - bin_and(maska2, 1) FROM visit WHERE V_NUM  = %v`, id)
	row := tx.QueryRow(sqlQuery)
	var data int
	err := row.Scan(&data)
	if err != nil {
		return 0, err
	}

	return data, nil
}

type doctorVisitSection struct {
	Id         int
	DoctorName string
	Section    int
}

func getDoctorsVisitingByUnit(d1 string, d2 string, unit int, tx *sql.Tx) (*[]doctorVisitSection, error) {
	sqlQuery := fmt.Sprintf(`select cast(trim(sd.kod_dock) AS integer), '',
  case when trim(du.uch_dock) = '' then 0 else trim(uch_dock) end uch_dock
  from dock_prava dp, spr_doct sd,
  (select v.uch_dock, v.name_doct, v.maska2 from visit v
                where v_date between '%s' and '%s'
                and case when maska2 - bin_and(maska2,1) = 0 then 1
                else maska2 - bin_and(maska2,1) end = %v) du
  where dp.podr = %v
  and du.name_doct = dp.kod_doct
  and dp.prava = 1
  and ('%s' between dp.date_n and dp.date_e
      or '%s' between dp.date_n and dp.date_e
      or dp.date_n between '%s' and '%s'
      or dp.date_e between '%s' and '%s')
  and sd.kod_dock_i = dp.kod_doct
 union distinct  select cast(trim(d.dock) AS integer) kod_dock, '', d.uch uch_dock
  from dock_uch d, spr_doct sd
  where d.podraz = %v
  and sd.kod_dock = d.dock
  and ((d.priz=1 and d.dat <='%s') or ((select per from date_period('%s','%s',d.dat,d.dat_upd))=1 and d.dat_upd is not null))
  and d.dat >= '01.01.2016'`, d1, d2,
		unit,
		unit,
		d1,
		d2,
		d1, d2,
		d1, d2,
		unit,
		d2, d1, d2)
	INFO.Println(sqlQuery)

	rows, err := tx.Query(sqlQuery)
	if err != nil {
		return nil, err
	}
	data := make([]doctorVisitSection, 0)
	for rows.Next() {
		row := doctorVisitSection{}
		rows.Scan(&row.Id, &row.DoctorName, &row.Section)
		row.DoctorName, _ = utils.ToUTF8(row.DoctorName)

		data = append(data, row)
	}

	return &data, nil
}

func form39GenerateData(d1 string, d2 string, tx *sql.Tx) error {
	incSvod := func(m map[int]map[int]int, column int, day int) map[int]map[int]int {
		_, isDay := m[day]
		if !isDay {
			m[day] = make(map[int]int)
		}
		_, isColumn := m[day][column]
		if !isColumn {
			m[day][column] = 0
		}
		m[day][column] += 1
		return m
	}
	units := []int{1, 2048, 16, 16777216, 33554432, 2, 512, 4, 8, 1024}

	// unit -> doctor -> section -> day -> column -> visit count
	var svod = make(map[int]map[int]map[int]map[int]map[int]int)
	for _, unit := range units {
		s, err := getDoctorsVisitingByUnit(d1, d2, unit, tx)
		if err != nil {
			ERROR.Println(err)
			return err
		}
		for _, doctor := range *s {
			_, found := svod[unit]
			if !found {
				svod[unit] = make(map[int]map[int]map[int]map[int]int)
			}
			_, found = svod[unit][doctor.Id]
			if !found {
				svod[unit][doctor.Id] = make(map[int]map[int]map[int]int)
			}
			_, found = svod[unit][doctor.Section]
			if !found {
				svod[unit][doctor.Id][doctor.Section] = make(map[int]map[int]int)
			}
		}
	}

	visits, err := form39SvodVisit(d1, d2, tx)
	if err != nil {
		ERROR.Println(err)
		return err
	}
	i := 0
	var visit form39SvodType
	defer func() {
		if err := recover(); err != nil {
			log.Println("panic occurred:", err)
			ERROR.Println(i, visit)
		}
	}()
	for i, visit = range *visits {
		fmt.Println(i, visit)
		date, err := time.Parse(time.RFC3339, visit.Date)
		if err != nil {
			ERROR.Println(err)
			return err
		}
		svod[visit.Unit][visit.DoctorId][visit.Section] = incSvod(svod[visit.Unit][visit.DoctorId][visit.Section], 1, date.Day())
		svod[visit.Unit][visit.DoctorId][visit.Section] = incSvod(svod[visit.Unit][visit.DoctorId][visit.Section], 18, date.Day())
		if visit.IsHome == 0 {
			svod[visit.Unit][visit.DoctorId][visit.Section] = incSvod(svod[visit.Unit][visit.DoctorId][visit.Section], 2, date.Day())
			if visit.PatientId != 306258 {
				if visit.Domicile == 2 {
					svod[visit.Unit][visit.DoctorId][visit.Section] = incSvod(svod[visit.Unit][visit.DoctorId][visit.Section], 3, date.Day())
				}
			}
			if visit.Old <= 17 {
				svod[visit.Unit][visit.DoctorId][visit.Section] = incSvod(svod[visit.Unit][visit.DoctorId][visit.Section], 4, date.Day())
			}
			if visit.Old >= 60 && visit.Old <= 99 {
				svod[visit.Unit][visit.DoctorId][visit.Section] = incSvod(svod[visit.Unit][visit.DoctorId][visit.Section], 5, date.Day())
			}
			if visit.TypeVisit&1024 == 0 {
				svod[visit.Unit][visit.DoctorId][visit.Section] = incSvod(svod[visit.Unit][visit.DoctorId][visit.Section], 6, date.Day())
				if visit.Old <= 17 {
					svod[visit.Unit][visit.DoctorId][visit.Section] = incSvod(svod[visit.Unit][visit.DoctorId][visit.Section], 7, date.Day())
				}
				if visit.Old >= 60 && visit.Old <= 99 {
					svod[visit.Unit][visit.DoctorId][visit.Section] = incSvod(svod[visit.Unit][visit.DoctorId][visit.Section], 8, date.Day())
				}
			}
			if visit.TypeVisit&1024 == 1024 {
				svod[visit.Unit][visit.DoctorId][visit.Section] = incSvod(svod[visit.Unit][visit.DoctorId][visit.Section], 9, date.Day())
				if visit.PatientId == 306258 {
					svod[visit.Unit][visit.DoctorId][visit.Section] = incSvod(svod[visit.Unit][visit.DoctorId][visit.Section], 21, date.Day())
				}
			}
		}
		if visit.IsHome == 1 {
			svod[visit.Unit][visit.DoctorId][visit.Section] = incSvod(svod[visit.Unit][visit.DoctorId][visit.Section], 10, date.Day())
			if visit.TypeVisit&1024 == 0 {
				svod[visit.Unit][visit.DoctorId][visit.Section] = incSvod(svod[visit.Unit][visit.DoctorId][visit.Section], 11, date.Day())
				if visit.Old <= 17 {
					svod[visit.Unit][visit.DoctorId][visit.Section] = incSvod(svod[visit.Unit][visit.DoctorId][visit.Section], 12, date.Day())
				}
				if visit.Old <= 1 {
					svod[visit.Unit][visit.DoctorId][visit.Section] = incSvod(svod[visit.Unit][visit.DoctorId][visit.Section], 13, date.Day())
				}
				if visit.Old >= 60 && visit.Old <= 99 {
					svod[visit.Unit][visit.DoctorId][visit.Section] = incSvod(svod[visit.Unit][visit.DoctorId][visit.Section], 14, date.Day())
				}
			}
			if visit.TypeVisit&1024 == 1024 {
				svod[visit.Unit][visit.DoctorId][visit.Section] = incSvod(svod[visit.Unit][visit.DoctorId][visit.Section], 9, date.Day())
				if visit.Old <= 17 {
					svod[visit.Unit][visit.DoctorId][visit.Section] = incSvod(svod[visit.Unit][visit.DoctorId][visit.Section], 15, date.Day())
				}
				if visit.Old <= 1 {
					svod[visit.Unit][visit.DoctorId][visit.Section] = incSvod(svod[visit.Unit][visit.DoctorId][visit.Section], 16, date.Day())
				}
				if visit.PatientId == 306258 {
					svod[visit.Unit][visit.DoctorId][visit.Section] = incSvod(svod[visit.Unit][visit.DoctorId][visit.Section], 21, date.Day())
				}
			}
		}
	}

	spravList, err := form39SpravModif(d1, d2, tx)
	if err != nil {
		ERROR.Println(err)
		return err
	}
	for i, row := range *spravList {
		fmt.Println(i, row)
		date, err := time.Parse(time.RFC3339, row.Date)
		if err != nil {
			ERROR.Println(err)
			return err
		}
		section := 0
		for k, _ := range svod[row.Unit][row.DoctorId] {
			section = k
			break
		}
		svod[row.Unit][row.DoctorId][section] = incSvod(svod[row.Unit][row.DoctorId][section], 21, date.Day())
		switch row.Unit {
		case 1, 2048, 2, 4, 8, 1024:
			if row.Cost == 0 {
				svod[row.Unit][row.DoctorId][section] = incSvod(svod[row.Unit][row.DoctorId][section], 1, date.Day())
				svod[row.Unit][row.DoctorId][section] = incSvod(svod[row.Unit][row.DoctorId][section], 2, date.Day())
				svod[row.Unit][row.DoctorId][section] = incSvod(svod[row.Unit][row.DoctorId][section], 9, date.Day())
				svod[row.Unit][row.DoctorId][section] = incSvod(svod[row.Unit][row.DoctorId][section], 18, date.Day())
			}
		}
		switch row.Unit {
		case 1, 2048, 2, 16, 16777216, 33554432, 4, 8, 1024:
			if row.Cost > 0 {
				svod[row.Unit][row.DoctorId][section] = incSvod(svod[row.Unit][row.DoctorId][section], 1, date.Day())
				svod[row.Unit][row.DoctorId][section] = incSvod(svod[row.Unit][row.DoctorId][section], 19, date.Day())
			}
		}
	}
	queryDeleteArch := fmt.Sprintf(`delete from ARHIV_OTCHET where d1 >= '%s' and d2 <= '%s' and n_type = %v`, d1, d2, 1)
	_, err = tx.Exec(queryDeleteArch)
	if err != nil {
		ERROR.Println(err)
		return err
	}

	getColumnt := func(arr map[int]int, column int) int {
		value, isHave := arr[column]
		if !isHave {
			return 0
		}
		return value
	}
	headUnit := map[int]int{
		1:        1,
		2048:     1,
		2:        2,
		16:       16,
		16777216: 16,
		33554432: 16,
		512:      512,
		4:        4,
		8:        8,
		1024:     1,
	}
	for unit, unitValue := range svod {
		for doctor, doctorValue := range unitValue {
			for section, sectionValue := range doctorValue {
				for day, dayValue := range sectionValue {
					fmt.Println(fmt.Sprintf(`unit: %v, doctor: %v, section: %v, day: %v, count: %v %v`, unit, doctor, section, day, dayValue[1], dayValue[1]))
					queryIns := fmt.Sprintf(
						`insert into ARHIV_OTCHET (d1, d2, n_type, n_podr, maska1, kod_Doct, uch_doct, 
p1,p2,p3,p4,p5,p6,p7,p8,p9,p10,p11,p12,p13,p14,p15,p16,p17,p18,p19,p20,p21)
values ('%s', '%s', %v, %v, %v, %v, %v,
%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v)`,
						d1, d2, 1, headUnit[unit], unit, doctor, section,
						getColumnt(dayValue, 1), getColumnt(dayValue, 2), getColumnt(dayValue, 3), getColumnt(dayValue, 4), getColumnt(dayValue, 5), getColumnt(dayValue, 6), getColumnt(dayValue, 7), getColumnt(dayValue, 8), getColumnt(dayValue, 9), getColumnt(dayValue, 10), getColumnt(dayValue, 11), getColumnt(dayValue, 12), getColumnt(dayValue, 13),
						getColumnt(dayValue, 14), getColumnt(dayValue, 15), getColumnt(dayValue, 16), getColumnt(dayValue, 17), getColumnt(dayValue, 18), getColumnt(dayValue, 19), getColumnt(dayValue, 20), getColumnt(dayValue, 21),
					)
					_, err = tx.Exec(queryIns)
					if err != nil {
						INFO.Println(err)
						return err
					}
				}
			}

		}
	}

	return nil
}

type form39SvodType struct {
	Date      string
	Section   int
	DoctorId  int
	Unit      int
	IsHome    int
	PatientId int
	Domicile  int
	Old       int
	TypeVisit int
}

func form39SvodVisit(d1 string, d2 string, tx *sql.Tx) (*[]form39SvodType, error) {
	sqlQuery := fmt.Sprintf(`SELECT EXTRACT(YEAR FROM CAST('%s' AS date)) - EXTRACT(YEAR FROM bday) vozrast,
      (select domicile from general where patient_id = v.patient_id) domicile,
      v.patient_id, v.v_date, cast(trim(v.NAME_DOCT) as integer), maska1, maska2 - bin_and(maska2, 1) unit, bin_and(maska2, 1) home,
      CAST(case when uch_dock='' then 0 else uch_dock END AS integer) uch_doct from visit v where v_date between '%s' and '%s'`, d2, d1, d2)

	INFO.Println(sqlQuery)

	rows, err := tx.Query(sqlQuery)
	if err != nil {
		ERROR.Println(err)
		return nil, err
	}
	data := make([]form39SvodType, 0)
	for rows.Next() {
		row := form39SvodType{}
		err = rows.Scan(&row.Old, &row.Domicile, &row.PatientId, &row.Date, &row.DoctorId, &row.TypeVisit, &row.Unit, &row.IsHome, &row.Section)
		if err != nil {
			ERROR.Println(err)
			return nil, err
		}
		if row.Unit == 0 {
			row.Unit = 1
		}

		data = append(data, row)
	}
	return &data, nil
}

type spravModifType struct {
	Cost     int
	DoctorId int
	Unit     int
	Date     string
}

func form39SpravModif(d1 string, d2 string, tx *sql.Tx) (*[]spravModifType, error) {
	sqlQuery := fmt.Sprintf(`select cost, doct, podr, date_out from F39_SPRAV_MODIF('%s','%s')`, d1, d2)
	INFO.Println(sqlQuery)

	rows, err := tx.Query(sqlQuery)
	if err != nil {
		ERROR.Println(err)
		return nil, err
	}
	data := make([]spravModifType, 0)
	for rows.Next() {
		row := spravModifType{}
		err = rows.Scan(&row.Cost, &row.DoctorId, &row.Unit, &row.Date)
		if err != nil {
			ERROR.Println(err)
			return nil, err
		}
		data = append(data, row)
	}
	return &data, nil
}
