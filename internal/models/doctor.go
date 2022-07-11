package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"pmain2/internal/types"
	"time"
)

type doctorModel struct {
}

func createDoctor() *doctorModel {
	return &doctorModel{}
}

func (m *doctorModel) GetRate(data types.DoctorFindParams, tx *sql.Tx) (*[]types.DoctorRate, error) {
	result := make([]types.DoctorRate, 0)
	sql := fmt.Sprintf(`select distinct podrazd ,trim(kod_dock), stavka,nz 
from dock_stavka a where trim(kod_dock) = %v and mesec = %v and god = %v
 and bin_or(podrazd,%v)=%v`, data.DoctorId, data.Month, data.Year, data.Unit, data.Unit)
	INFO.Println(sql)
	rows, err := tx.Query(sql)
	defer rows.Close()
	defer rows.Close()
	if err != nil {
		ERROR.Println(err)
		return nil, err
	}

	for rows.Next() {
		row := types.DoctorRate{}
		err := rows.Scan(&row.Unit, &row.DoctorId, &row.Rate, &row.Id)
		if err != nil {
			ERROR.Println(err.Error())
			return nil, err
		}
		result = append(result, row)
	}

	return &result, nil
}

func (m *doctorModel) VisitCountPlan(data types.DoctorFindParams, tx *sql.Tx) (*[]types.DoctorVisitCountPlan, error) {
	result := make([]types.DoctorVisitCountPlan, 0)
	startDate := time.Date(data.Year, time.Month(data.Month), 1, 0, 0, 0, 0, time.UTC)
	date := time.Date(data.Year, time.Month(data.Month+1), 0, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(data.Year, time.Month(data.Month), date.Day(), 0, 0, 0, 0, time.UTC)
	sqlQuery := fmt.Sprintf(`SELECT maska2-bin_and(maska2, 1), count(*), 
   (select count_plan from doct_plan(2018, 1, maska2-bin_and(maska2,1),15))
from visit 
where name_doct = %v and v_date between '%s' and '%s'
group by 1`, data.DoctorId, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	INFO.Println(sqlQuery)
	rows, err := tx.Query(sqlQuery)
	defer rows.Close()
	if err != nil {
		ERROR.Println(err)
		return nil, err
	}

	for rows.Next() {
		row := types.DoctorVisitCountPlan{}
		var plan sql.NullFloat64
		var visit sql.NullInt64
		err := rows.Scan(&row.Unit, &visit, &plan)
		row.Plan = plan.Float64
		row.Visit = visit.Int64
		if row.Unit == 0 {
			row.Unit = 1
		}
		if err != nil {
			ERROR.Println(err.Error())
			return nil, err
		}
		result = append(result, row)
	}

	return &result, nil
}

func (m *doctorModel) GetUnits(data types.DoctorFindParams, tx *sql.Tx) (*[]int, error) {
	result := make([]int, 0)
	sqlQuery := fmt.Sprintf(`select podr from dock_prava a,spr_doct b, spr_prava s
where a.kod_doct = b.kod_dock_i and s.kod1 = 1
and b.kod_dock = %v
and s.maska1 = a.podr
and 'now' between a.date_n and a.date_e
and bin_or(a.podr, %v) = %v
GROUP BY 1`, data.DoctorId, data.Unit, data.Unit)
	INFO.Println(sqlQuery)
	rows, err := tx.Query(sqlQuery)
	if err != nil {
		ERROR.Println(err)
		return nil, err
	}

	for rows.Next() {
		row := 0
		err := rows.Scan(&row)
		if err != nil {
			ERROR.Println(err.Error())
			return nil, err
		}
		result = append(result, row)
	}

	res, err := json.Marshal(&result)
	if err != nil {
		ERROR.Println(err.Error())
		return nil, err
	}
	INFO.Println(string(res))

	return &result, nil
}

func (m *doctorModel) UpdRate(data types.DoctorQueryUpdRate, tx *sql.Tx) (sql.Result, error) {
	sqlQuery := fmt.Sprintf(`update dock_stavka set stavka = '%s'
where kod_dock = %v and podrazd = %v and mesec= %v and god= %v`, data.Rate, data.DoctorId, data.Unit, data.Month, data.Year)
	INFO.Println(sqlQuery)
	return tx.Exec(sqlQuery)
}

func (m *doctorModel) AddRate(data types.DoctorQueryUpdRate, tx *sql.Tx) (sql.Result, error) {
	sqlQuery := fmt.Sprintf(`insert into dock_stavka (kod_dock,podrazd,mesec,god,stavka)
values (%v, %v, %v, %v, '%s')`, data.DoctorId, data.Unit, data.Month, data.Year, data.Rate)
	INFO.Println(sqlQuery)
	return tx.Exec(sqlQuery)
}

func (m *doctorModel) DelRate(data types.DoctorQueryUpdRate, tx *sql.Tx) (sql.Result, error) {
	sqlQuery := fmt.Sprintf(`delete from dock_stavka where kod_dock = %v and podrazd = %v and mesec = %v and god = %v`,
		data.DoctorId, data.Unit, data.Month, data.Year)
	INFO.Println(sqlQuery)
	return tx.Exec(sqlQuery)
}
