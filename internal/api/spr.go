package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pmain2/internal/database"
	"pmain2/internal/models"
)

type sprApi struct{}

func sprApiInit() *sprApi {
	return &sprApi{}
}

func (s *sprApi) GetPodr(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return nil
	}

	conn, err := database.Connect()
	if err != nil {
		return err
	}
	defer conn.Close()

	model := models.CreateSpr(conn.DB)
	data, err := model.GetPodr()
	if err != nil {
		return err
	}

	res, err := json.Marshal(data)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, string(res))
	return nil

}
