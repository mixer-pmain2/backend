package controller

import (
	"fmt"
	"strings"
	"time"

	"pmain2/internal/consts"
	"pmain2/internal/models"
	"pmain2/internal/types"
	"pmain2/pkg/cache"
	"pmain2/pkg/utils"
)

var (
	cachePat = cache.CreateCache(time.Minute, time.Minute)
)

type patient struct{}

func initPatientController() *patient {
	return &patient{}
}

func (p *patient) New(newPatient *types.NewPatient) (int, error, *[]types.Patient) {

	newPatient.Lname, _ = utils.ToWin1251(strings.ToUpper(newPatient.Lname))
	newPatient.Fname, _ = utils.ToWin1251(strings.ToUpper(newPatient.Fname))
	newPatient.Sname, _ = utils.ToWin1251(strings.ToUpper(newPatient.Sname))
	newPatient.Sex, _ = utils.ToWin1251(strings.ToUpper(newPatient.Sex))

	if newPatient.IsAnonim {
		newPatient.Lname = "-" + strings.Trim(newPatient.Lname, "-")
	}

	model := models.Model.Patient
	err, tx := models.Model.CreateTx()
	if err != nil {
		return 21, err, nil
	}
	defer tx.Rollback()

	if !newPatient.IsForced {
		found, err := model.FindByFIO(newPatient.Lname, newPatient.Fname, newPatient.Sname, tx)
		if err != nil {
			tx.Rollback()
			return -1, err, nil
		}
		if len(*found) > 0 {
			return 0, nil, found
		}
	}

	id, err := model.GetMaxPatientId(tx)
	if err != nil {
		tx.Rollback()
		return -1, err, nil
	}
	newPatient.PatientId = id + 1
	_, err = model.New(newPatient, tx)
	if err != nil {
		tx.Rollback()
		return -1, err, nil
	}
	newPatient.IsForced = true

	err = tx.Commit()
	if err != nil {
		return 22, err, nil
	}

	d, err := p.FindById(newPatient.PatientId, false)
	if err != nil {
		return 51, err, nil
	}

	data := make([]types.Patient, 0)
	data = append(data, *d)

	return 0, nil, &data
}

func (p *patient) FindByFio(lname, fname, sname string) (*[]types.Patient, error) {
	cacheName := lname + " " + fname + " " + sname

	item, ok := cachePat.Get(cacheName)
	if ok {
		res := item.(*[]types.Patient)
		return res, nil
	}

	model := models.Model.Patient
	err, tx := models.Model.CreateTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	lname, _ = utils.ToWin1251(lname)
	fname, _ = utils.ToWin1251(fname)
	sname, _ = utils.ToWin1251(sname)
	data, err := model.FindByFIO(lname, fname, sname, tx)
	if err != nil {
		return nil, err
	}
	tx.Commit()

	cachePat.Set(cacheName, data, 0)
	return data, nil
}

func (p *patient) FindByAddress(address types.Patient) (*[]types.Patient, error) {
	cacheName := address

	item, ok := cachePat.Get(cacheName)
	if ok {
		res := item.(*[]types.Patient)
		return res, nil
	}

	model := models.Model.Patient
	err, tx := models.Model.CreateTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	data, err := model.FindByAddress(address, tx)
	if err != nil {
		return nil, err
	}
	tx.Commit()

	cachePat.Set(cacheName, data, 0)
	return data, nil
}

