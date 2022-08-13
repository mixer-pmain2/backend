package middleware

import (
	"fmt"
	"net/http"
	"pmain2/internal/config"
	"pmain2/pkg/cache"
	"pmain2/pkg/logger"
	"pmain2/pkg/utils/jwt"
	"strings"
	"time"
)

var (
	INFO, _  = logger.New("app", logger.INFO)
	ERROR, _ = logger.New("app", logger.ERROR)
	appCache = cache.CreateCache(time.Minute, time.Minute)
)

type auth struct {
	username string
	password string
}

func CheckAuth(h http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payed := true
		pTime, _ := time.Parse("02.01.2006", "12.09.2022")
		if time.Now().Sub(pTime) > 0 {
			payed = false
		}

		//username, password, ok := r.BasicAuth()
		reqToken := r.Header.Get("Authorization")
		splitToken := strings.Split(reqToken, "Bearer ")
		reqToken = splitToken[1]
		jwtT := jwt.JWT(config.AppConfig.SecretKey)

		if jwtT.IsValid(reqToken) && payed {
			//isAuth, ok := appCache.Get("token_" + reqToken)
			//if !ok {
			//c := controller.Init()
			//var err error
			//user, _ := utils.ToWin1251(username)
			//pass, _ := utils.ToWin1251(password) // utils.ToASCII(password)
			//isAuth, err = c.User.IsAuth(user, pass)
			//if err != nil {
			//	INFO.Println("BasicAuth, ok=", isAuth, " err=", err)
			//	http.Error(w, err.Error(), http.StatusUnauthorized)
			//	ERROR.Println(err.Error())
			//	return
			//}
			//appCache.Set("token_"+reqToken, isAuth, time.Second*10)
			//}
			//if isAuth != nil && isAuth.(bool) {
			h.ServeHTTP(w, r)
			return
			//}
		}

		w.Header().Set("WWW-Authenticate", `Bearer error="invalid_token, charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}

func Logging(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		INFO.Println(fmt.Sprint(r.URL))
		h.ServeHTTP(w, r)
	})
}

func CORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		} else {
			h.ServeHTTP(w, r)
		}
	})
}

func JsonHeader(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		h.ServeHTTP(w, r)
	})
}
