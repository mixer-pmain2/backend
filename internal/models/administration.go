package models

import (
	"database/sql"
	"fmt"
	"pmain2/internal/types"
	"time"
)

type administrationModel struct {
}

func (m *administrationModel) DoctorLocation(unit int, date time.Time, data types.DoctorBySection, tx *sql.Tx) (sql.Result, error) {
	sql := fmt.Sprintf(`insert into dock_uch(DAT, DOCK, UCH, PRIZ, PODRAZ, FIO_DOCK) 
values('%s', %v, %v, 1, %v, '')`, date.Format("2006-01-02"), data.DoctorId, data.Section, unit)
	INFO.Println(sql)
	result, err := tx.Exec(sql)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (m *administrationModel) DisableSections(unit int, tx *sql.Tx) (sql.Result, error) {
	sql := fmt.Sprintf(`update dock_uch 
set priz = 0 
where bin_and(podraz, %v) = %v 
and podraz <> 2147483647 and uch > 10 
and uch <> 299 and priz = 1`, unit, unit)
	INFO.Println(sql)
	result, err := tx.Exec(sql)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (m *administrationModel) DeleteSectionsByDate(unit int, date time.Time, tx *sql.Tx) (sql.Result, error) {
	sql := fmt.Sprintf(`delete from dock_uch 
where bin_and(podraz, %v) = %v 
and podraz <> 2147483647 and uch > 10 
and dat >= '%s' and dat < '%s'
and uch <> 299 and priz = 1`, unit, unit, date.Format("2006-01-02"), date.Add(time.Hour*24).Format("2006-01-02"))
	INFO.Println(sql)
	result, err := tx.Exec(sql)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (m *administrationModel) DoctorLeadSection(unit int, d1 time.Time, d2 time.Time, data types.DoctorBySection, tx *sql.Tx) (sql.Result, error) {
	sql := fmt.Sprintf(`insert into dock_uch_s (date_s, date_e, DOCK, UCH, PODRAZ, FIO_DOCK)
values('%s', '%s', %v, %v, %v, '')`, d1.Format("2006-01-02"), d2.Format("2006-01-02"), data.DoctorId, data.Section, unit)
	INFO.Println(sql)
	result, err := tx.Exec(sql)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (m *administrationModel) DeleteLeadSectionsByDate(unit int, d1 time.Time, d2 time.Time, tx *sql.Tx) (sql.Result, error) {
	sql := fmt.Sprintf(`delete from dock_uch_s 
where bin_and(podraz, %v) = %v 
and podraz <> 2147483647 and uch > 10 
and date_s = '%s' and date_e = '%s'
and uch <> 299`, unit, unit, d1.Format("2006-01-02"), d2.Format("2006-01-02"))
	INFO.Println(sql)
	result, err := tx.Exec(sql)
	if err != nil {
		return nil, err
	}

	return result, nil
}