func (p *patient) FindById(id int64, isCache bool) (*types.Patient, error) {
	cacheName := fmt.Sprintf("patient_id_%v", id)

	item, ok := cache.AppCache.Get(cacheName)
	if ok && isCache {
		return item.(*types.Patient), nil
	}

	err, tx := models.Model.CreateTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	model := models.Model.Patient
	data, err := model.Get(id, tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	cache.AppCache.Set(cacheName, data, 0)
	return data, nil
}

func (p *patient) FindUchet(id int64, isCache bool) (*[]models.FindUchetS, error) {
	cacheName := fmt.Sprintf("find_uchet_%v", id)
	item, ok := cache.AppCache.Get(cacheName)
	if ok && isCache {
		return item.(*[]models.FindUchetS), nil
	}
	model := models.Model.Patient
	err, tx := models.Model.CreateTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	data, err := model.FindUchet(id, 1000, 0, tx)
	if err != nil {
		ERROR.Println(err.Error())
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	cache.AppCache.Set(cacheName, data, 0)
	return data, nil
}

func (p *patient) FindLastUchet(id int64, isCache bool) (*models.FindUchetS, error) {
	cacheName := fmt.Sprintf("find_last_uchet_%v", id)
	item, ok := cache.AppCache.Get(cacheName)
	if ok && isCache {
		return item.(*models.FindUchetS), nil
	}
	model := models.Model.Patient
	err, tx := models.Model.CreateTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	data, err := model.FindLastUchet(id, tx)
	if err != nil {
		ERROR.Println(err.Error())
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	cache.AppCache.Set(cacheName, data, 30)
	return data, nil
}

func (p *patient) GetAddress(id int64, isCache bool) (string, error) {
	cacheName := fmt.Sprintf("patient_address_%v", id)
	item, ok := cache.AppCache.Get(cacheName)
	if ok && isCache {
		return item.(string), nil
	}

	model := models.Model.Patient
	err, tx := models.Model.CreateTx()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()
	data, err := model.GetAddress(id, tx)
	if err != nil {
		ERROR.Println(err.Error())
		return "", err
	}
	err = tx.Commit()
	if err != nil {
		return "", err
	}
	cache.AppCache.Set(cacheName, data, 30)
	return data, nil
}

func (p *patient) HistoryVisits(id int, isCache bool) (*[]models.HistoryVisit, error) {
	cacheName := fmt.Sprintf("disp_history_Visit_%v", id)
	if isCache {
		item, ok := cache.AppCache.Get(cacheName)
		if ok {
			return item.(*[]models.HistoryVisit), nil
		}
	}
	model := models.Model.Patient
	err, tx := models.Model.CreateTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	data, err := model.HistoryVisits(id, tx)
	if err != nil {
		ERROR.Println(err.Error())
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	cache.AppCache.Set(cacheName, data, 0)
	return data, nil
}

func (p *patient) HistoryHospital(id int) (*[]models.HistoryHospital, error) {
	cacheName := fmt.Sprintf("disp_history_hospital_%v", id)
	item, ok := cache.AppCache.Get(cacheName)
	if ok {
		return item.(*[]models.HistoryHospital), nil
	}
	model := models.Model.Patient
	err, tx := models.Model.CreateTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	data, err := model.HistoryHospital(id, tx)
	if err != nil {
		ERROR.Println(err.Error())
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	cache.AppCache.Set(cacheName, data, 0)
	return data, nil
}

func (p *patient) NewVisit(visit *types.NewVisit) (int, error) {
	visit.Normalize()
	model := models.Model.Patient
	err, tx := models.Model.CreateTx()
	if err != nil {
		return 20, err
	}
	defer tx.Rollback()
	lastUchet, err := model.FindLastUchet(visit.PatientId, tx)
	if err != nil {
		return 100, err
	}

	//-проверить что пациент не мертв или это работа с документами
	if (lastUchet != nil && lastUchet.Reason == consts.EXIT_REAS_DEAD) && visit.Visit&consts.VISIT_WORK_WITH_DOCUMENTS == 0 {
		return 101, nil
	}

	//-в этот день не было посещений
	isVisisted, err := model.IsVisited(visit, tx)
	if err != nil {
		tx.Rollback()
		return 102, err
	}

	if isVisisted {
		return 202, nil
	}

	err, tx = models.Model.CreateTx()
	if err != nil {
		return 21, err
	}
	defer tx.Rollback()

	//Обрезаем до 10, т.к. в посещениях длина диагноза 10
	if len(visit.Diagnose) > 10 {
		visit.Diagnose = visit.Diagnose[0:10]
	}

	model = models.Model.Patient
	_, err = model.NewVisit(*visit, tx)
	if err != nil {
		tx.Rollback()
		return 200, err
	}
	if visit.SRC >= 0 {
		_, err = model.NewSRC(&types.NewSRC{
			PatientId: visit.PatientId,
			DateAdd:   visit.Date,
			DockId:    visit.DockId,
			Unit:      visit.Unit,
			Zakl:      visit.SRC,
		}, tx)
		if err != nil {
			tx.Rollback()
			return 201, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return 22, err
	}

	return 0, nil
}

func (p *patient) NewProf(visit *types.NewProf) (int, error) {
	visit.Normalize()

	model := models.Model.Patient
	err, tx := models.Model.CreateTx()
	if err != nil {
		return 20, err
	}
	defer tx.Rollback()

	if visit.Count == 0 {
		return 203, nil
	}

	model = models.Model.Patient
	for i := 0; i < visit.Count; i++ {
		_, err = model.NewProf(*visit, tx)
		if err != nil {
			tx.Rollback()
			return 200, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return 22, err
	}

	return 0, nil
}

func (p *patient) NewReg(reg *types.NewRegister) (int, error) {
	pat, err := p.FindById(reg.PatientId, false)
	if err != nil {
		return -1, err
	}

	sprModel := models.Model.Spr
	patientModel := models.Model.Patient
	err, tx := models.Model.CreateTx()
	if err != nil {
		return 20, err
	}
	defer tx.Rollback()

	address, err := p.GetAddress(pat.Id, false)
	if err != nil {
		return -1, err
	}
	if reg.Reason == "001" {
		if len(address) < 10 {
			return 301, nil
		}
	}

	regDate, err := time.Parse("2006-01-02", reg.Date)
	if err != nil {
		return -1, err
	}

	lastReg, err := p.FindLastUchet(reg.PatientId, false)
	if err != nil {
		return -1, err
	}

	if lastReg.Date != "" {
		lastRegDate, err := time.Parse("2006-01-02", lastReg.Date)
		if err != nil {
			return -1, err
		}

		if lastRegDate.Sub(regDate) > 0 {
			return 360, nil
		}
		isClose, err := sprModel.IsClosedSection(lastReg.Section, tx)
		if err != nil {
			return -1, err
		}
		if isClose {
			return 303, nil
		}

		if reg.Diagnose == "" {
			reg.Diagnose = lastReg.Diagnose
		}

		if reg.Section == 0 {
			reg.Section = lastReg.Section
		}

		if reg.Category == 0 {
			reg.Category = lastReg.Category
		}

		if lastReg.Reason == consts.EXIT_REAS_DEAD {
			return 305, nil
		}

		if reg.Reason == consts.REAS_SWITCH_CATEG_GROUP && (reg.Category == 10 || reg.Category == lastReg.Category) {
			return 310, nil
		}

		if (reg.Reason == consts.REAS_SWITCH_CATEG_TO_AMBULANC && lastReg.Category > 0 && lastReg.Category < 9) ||
			(reg.Reason == consts.REAS_SWITCH_CATEG_TO_CONSULTANT && lastReg.Category == 10) {
			return 311, nil
		}
	}

	if reg.Category == 0 {
		return 317, err
	}

	countJudgment, err := patientModel.GetCountJudgment(reg.PatientId, tx)
	if err != nil {
		return -1, err
	}
	if reg.Category == 7 && countJudgment == 0 {
		return 304, nil
	}

	if reg.Section < 10 && reg.Reason == consts.REAS_NEW {
		return 306, nil
	}

	if (reg.Category == 7 || reg.Category == 8) && (reg.Section < 18 || reg.Section > 19) && reg.Section < 130 {
		return 307, nil
	}

	if reg.Reason == consts.REAS_SWITCH_CATEG_TO_AMBULANC && reg.Category == 10 {
		return 308, nil
	}

	if reg.Reason == consts.REAS_SWITCH_CATEG_TO_CONSULTANT && reg.Category != 10 {
		return 309, nil
	}

	if reg.Reason == consts.REAS_EXIT {
		if reg.ExitReason == "" {
			return 316, nil
		}
		reg.Reason = reg.ExitReason
	}

	if reg.Reason == consts.EXIT_REAS_NO_PSIH_DIAG {
		//TODO для всех prava <> 2147483647
		return 312, nil
	}

	isClose, err := sprModel.IsClosedSection(reg.Section, tx)
	if err != nil {
		return -1, err
	}
	if isClose {
		return 313, nil
	}

	if reg.Reason[0] == 'S' {
		inHospital, err := patientModel.IsInHospital(reg.PatientId, tx)
		if err != nil {
			return -1, err
		}
		if inHospital {
			return 314, nil
		}
	}

	countReg, err := patientModel.GetCountRegDataInDate(reg.PatientId, reg.Section, regDate, tx)
	if err != nil {
		return -1, err
	}
	if countReg > 0 {
		return 315, nil
	}

	_, err = patientModel.InsertReg(*reg, tx)
	if err != nil {
		tx.Rollback()
		return 350, err
	}
	if reg.Reason == consts.REAS_NEW {
		_, err = patientModel.UpdPatientVisible(reg.PatientId, 0, tx)
		if err != nil {
			return 351, err
		}
	}
	err = tx.Commit()
	if err != nil {
		return 350, err
	}

	return 0, nil
}

func (p *patient) NewRegisterTransfer(reg *types.NewRegisterTransfer) (int, error) {
	sprModel := models.Model.Spr
	patientModel := models.Model.Patient
	err, tx := models.Model.CreateTx()
	if err != nil {
		return 20, err
	}
	defer tx.Rollback()

	lastReg, err := p.FindLastUchet(reg.PatientId, false)
	if err != nil {
		return -1, err
	}

	isClose, err := sprModel.IsClosedSection(lastReg.Section, tx)
	if err != nil {
		return -1, err
	}
	if isClose {
		return 303, nil
	}

	isClose, err = sprModel.IsClosedSection(reg.Section, tx)
	if err != nil {
		return -1, err
	}
	if isClose {
		return 303, nil
	}

	if ((lastReg.Section == 19 && reg.Section == 18) || (lastReg.Section == 18 && reg.Section == 19)) && reg.Category < 7 {
		return 370, nil
	}
	if ((reg.Section == 18) || (reg.Section == 19)) && reg.Category < 7 {
		return 370, nil
	}
	if ((lastReg.Section == 18) || (lastReg.Section == 19)) && reg.Category > 6 {
		return 370, nil
	}

	reasonPrev := consts.REAS_FROM
	reasonNext := consts.REAS_TO

	if (lastReg.Section == 481 && reg.Section == 480) ||
		(lastReg.Section == 591 && reg.Section == 590) ||
		(lastReg.Section == 661 && reg.Section == 660) {
		reasonPrev = consts.REAS_FROM
		reasonNext = consts.REAS_TO
	} else {
		if lastReg.Section >= 400 && reg.Section < 400 {
			reasonPrev = consts.EXIT_REAS_EXIT
			reasonNext = consts.REAS_NEW
		}
		if lastReg.Section < 400 && reg.Section >= 400 {
			reasonPrev = consts.EXIT_REAS_EXIT
			reasonNext = consts.REAS_NEW
		}
	}
	regDate, err := time.Parse("2006-01-02", reg.Date)
	if err != nil {
		return -1, err
	}

	lastRegDate, err := time.Parse("2006-01-02", lastReg.Date)
	if err != nil {
		return -1, err
	}

	if lastRegDate.Sub(regDate) > 0 {
		return 360, nil
	}

	err, tx = models.Model.CreateTx()
	if err != nil {
		return 20, err
	}
	defer tx.Rollback()

	if err != nil {
		return -1, err
	}
	regPrev := types.NewRegister{
		PatientId:  reg.PatientId,
		Reason:     reasonPrev,
		ExitReason: "",
		Section:    lastReg.Section,
		Category:   lastReg.Category,
		Diagnose:   lastReg.Diagnose,
		Date:       reg.Date,
		DockId:     reg.DockId,
	}
	_, err = patientModel.InsertReg(regPrev, tx)
	if err != nil {
		return 350, err
	}
	regNext := types.NewRegister{
		PatientId:  reg.PatientId,
		Reason:     reasonNext,
		ExitReason: "",
		Section:    reg.Section,
		Category:   reg.Category,
		Diagnose:   lastReg.Diagnose,
		Date:       reg.Date,
		DockId:     reg.DockId,
	}
	_, err = patientModel.InsertReg(regNext, tx)
	if err != nil {
		return 350, err
	}
	err = tx.Commit()
	if err != nil {
		return 350, err
	}

	return 0, nil
}

func (p *patient) HistorySindrom(id int, isCache bool) (*[]models.HistorySindrom, error) {
	cacheName := fmt.Sprintf("disp_sindrom_%v", id)
	if isCache {
		item, ok := cache.AppCache.Get(cacheName)
		if ok {
			return item.(*[]models.HistorySindrom), nil
		}
	}
	err, tx := models.Model.CreateTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	data, err := models.Model.Patient.HistorySindrom(id, tx)
	if err != nil {
		ERROR.Println(err)
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		ERROR.Println(err)
		return nil, err
	}
	cache.AppCache.Set(cacheName, data, 0)
	return data, nil
}

func (p *patient) NewSindrom(sindrom *types.Sindrom) (int, error) {
	model := models.Model.Patient
	err, tx := models.Model.CreateTx()
	if err != nil {
		return 20, err
	}

	history, err := model.HistorySindrom(sindrom.PatientId, tx)
	if err != nil {
		return -1, err
	}

	isSindrom := strings.Contains(sindrom.Diagnose, "F")

	count := 0
	for _, row := range *history {
		isSindromRow := strings.Contains(row.Diagnose, "F")
		if (isSindrom && isSindromRow) || (!isSindrom && !isSindromRow) {
			count += 1
		}
	}

	if isSindrom && count >= 4 {
		return 380, nil
	}
	if !isSindrom && count >= 4 {
		return 381, nil
	}

	_, err = model.NewSindrom(*sindrom, tx)
	if err != nil {
		tx.Rollback()
		return 200, err
	}

	err = tx.Commit()
	if err != nil {
		return 22, err
	}

	return 0, nil
}

func (p *patient) RemoveSindrom(sindrom *types.Sindrom) (int, error) {

	err, tx := models.Model.CreateTx()
	if err != nil {
		ERROR.Println(err)
		return 20, nil
	}

	model := models.Model.Patient

	_, err = model.RemoveSindrom(*sindrom, tx)
	if err != nil {
		tx.Rollback()
		ERROR.Println(err)
		return 200, nil
	}

	err = tx.Commit()
	if err != nil {
		ERROR.Println(err)
		return 22, nil
	}

	return 0, nil
}

func (p *patient) FindInvalid(id int64, isCache bool) (*[]models.FindInvalid, error) {
	cacheName := fmt.Sprintf("find_invalid_%v", id)
	if isCache {
		item, ok := cache.AppCache.Get(cacheName)
		if ok {
			return item.(*[]models.FindInvalid), nil
		}
	}
	model := models.Model.Patient
	err, tx := models.Model.CreateTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	data, err := model.FindInvalid(id, tx)
	if err != nil {
		ERROR.Println(err.Error())
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	cache.AppCache.Set(cacheName, data, 0)
	return data, nil
}

func (p *patient) NewInvalid(newInvalid *types.NewInvalid) (int, error) {
	err, tx := models.Model.CreateTx()
	if err != nil {
		return 20, err
	}
	defer tx.Rollback()

	model := models.Model.Patient

	if newInvalid.IsInfinity {
		newInvalid.DateEnd = "2222-12-31"
	}

	d1, err := time.Parse("2006-01-02", newInvalid.DateStart)
	if err != nil {
		return 390, err
	}

	d2, err := time.Parse("2006-01-02", newInvalid.DateEnd)
	if err != nil {
		return 391, err
	}

	if d1.Sub(d2) > 0 {
		return 396, nil
	}

	if newInvalid.Kind == "10" {
		newInvalid.Reason = "1"
		_, err = model.NewChildInvalid(newInvalid, tx)
		if err != nil {
			tx.Rollback()
			return -1, err
		}
	}

	_, err = model.NewInvalid(newInvalid, tx)
	if err != nil {
		return -1, err
	}

	err = tx.Commit()
	if err != nil {
		return -1, err
	}
	return 0, nil
}

func (p *patient) UpdInvalid(newInvalid *types.NewInvalid) (int, error) {
	err, tx := models.Model.CreateTx()
	if err != nil {
		return 20, err
	}
	defer tx.Rollback()

	model := models.Model.Patient

	_, err = time.Parse("2006-01-02", newInvalid.DateDocument)
	if err != nil {
		return 397, err
	}

	invalids, err := model.FindInvalid(newInvalid.PatientId, tx)
	if err != nil {
		tx.Rollback()
		return 21, err
	}
	if len(*invalids) == 0 {
		return 398, nil
	}

	_, err = model.UpdInvalid(newInvalid, tx)
	if err != nil {
		tx.Rollback()
		return 21, err
	}

	err = tx.Commit()
	if err != nil {
		return 22, err
	}
	return 0, nil
}

func (p *patient) FindCustody(id int64, isCache bool) (*[]types.FindCustody, error) {
	cacheName := fmt.Sprintf("find_custody_%v", id)
	if isCache {
		item, ok := cache.AppCache.Get(cacheName)
		if ok {
			return item.(*[]types.FindCustody), nil
		}
	}
	model := models.Model.Patient
	err, tx := models.Model.CreateTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	data, err := model.FindCustody(id, tx)
	if err != nil {
		tx.Rollback()
		ERROR.Println(err.Error())
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	cache.AppCache.Set(cacheName, data, 0)
	return data, nil
}

func (p *patient) NewCustody(custody *types.NewCustody) (int, error) {

	_, err := time.Parse("2006-01-02", custody.DateStart)
	if err != nil {
		return 400, err
	}

	if custody.Custody == "" {
		return 402, nil
	}

	err, tx := models.Model.CreateTx()
	if err != nil {
		return 20, err
	}
	defer tx.Rollback()

	model := models.Model.Patient
	_, err = model.NewCustody(custody, tx)
	if err != nil {
		tx.Rollback()
		ERROR.Println(err.Error())
		return 21, err
	}

	err = tx.Commit()
	if err != nil {
		return 22, err
	}
	return 0, nil
}

func (p *patient) UpdCustody(custody *types.NewCustody) (int, error) {

	_, err := time.Parse("2006-01-02", custody.DateStart)
	if err != nil {
		return 400, err
	}

	_, err = time.Parse("2006-01-02", custody.DateEnd)
	if err != nil {
		return 401, err
	}

	model := models.Model.Patient
	err, tx := models.Model.CreateTx()
	if err != nil {
		return 21, err
	}
	defer tx.Rollback()

	_, err = model.UpdCustody(custody, tx)
	if err != nil {
		tx.Rollback()
		ERROR.Println(err.Error())
		return -1, err
	}

	err = tx.Commit()
	if err != nil {
		return -1, err
	}
	return 0, nil
}

func (p *patient) FindVaccination(id int64, isCache bool) (*[]types.FindVaccination, error) {
	cacheName := fmt.Sprintf("find_vaccination_%v", id)
	if isCache {
		item, ok := cache.AppCache.Get(cacheName)
		if ok {
			return item.(*[]types.FindVaccination), nil
		}
	}
	model := models.Model.Patient
	err, tx := models.Model.CreateTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	data, err := model.FindVaccination(id, tx)
	if err != nil {
		tx.Rollback()
		ERROR.Println(err.Error())
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	cache.AppCache.Set(cacheName, data, 0)
	return data, nil
}

func (p *patient) FindInfection(id int64, isCache bool) (*[]types.FindInfection, error) {
	cacheName := fmt.Sprintf("find_infection_%v", id)
	if isCache {
		item, ok := cache.AppCache.Get(cacheName)
		if ok {
			return item.(*[]types.FindInfection), nil
		}
	}
	model := models.Model.Patient
	err, tx := models.Model.CreateTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	data, err := model.FindInfection(id, tx)
	if err != nil {
		tx.Rollback()
		ERROR.Println(err.Error())
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	cache.AppCache.Set(cacheName, data, 0)
	return data, nil
}

func (p *patient) UpdPassport(passport *types.Patient) (int, error) {
	err, tx := models.Model.CreateTx()
	if err != nil {
		return 20, err
	}
	defer tx.Rollback()

	model := models.Model.Patient
	_, err = model.UpdPassport(passport, tx)
	if err != nil {
		tx.Rollback()
		ERROR.Println(err.Error())
		return 21, err
	}

	err = tx.Commit()
	if err != nil {
		return 22, err
	}
	return 0, nil
}

func (p *patient) UpdAddress(address *types.Patient) (int, error) {
	model := models.Model.Patient
	err, tx := models.Model.CreateTx()
	if err != nil {
		return 20, err
	}
	defer tx.Rollback()
	if err != nil {
		return -1, err
	}
	patient, err := model.Get(address.Id, tx)
	tx.Commit()
	if address.Republic == 0 {
		address.Republic = patient.Republic
	}
	if address.Region == 0 {
		address.Region = patient.Region
	}
	if address.District == 0 {
		address.District = patient.District
	}
	if address.Area == 0 {
		address.Area = patient.Area
	}
	if address.Street == 0 {
		address.Street = patient.Street
	}
	if address.Domicile == 0 {
		address.Domicile = patient.Domicile
	}

	address.Build, _ = utils.ToWin1251(address.Build)
	address.Flat, _ = utils.ToWin1251(address.Flat)
	address.House, _ = utils.ToWin1251(address.House)

	err, tx = models.Model.CreateTx()
	if err != nil {
		return 20, err
	}
	defer tx.Rollback()
	_, err = model.UpdAddress(address, tx)
	if err != nil {
		ERROR.Println(err.Error())
		return 9, err
	}

	err = tx.Commit()
	if err != nil {
		return 21, err
	}
	return 0, nil
}

func (p *patient) GetSection22(id int64, isCache bool) (*[]types.ST22, error) {
	cacheName := fmt.Sprintf("find_section22_%v", id)
	if isCache {
		item, ok := cache.AppCache.Get(cacheName)
		if ok {
			return item.(*[]types.ST22), nil
		}
	}
	model := models.Model.Patient
	err, tx := models.Model.CreateTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	data, err := model.GetSection22(id, tx)
	if err != nil {
		ERROR.Println(err.Error())
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	cache.AppCache.Set(cacheName, data, 0)
	return data, nil
}

func (p *patient) NewSection22(section *types.ST22) (int, error) {
	model := models.Model.Patient

	_, err := time.Parse("2006-01-02", section.DateStart)
	if err != nil {
		return 410, err
	}

	_, err = time.Parse("2006-01-02", section.DateEnd)
	if err != nil {
		return 411, err
	}

	err, tx := models.Model.CreateTx()
	if err != nil {
		return 20, err
	}
	defer tx.Rollback()
	_, err = model.NewSection22(section, tx)
	if err != nil {
		tx.Rollback()
		ERROR.Println(err.Error())
		return 21, err
	}

	err = tx.Commit()
	if err != nil {
		return 22, err
	}
	return 0, nil
}

func (p *patient) SOD(id int64, isCache bool) (*[]types.SOD, error) {
	cacheName := fmt.Sprintf("find_sod_%v", id)
	if isCache {
		item, ok := cache.AppCache.Get(cacheName)
		if ok {
			return item.(*[]types.SOD), nil
		}
	}
	model := models.Model.Patient
	err, tx := models.Model.CreateTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	data, err := model.SOD(id, tx)
	if err != nil {
		ERROR.Println(err.Error())
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	cache.AppCache.Set(cacheName, data, 0)
	return data, nil
}

func (p *patient) OODLast(id int64, isCache bool) (*types.OOD, error) {
	cacheName := fmt.Sprintf("find_ood_%v", id)
	if isCache {
		item, ok := cache.AppCache.Get(cacheName)
		if ok {
			return item.(*types.OOD), nil
		}
	}
	model := models.Model.Patient
	err, tx := models.Model.CreateTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	data, err := model.OODLast(id, tx)
	if err != nil {
		ERROR.Println(err.Error())
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	cache.AppCache.Set(cacheName, data, 0)
	return data, nil
}

func (p *patient) FindSection29(id int64, isCache bool) (*[]types.FindSection29, error) {
	cacheName := fmt.Sprintf("find_findSection29_%v", id)
	if isCache {
		item, ok := cache.AppCache.Get(cacheName)
		if ok {
			return item.(*[]types.FindSection29), nil
		}
	}
	model := models.Model.Patient
	err, tx := models.Model.CreateTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	data, err := model.FindSection29(id, tx)
	if err != nil {
		ERROR.Println(err.Error())
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	cache.AppCache.Set(cacheName, data, 0)
	return data, nil
}

func (p *patient) NewOOD(ood *types.OOD) (int, error) {
	model := models.Model.Patient

	err, tx := models.Model.CreateTx()
	if err != nil {
		return 20, err
	}
	defer tx.Rollback()
	_, err = model.NewOOD(ood, tx)
	if err != nil {
		tx.Rollback()
		ERROR.Println(err.Error())
		return 21, err
	}

	err = tx.Commit()
	if err != nil {
		return 22, err
	}
	return 0, nil
}

func (p *patient) NewSOD(sod *types.SOD) (int, error) {
	model := models.Model.Patient

	err, tx := models.Model.CreateTx()
	if err != nil {
		return 20, err
	}
	defer tx.Rollback()
	_, err = model.NewSOD(sod, tx)
	if err != nil {
		tx.Rollback()
		ERROR.Println(err.Error())
		return 21, err
	}

	err = tx.Commit()
	if err != nil {
		return 22, err
	}
	return 0, nil
}

func (p *patient) GetDoctorsVisitByPatient(id int, date time.Time, isCache bool) (*[]types.Doctor, error) {
	cacheName := fmt.Sprintf("GetDoctorsVisitByPatient%v%s", id, date)
	if isCache {
		item, ok := cache.AppCache.Get(cacheName)
		if ok {
			return item.(*[]types.Doctor), nil
		}
	}
	model := models.Model.Patient
	err, tx := models.Model.CreateTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	date = date.AddDate(0, 0, -365)
	data, err := model.GetDoctorsVisitByPatient(id, date, tx)
	if err != nil {
		ERROR.Println(err.Error())
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	cache.AppCache.Set(cacheName, data, 0)
	return data, nil
}

func (p *patient) GetLastUKLByVisitPatient(id int, isCache bool) (*types.UKLData, error) {
	cacheName := fmt.Sprintf("GetLastVisitByPatient%v", id)
	if isCache {
		item, ok := cache.AppCache.Get(cacheName)
		if ok {
			return item.(*types.UKLData), nil
		}
	}
	model := models.Model.Patient
	err, tx := models.Model.CreateTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	data, err := model.GetLastUKLByVisitPatient(id, tx)
	if err != nil {
		ERROR.Println(err.Error())
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	cache.AppCache.Set(cacheName, data, 0)
	return data, nil
}

func (p *patient) NewUKLByVisitPatient(data *types.NewUKL) (int, error) {

	if data.DoctorId == 0 {
		return 800, nil
	}

	model := models.Model.Patient
	err, tx := models.Model.CreateTx()
	if err != nil {
		return 21, err
	}
	defer tx.Rollback()
	lastUkl, err := model.GetLastUKLByVisitPatient(data.PatientId, tx)
	if err != nil {
		ERROR.Println(err.Error())
		return 21, err
	}
	if data.Level == 1 {
		dateLastUkl, _ := time.Parse(consts.DATE_FORMAT_DB, lastUkl.Date1)
		duration := time.Now().Sub(dateLastUkl)
		if data.Unit == consts.UNIT_APL {
			if duration < time.Hour*24*90 {
				return 801, nil
			}
		}
		if data.Unit == consts.UNIT_CHILD {
			if duration < time.Hour*24*30 {
				return 801, nil
			}
		}
		lastUchet, err := model.FindLastUchet(int64(data.PatientId), tx)
		if err != nil {
			return -1, err
		}
		if lastUchet.Id == 0 {
			return 805, nil
		}
		if data.Unit != consts.UNIT_APL && data.Unit != consts.UNIT_CHILD {
			if lastUchet.Id > 0 && lastUchet.Id == lastUkl.RegistratId {
				return 801, nil
			}
		}
		allowEditDate := false
		if service, err := models.Model.Spr.GetParams(tx); err != nil {
			for _, s := range *service {
				if s.Param == "UKL_EDIT_DATE" && s.ParamI == data.UserId {
					allowEditDate = true
					break
				}
			}
		}
		if !allowEditDate {
			data.Date = time.Now().Format(consts.DATE_FORMAT_DB)
		}
		_, err = model.NewUKLByVisitPatientLvl1(data, lastUchet.Id, tx)
		if err != nil {
			return -1, err
		}
	}
	data.Id = lastUkl.Id
	if data.Level == 2 {
		if lastUkl.Id == 0 {
			return 802, nil
		}
		_, err = model.NewUKLByVisitPatientLvl2(data, tx)
		if err != nil {
			return -1, err
		}
	}
	if data.Level == 3 {
		if lastUkl.Id == 0 {
			return 802, nil
		}
		if lastUkl.User2 == 0 {
			return 803, nil
		}
		_, err = model.NewUKLByVisitPatientLvl3(data, tx)
		if err != nil {
			return -1, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return 22, err
	}

	return 0, nil
}

func (p *patient) GetLastUKLBySuicide(id int, isCache bool) (*types.UKLData, error) {
	cacheName := fmt.Sprintf("GetLastUKLBySuicide%v", id)
	if isCache {
		item, ok := cache.AppCache.Get(cacheName)
		if ok {
			return item.(*types.UKLData), nil
		}
	}
	model := models.Model.Patient
	err, tx := models.Model.CreateTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	data, err := model.GetLastUKLBySuicide(id, tx)
	if err != nil {
		ERROR.Println(err.Error())
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	cache.AppCache.Set(cacheName, data, 0)
	return data, nil
}

func (p *patient) NewUKLBySuicide(data *types.NewUKL) (int, error) {

	if data.DoctorId == 0 {
		return 800, nil
	}

	model := models.Model.Patient
	err, tx := models.Model.CreateTx()
	if err != nil {
		return 21, err
	}
	defer tx.Rollback()
	lastUkl, err := model.GetLastUKLBySuicide(data.PatientId, tx)
	if err != nil {
		ERROR.Println(err.Error())
		return 21, err
	}
	lastVisit, err := model.CheckUKLLastVisit(data.PatientId, data.Unit, tx)
	if lastVisit.Id == 0 {
		return 804, nil
	}
	if data.Level == 1 {
		if err != nil {
			return -1, err
		}
		if lastVisit.Id == lastUkl.VisitId {
			return 801, nil
		}
		data.Date = time.Now().Format(consts.DATE_FORMAT_DB)
		_, err = model.NewUKLBySuicide1(data, lastVisit.Id, tx)
		if err != nil {
			return -1, err
		}
	}
	data.Id = lastUkl.Id
	if data.Level == 2 {
		if lastUkl.Id == 0 {
			return 802, nil
		}
		_, err = model.NewUKLBySuicide2(data, tx)
		if err != nil {
			return -1, err
		}
	}
	if data.Level == 3 {
		if lastUkl.Id == 0 {
			return 802, nil
		}
		if lastUkl.User2 == 0 {
			return 803, nil
		}
		_, err = model.NewUKLBySuicide3(data, tx)
		if err != nil {
			return -1, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return 22, err
	}

	return 0, nil
}

func (p *patient) GetLastUKLByPsychotherapy(id int, isCache bool) (*types.UKLData, error) {
	cacheName := fmt.Sprintf("GetLastUKLByPsychotherapy%v", id)
	if isCache {
		item, ok := cache.AppCache.Get(cacheName)
		if ok {
			return item.(*types.UKLData), nil
		}
	}
	model := models.Model.Patient
	err, tx := models.Model.CreateTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	data, err := model.GetLastUKLByPsychotherapy(id, tx)
	if err != nil {
		ERROR.Println(err.Error())
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	cache.AppCache.Set(cacheName, data, 0)
	return data, nil
}

func (p *patient) NewUKLByPsychotherapy(data *types.NewUKL) (int, error) {

	if data.DoctorId == 0 {
		return 800, nil
	}

	model := models.Model.Patient
	err, tx := models.Model.CreateTx()
	if err != nil {
		return 21, err
	}
	defer tx.Rollback()
	lastUkl, err := model.GetLastUKLByPsychotherapy(data.PatientId, tx)
	if err != nil {
		ERROR.Println(err.Error())
		return 21, err
	}
	lastVisit, err := model.CheckUKLLastVisit(data.PatientId, data.Unit, tx)
	if lastVisit.Id == 0 {
		return 804, nil
	}
	if data.Level == 1 {
		if err != nil {
			return -1, err
		}
		if lastVisit.Id == lastUkl.VisitId {
			return 801, nil
		}
		data.Date = time.Now().Format(consts.DATE_FORMAT_DB)
		_, err = model.NewUKLByPsychotherapy1(data, lastVisit.Id, tx)
		if err != nil {
			return -1, err
		}
	}
	data.Id = lastUkl.Id
	if data.Level == 2 {
		if lastUkl.Id == 0 {
			return 802, nil
		}
		_, err = model.NewUKLByPsychotherapy2(data, tx)
		if err != nil {
			return -1, err
		}
	}
	if data.Level == 3 {
		if lastUkl.Id == 0 {
			return 802, nil
		}
		if lastUkl.User2 == 0 {
			return 803, nil
		}
		_, err = model.NewUKLByPsychotherapy3(data, tx)
		if err != nil {
			return -1, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return 22, err
	}

	return 0, nil
}

func (p *patient) GetListUKLByPatient(id int, isType int, isCache bool) (*[]types.UKLData, error) {
	cacheName := fmt.Sprintf("GetListUKLByPatient%v_%v", id, isType)
	if isCache {
		item, ok := cache.AppCache.Get(cacheName)
		if ok {
			return item.(*[]types.UKLData), nil
		}
	}
	model := models.Model.Patient
	err, tx := models.Model.CreateTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	data, err := model.GetListUKLByPatient(id, isType, tx)
	if err != nil {
		ERROR.Println(err.Error())
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	cache.AppCache.Set(cacheName, data, 0)
	return data, nil
}

func (p *patient) GetForcedByPatient(id int, isCache bool) (*[]types.ForcedM, error) {
	cacheName := fmt.Sprintf("GetForcedByPatient%v", id)
	if isCache {
		item, ok := cache.AppCache.Get(cacheName)
		if ok {
			return item.(*[]types.ForcedM), nil
		}
	}
	model := models.Model.Patient
	err, tx := models.Model.CreateTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	data, err := model.GetForcedByPatient(id, tx)
	if err != nil {
		ERROR.Println(err.Error())
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	cache.AppCache.Set(cacheName, data, 0)
	return data, nil
}

func (p *patient) GetForcedLastByPatient(patientId int, isCache bool) (*types.Forced, error) {
	cacheName := fmt.Sprintf("GetForcedLastByPatient%v", patientId)
	if isCache {
		item, ok := cache.AppCache.Get(cacheName)
		if ok {
			return item.(*types.Forced), nil
		}
	}
	model := models.Model.Patient
	err, tx := models.Model.CreateTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	data, err := model.GetForcedLastByPatient(patientId, tx)
	if err != nil {
		ERROR.Println(err.Error())
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	cache.AppCache.Set(cacheName, data, 0)
	return data, nil
}

func (p *patient) GetViewed(id int, number int, isCache bool) (*[]types.ViewedM, error) {
	cacheName := fmt.Sprintf("GetViewed%v_%v", id, number)
	if isCache {
		item, ok := cache.AppCache.Get(cacheName)
		if ok {
			return item.(*[]types.ViewedM), nil
		}
	}
	model := models.Model.Patient
	err, tx := models.Model.CreateTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	data, err := model.GetViewed(id, number, tx)
	if err != nil {
		ERROR.Println(err.Error())
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	cache.AppCache.Set(cacheName, data, 0)
	return data, nil
}

func (p *patient) GetPolicy(id int, isCache bool) (*types.Policy, error) {
	cacheName := fmt.Sprintf("GetPolicy_%v", id)
	if isCache {
		item, ok := cache.AppCache.Get(cacheName)
		if ok {
			return item.(*types.Policy), nil
		}
	}
	model := models.Model.Patient
	err, tx := models.Model.CreateTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	data, err := model.GetPolicy(id, tx)
	if err != nil {
		ERROR.Println(err.Error())
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	cache.AppCache.Set(cacheName, data, 0)
	return data, nil
}

func (p *patient) UpdatePolicy(policy types.Policy) (int, error) {
	cacheName := fmt.Sprintf("GetPolicy_%v", policy.PatientId)
	cache.AppCache.Delete(cacheName)

	model := models.Model.Patient
	err, tx := models.Model.CreateTx()
	if err != nil {
		ERROR.Println(err.Error())
		return 22, err
	}
	defer tx.Rollback()
	_, err = model.UpdatePolicy(policy, tx)
	if err != nil {
		ERROR.Println(err.Error())
		return 22, err
	}
	err = tx.Commit()
	if err != nil {
		ERROR.Println(err)
		return 22, err
	}
	return 0, nil
}

func (p *patient) GetForced(id int, isCache bool) (*types.Forced, error) {
	cacheName := fmt.Sprintf("GetForced%v", id)
	if isCache {
		item, ok := cache.AppCache.Get(cacheName)
		if ok {
			return item.(*types.Forced), nil
		}
	}
	model := models.Model.Patient
	err, tx := models.Model.CreateTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	data, err := model.GetForced(id, tx)
	if err != nil {
		ERROR.Println(err.Error())
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	cache.AppCache.Set(cacheName, data, 0)
	return data, nil
}

func (p *patient) PostForcedByPatient(forced *types.Forced) (int, error) {
	dateNull := "1899-12-30"

	if forced.TypeId != 4 {
		if forced.DoctorId1 == 0 || forced.DoctorId2 == 0 {
			return 850, nil
		}
		if forced.ConclusionId == 0 {
			return 851, nil
		}
	}
	//if forced.CourtDate == "" {
	//	return 852, nil
	//}
	//if forced.CourtConclusionDate == "" {
	//	return 853, nil
	//}

	if forced.TypeId == 4 {
		forced.DoctorId1 = 0
		forced.DoctorId2 = 0
	}

	forcedLast, err := p.GetForcedLastByPatient(forced.PatientId, false)
	if forced.ConclusionId == 7 {
		forced.ConclusionId = forcedLast.ConclusionId
	}

	forced.DateEnd = ""
	if forced.ViewId == 5 || forced.ViewId == 6 {
		forced.DateEnd = forced.CourtDate
	}

	if forced.ViewId == 7 {
		forced.ViewId = forcedLast.ConclusionId
	}

	if forced.TypeId != 4 {
		if forced.ActNumber != 0 {
			forced.TypeId = 2
		} else {
			forced.TypeId = 3
		}
	}

	model := models.Model.Patient
	err, tx := models.Model.CreateTx()
	if err != nil {
		return 20, err
	}

	defer tx.Rollback()
	model = models.Model.Patient

	currentForced, _ := model.GetForcedNumberByPatient(forced.PatientId, forced.Number, tx)
	if forced.CourtConclusionDate == "" {
		forced.CourtConclusionDate = dateNull
	}
	if forced.CourtDate == "" {
		forced.CourtDate = dateNull
	}
	if forced.DateEnd == "" {
		forced.DateEnd = dateNull
	}
	if forced.DateView == "" {
		forced.DateView = dateNull
	}
	if forced.DateEnd == "" {
		forced.DateEnd = dateNull
	}
	if forced.ActDate == "" {
		forced.ActDate = dateNull
	}
	if dateView, _ := time.Parse(consts.DATE_FORMAT_DB, forced.DateView); true {
		actDate, _ := time.Parse(consts.DATE_FORMAT_DB, forced.ActDate)

		if dateView.Sub(actDate) < 0 {

		}
	}

	dateEnd, _ := time.Parse(time.RFC3339, currentForced.DateEnd)
	dNull, _ := time.Parse(consts.DATE_FORMAT_DB, dateNull)
	if dateEnd.Format(consts.DATE_FORMAT_DB) != dNull.Format(consts.DATE_FORMAT_DB) && currentForced.DateEnd != "" {
		return 855, nil
	}
	if forced.Number != forcedLast.Number {
		return 855, nil
	}

	if forced.Number == 0 {
		//forced.Number, err = model.GetNumForcedByPatient(forced.PatientId, tx)
		forced.Number = forcedLast.Number
	}
	if err != nil {
		return 22, err
	}
	if forced.Id > 0 {
		forcedCur, _ := model.GetForced(forced.Id, tx)
		if forcedCur.Id > 0 {
			//forced.Number = forcedCur.Number
			//_, err = model.DeleteForcedByViewDate(forced, tx)
			//if err != nil {
			//	return 22, err
			//}
			//_, err = model.PostForcedByPatient(forced, tx)
			_, err = model.UpdForcedByPatient(forced, tx)
			if err != nil {
				return 22, err
			}
		}
	}
	if forced.Id == 0 {
		_, err = model.PostForcedByPatient(forced, tx)
		if err != nil {
			return 22, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return 22, err
	}

	return 0, nil
}

func (p *patient) PostNewForcedByPatient(forced *types.Forced) (int, error) {
	dateNull := "1899-12-30"

	model := models.Model.Patient
	err, tx := models.Model.CreateTx()
	if err != nil {
		return 20, err
	}

	defer tx.Rollback()

	forced.Number, err = model.GetNumForcedByPatient(forced.PatientId, tx)
	if err != nil {
		return 854, nil
	}

	if forced.CourtConclusionDate == "" {
		forced.CourtConclusionDate = dateNull
	}
	if forced.CourtDate == "" {
		forced.CourtDate = dateNull
	}
	if forced.DateEnd == "" {
		forced.DateEnd = dateNull
	}
	if forced.DateView == "" {
		forced.DateView = dateNull
	}
	if forced.DateEnd == "" {
		forced.DateEnd = dateNull
	}
	if forced.ActDate == "" {
		forced.ActDate = dateNull
	}
	if dateView, _ := time.Parse(consts.DATE_FORMAT_DB, forced.DateView); true {
		actDate, _ := time.Parse(consts.DATE_FORMAT_DB, forced.ActDate)

		if dateView.Sub(actDate) < 0 {

		}
	}

	if sod, err := p.SOD(int64(forced.PatientId), true); err == nil {
		if len(*sod) == 0 {
			return 859, nil
		}
	} else {
		ERROR.Println(err)
		return 0, err
	}

	if forced.CourtConclusionDate == dateNull || forced.CourtDate == dateNull {
		return 860, nil
	}

	pol_date, _ := time.Parse(consts.DATE_FORMAT_INPUT, forced.CourtConclusionDate)
	op_date, _ := time.Parse(consts.DATE_FORMAT_INPUT, forced.CourtDate)
	fmt.Println(pol_date, op_date)

	//pol_date >=op_date
	if pol_date.Sub(op_date) < 0 {
		return 858, nil
	}

	if err != nil {
		return 22, err
	}
	_, err = model.PostForcedByPatient(forced, tx)
	if err != nil {
		return 22, err
	}

	err = tx.Commit()
	if err != nil {
		return 22, err
	}

	return 0, nil
}

func (p *patient) EndForcedByPatient(forced *types.Forced) (int, error) {
	dateNull := "1899-12-30"

	forcedLast, err := p.GetForcedLastByPatient(forced.PatientId, false)
	forced.Number = forcedLast.Number
	if forced.Number == 0 {
		return 854, nil
	}
	if forced.DateEnd == "" {
		return 856, nil
	}

	if forced.CourtDate == "" {
		return 857, nil
	}
	forced.CourtConclusionDate = forced.CourtDate
	forced.TypeId = 4

	model := models.Model.Patient
	err, tx := models.Model.CreateTx()
	if err != nil {
		return 20, err
	}

	defer tx.Rollback()

	//forced.Number, err = model.GetNumForcedByPatient(forced.PatientId, tx)
	if err != nil {
		return 854, nil
	}

	if forced.CourtConclusionDate == "" {
		forced.CourtConclusionDate = dateNull
	}
	if forced.CourtDate == "" {
		forced.CourtDate = dateNull
	}
	if forced.DateEnd == "" {
		forced.DateEnd = dateNull
	}
	if forced.DateView == "" {
		forced.DateView = dateNull
	}
	if forced.DateEnd == "" {
		forced.DateEnd = dateNull
	}
	if forced.ActDate == "" {
		forced.ActDate = dateNull
	}
	if dateView, _ := time.Parse(consts.DATE_FORMAT_DB, forced.DateView); true {
		actDate, _ := time.Parse(consts.DATE_FORMAT_DB, forced.ActDate)

		if dateView.Sub(actDate) < 0 {

		}
	}

	if err != nil {
		return 22, err
	}
	_, err = model.PostForcedByPatient(forced, tx)
	if err != nil {
		return 22, err
	}

	err = tx.Commit()
	if err != nil {
		return 22, err
	}

	return 0, nil
}

func (p *patient) GetNumForcedByPatient(id int) (int, error) {
	model := models.Model.Patient
	err, tx := models.Model.CreateTx()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	data, err := model.GetNumForcedByPatient(id, tx)
	if err != nil {
		ERROR.Println(err.Error())
		return 0, err
	}
	err = tx.Commit()
	if err != nil {
		return 0, err
	}
	return data, nil
}
