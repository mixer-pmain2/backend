package report

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"pmain2/internal/database"
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
