package report

import (
	"encoding/json"
	"net/http"
	"pmain2/internal/consts"
	"pmain2/internal/types"
	"pmain2/pkg/logger"
	"strconv"
)

var (
	INFO, _  = logger.New("report", logger.INFO)
	ERROR, _ = logger.New("report", logger.ERROR)
)

func resSuccess(val int) []byte {
	res := types.HttpResponse{Success: true, Error: 0}

	if val > 0 {
		res.Success = false
		res.Error = val
		res.Message = consts.ArrErrors[val]
	}
	resMarshal, _ := json.Marshal(res)
	return resMarshal
}

type Api struct {
}

func getOrderReportParams(r *http.Request) *reportParams {
	p := reportParams{}
	json.NewDecoder(r.Body).Decode(&p)
	return &p
}

func (h *Api) HandleOrder(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return nil
	}

	param := getOrderReportParams(r)

	if param.Code == "" {
		w.Write(resSuccess(701))
		return nil
	}

	err, tx := CreateTx()
	defer tx.Rollback()
	if err != nil {
		ERROR.Println(err)
		return err
	}

	_, err = newJob(param, tx)
	if err != nil {
		val := 700
		res, _ := json.Marshal(orderResponse{Success: false, Error: val, Message: consts.ArrErrors[val]})
		ERROR.Println(err)
		w.Write(res)
		return err
	}
	err = tx.Commit()
	if err != nil {
		val := 700
		res, _ := json.Marshal(orderResponse{Success: false, Error: val, Message: consts.ArrErrors[val]})
		ERROR.Println(err)
		w.Write(res)
		return err
	}

	res, _ := json.Marshal(orderResponse{Success: true})
	w.Write(res)
	return nil
}

func (h *Api) HandleDownload(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return nil
	}

	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	_, tx := CreateTx()

	byteFile, err := getJob(id, tx)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	_, err = w.Write(*byteFile)
	if err != nil {
		ERROR.Println(err)
		return err
	}
	return nil
}

func (h *Api) HandleList(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return nil
	}

	userId, err := strconv.Atoi(r.URL.Query().Get("userId"))

	err, tx := CreateTx()
	defer tx.Rollback()
	if err != nil {
		ERROR.Println(err)
		return err
	}

	data, err := getJobs(userId, tx)
	if err != nil {
		ERROR.Println(err)
		return err
	}

	res, _ := json.Marshal(data)
	w.Write(res)
	return nil
}
