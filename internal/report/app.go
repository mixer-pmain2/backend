package report

import (
	"bytes"
	"log"
	"time"
)

func Run() error {
	initDBConnect()
	err := process()
	if err != nil {
		return err
	}
	return nil
}

func createReport() error {
	defer func() {
		if err := recover(); err != nil {
			log.Println("panic occurred:", err)
			ERROR.Println("panic occurred:", err)
		}
	}()
	//INFO.Print("Check report job")
	err, tx := CreateTx()
	if err != nil {
		ERROR.Println(err)
		tx.Rollback()
		return err
	}
	data, err := getNewJobs(tx)
	if err != nil {
		ERROR.Println(err)
		tx.Rollback()
		return err
	}
	tx.Commit()
	job := ReportsJob{}
	for _, row := range *data {
		err, tx = CreateTx()
		if err != nil {
			ERROR.Println(err)
			return err
		}
		//отмечам что отчет выполняется
		_, err = setStatusByJob(row, statusType.progress, tx)
		if err != nil {
			ERROR.Println(err)
			tx.Rollback()
			return err
		}
		err = tx.Commit()
		if err != nil {
			ERROR.Println(err)
			return err
		}

		err, tx = CreateTx()
		if err != nil {
			ERROR.Println(err)
			return err
		}
		//выполняем нужный отчет
		var buf *bytes.Buffer
		var err error = nil
		INFO.Println("building report", row.Code)
		switch row.Code {
		case reportType.ReceptionLog:
			buf, err = job.ReceptionLog(row, tx)
		case reportType.VisitsPerPeriod:
			buf, err = job.VisitsPerPeriod(row, tx)
		case reportType.AdmittedToTheHospital:
			buf, err = job.AdmittedToTheHospital(row, tx)
		case reportType.DischargedFromHospital:
			buf, err = job.DischargedFromHospital(row, tx)
		case reportType.Unvisited:
			buf, err = job.Unvisited(row, tx)
		case reportType.Registered:
			buf, err = job.Registered(row, tx)
		case reportType.Deregistered:
			buf, err = job.Deregistered(row, tx)
		case reportType.ConsistingOnTheSite:
			buf, err = job.ConsistingOnTheSite(row, tx)
		case reportType.ThoseInTheHospital:
			buf, err = job.ThoseInTheHospital(row, tx)
		case reportType.HospitalTreatment:
			buf, err = job.HospitalTreatment(row, tx)
		case reportType.AmbulatoryTreatment:
			buf, err = job.AmbulatoryTreatment(row, tx)
		case reportType.PBSTIN:
			buf, err = job.PBSTIN(row, tx)
		case reportType.TakenOnADN:
			buf, err = job.Registered(row, tx)
		case reportType.TakenFromADN:
			buf, err = job.Deregistered(row, tx)
		case reportType.TakenForADNAccordingToClinical:
			buf, err = job.TakenForADNAccordingToClinical(row, tx)
		case reportType.ProtocolUKL:
			buf, err = job.ProtocolUKL(row, tx)
		case reportType.Form39General:
			buf, err = job.Form39General(row, tx)

		default:
			err, tx = CreateTx()
			if err != nil {
				ERROR.Println(err)
				continue
			}
			INFO.Println("not found report", row.Code)
			_, err = setStatusByJob(row, statusType.error, tx)
			if err != nil {
				ERROR.Println(err)
				tx.Rollback()
				return err
			}
			tx.Commit()
			return err
		}
		if err != nil {
			ERROR.Println("error by report", row.Code)
			ERROR.Println("\t", err, row)
			_, err = setStatusByJob(row, statusType.error, tx)
			if err != nil {
				ERROR.Println(err)
				tx.Rollback()
				return err
			}
			tx.Commit()
			return err
		}
		//сохраняем отчет и помечаем что он выполнен
		//buf := toExcel(dataR)
		_, err = saveReport(row, buf, tx)
		if err != nil {
			ERROR.Println("save error by report", row.Code)
			ERROR.Println("\t", err)
			tx.Rollback()

			err, tx = CreateTx()
			if err != nil {
				ERROR.Println(err)
				return err
			}
			setStatusByJob(row, statusType.error, tx)
			tx.Commit()
			return err
		}
		INFO.Println("built report", row.Code)
		tx.Commit()
	}
	return nil
}

func process() error {
	var err error

	go func() {
		for true {
			err = createReport()
			if err != nil {
				continue
			}
			//visitLastDate("01.01.2022", "16.07.2022")
			<-time.Tick(time.Second * 10)
		}
	}()
	go func() {
		for true {
			err, tx := CreateTx()
			if err != nil {
				ERROR.Println(err)
			}

			_, err = deleteOlderReport(7, tx)
			if err != nil {
				ERROR.Println(err)
			}
			tx.Commit()

			<-time.Tick(time.Hour * 24)
		}
	}()
	return err
}
