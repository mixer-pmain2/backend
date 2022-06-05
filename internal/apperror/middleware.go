package apperror

import (
	"encoding/json"
	"errors"
	"net/http"
	"pmain2/internal/types"
	"pmain2/pkg/logger"
)

var (
	ERROR, _ = logger.New("errors", logger.ERROR)
)

type appHandler func(w http.ResponseWriter, r *http.Request) error

func Middleware(handle appHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var appErr *AppError

		err := handle(w, r)
		if err != nil {
			ERROR.Println(err.Error())
			if errors.As(err, &appErr) {
				if errors.Is(err, ErrDataNotFound) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("{}"))
					return
				}
			}
			//var status int
			//w.WriteHeader(status)
			//if status == http.StatusNotFound {
			//	fmt.Println("StatusNotFound")
			//}
			res, _ := json.Marshal(types.HttpResponse{
				Success: false,
				Error:   0,
				Message: "Ошибка выполнения запроса",
			})

			w.WriteHeader(http.StatusBadRequest)
			ERROR.Println(err.Error())
			w.Write(res)
		}
	}
}
