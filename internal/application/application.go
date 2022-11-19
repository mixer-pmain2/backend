package application

import (
	"log"
	"net/http"
	"pmain2/internal/api"
	"pmain2/internal/apperror"
	"pmain2/internal/middleware"
	"pmain2/internal/server"
	"pmain2/internal/web"
	"pmain2/pkg/logger"
	"pmain2/pkg/version"
)

const (
	Version version.Version = "0.0.0"
)

var (
	INFO  *log.Logger
	ERROR *log.Logger
)

func InitLogger() {
	INFO, _ = logger.New("app", logger.INFO)
	ERROR, _ = logger.New("app", logger.ERROR)
}

func CreateRouters(server *server.Server) {
	INFO.Println("Init routes")
	apiHandlers := api.Init()
	apiRouter := server.Router.PathPrefix("/api/v0").Subrouter()
	apiRouter.Use(middleware.CORS)
	apiRouter.Use(middleware.CheckAuth)
	apiRouter.Use(middleware.JsonHeader)
	apiRouter.Use(middleware.Logging)

	apiRouter.HandleFunc("/user/{id:[0-9]*}/uch/", apperror.Middleware(apiHandlers.User.GetUch)).Methods(http.MethodGet, http.MethodOptions)
	apiRouter.HandleFunc("/user/{id:[0-9]*}/prava/", apperror.Middleware(apiHandlers.User.GetPrava)).Methods(http.MethodGet, http.MethodOptions)
	apiRouter.HandleFunc("/user/{id:[0-9]*}/changepassword", apperror.Middleware(apiHandlers.User.ChangePassword)).Methods(http.MethodPost, http.MethodOptions)
	apiRouter.HandleFunc("/user/{id:[0-9]*}/", apperror.Middleware(apiHandlers.User.GetUser)).Methods(http.MethodGet, http.MethodOptions)
	apiRouter.HandleFunc("/auth/signin/", apperror.Middleware(apiHandlers.User.Signin)).Methods(http.MethodGet, http.MethodOptions)
	apiRouter.HandleFunc("/patient/find/", apperror.Middleware(apiHandlers.Patient.Find)).Methods(http.MethodGet, http.MethodOptions)
	apiRouter.HandleFunc("/patient/findByAddress/", apperror.Middleware(apiHandlers.Patient.FindByAddress)).Methods(http.MethodGet, http.MethodOptions)
	apiRouter.HandleFunc("/patient/new/", apperror.Middleware(apiHandlers.Patient.New)).Methods(http.MethodPost, http.MethodOptions)
	apiRouter.HandleFunc("/patient/{id:[0-9]*}/hospital/", apperror.Middleware(apiHandlers.Patient.HistoryHospital)).Methods(http.MethodGet, http.MethodOptions)
	apiRouter.HandleFunc("/patient/{id:[0-9]*}/uchet/", apperror.Middleware(apiHandlers.Patient.FindUchet)).Methods(http.MethodGet, http.MethodOptions)
	apiRouter.HandleFunc("/patient/{id:[0-9]*}/uchet/", apperror.Middleware(apiHandlers.Patient.NewReg)).Methods(http.MethodPost, http.MethodOptions)
	apiRouter.HandleFunc("/patient/{id:[0-9]*}/uchet/transfer/", apperror.Middleware(apiHandlers.Patient.NewRegTransfer)).Methods(http.MethodPost, http.MethodOptions)
	apiRouter.HandleFunc("/patient/{id:[0-9]*}/visit/", apperror.Middleware(apiHandlers.Patient.HistoryVisits)).Methods(http.MethodGet, http.MethodOptions)
	apiRouter.HandleFunc("/patient/{id:[0-9]*}/visit/", apperror.Middleware(apiHandlers.Patient.NewVisit)).Methods(http.MethodPost, http.MethodOptions)
	apiRouter.HandleFunc("/patient/{id:[0-9]*}/syndrome/", apperror.Middleware(apiHandlers.Patient.GetSindrom)).Methods(http.MethodGet, http.MethodOptions)
	apiRouter.HandleFunc("/patient/{id:[0-9]*}/syndrome/", apperror.Middleware(apiHandlers.Patient.NewSindrom)).Methods(http.MethodPost, http.MethodOptions)
	apiRouter.HandleFunc("/patient/{id:[0-9]*}/syndrome/", apperror.Middleware(apiHandlers.Patient.RemoveSindrom)).Methods(http.MethodDelete, http.MethodOptions)
	apiRouter.HandleFunc("/patient/{id:[0-9]*}/invalid/", apperror.Middleware(apiHandlers.Patient.FindInvalid)).Methods(http.MethodGet, http.MethodOptions)
	apiRouter.HandleFunc("/patient/{id:[0-9]*}/invalid/", apperror.Middleware(apiHandlers.Patient.NewInvalid)).Methods(http.MethodPost, http.MethodOptions)
	apiRouter.HandleFunc("/patient/{id:[0-9]*}/invalid/", apperror.Middleware(apiHandlers.Patient.UpdInvalid)).Methods(http.MethodPut, http.MethodOptions)
	apiRouter.HandleFunc("/patient/{id:[0-9]*}/custody/", apperror.Middleware(apiHandlers.Patient.FindCustody)).Methods(http.MethodGet, http.MethodOptions)
	apiRouter.HandleFunc("/patient/{id:[0-9]*}/custody/", apperror.Middleware(apiHandlers.Patient.NewCustody)).Methods(http.MethodPost, http.MethodOptions)
	apiRouter.HandleFunc("/patient/{id:[0-9]*}/custody/", apperror.Middleware(apiHandlers.Patient.UpdCustody)).Methods(http.MethodPut, http.MethodOptions)
	apiRouter.HandleFunc("/patient/{id:[0-9]*}/vaccination/", apperror.Middleware(apiHandlers.Patient.FindVaccination)).Methods(http.MethodGet, http.MethodOptions)
	apiRouter.HandleFunc("/patient/{id:[0-9]*}/infection/", apperror.Middleware(apiHandlers.Patient.FindInfection)).Methods(http.MethodGet, http.MethodOptions)
	apiRouter.HandleFunc("/patient/{id:[0-9]*}/passport/", apperror.Middleware(apiHandlers.Patient.UpdPassport)).Methods(http.MethodPut, http.MethodOptions)
	apiRouter.HandleFunc("/patient/{id:[0-9]*}/address/", apperror.Middleware(apiHandlers.Patient.GetAddress)).Methods(http.MethodGet, http.MethodOptions)
	apiRouter.HandleFunc("/patient/{id:[0-9]*}/address/", apperror.Middleware(apiHandlers.Patient.UpdAddress)).Methods(http.MethodPut, http.MethodOptions)
	apiRouter.HandleFunc("/patient/{id:[0-9]*}/section22/", apperror.Middleware(apiHandlers.Patient.GetSection22)).Methods(http.MethodGet, http.MethodOptions)
	apiRouter.HandleFunc("/patient/{id:[0-9]*}/section22/", apperror.Middleware(apiHandlers.Patient.NewSection22)).Methods(http.MethodPost, http.MethodOptions)
	apiRouter.HandleFunc("/patient/{id:[0-9]*}/sod/", apperror.Middleware(apiHandlers.Patient.SOD)).Methods(http.MethodGet, http.MethodOptions)
	apiRouter.HandleFunc("/patient/{id:[0-9]*}/sod/", apperror.Middleware(apiHandlers.Patient.NewSOD)).Methods(http.MethodPost, http.MethodOptions)
	apiRouter.HandleFunc("/patient/{id:[0-9]*}/ood/last/", apperror.Middleware(apiHandlers.Patient.OODLast)).Methods(http.MethodGet, http.MethodOptions)
	apiRouter.HandleFunc("/patient/{id:[0-9]*}/ood/", apperror.Middleware(apiHandlers.Patient.NewOOD)).Methods(http.MethodPost, http.MethodOptions)
	apiRouter.HandleFunc("/patient/{id:[0-9]*}/section29/find/", apperror.Middleware(apiHandlers.Patient.FindSection29)).Methods(http.MethodGet, http.MethodOptions)
	apiRouter.HandleFunc("/patient/{id:[0-9]*}/forced/", apperror.Middleware(apiHandlers.Patient.GetForcedByPatient)).Methods(http.MethodGet, http.MethodOptions)
	apiRouter.HandleFunc("/patient/{id:[0-9]*}/forced/last/", apperror.Middleware(apiHandlers.Patient.GetForcedLastByPatient)).Methods(http.MethodGet, http.MethodOptions)
	apiRouter.HandleFunc("/patient/{id:[0-9]*}/forced/number/", apperror.Middleware(apiHandlers.Patient.GetNumForcedByPatient)).Methods(http.MethodGet, http.MethodOptions)
	apiRouter.HandleFunc("/patient/{id:[0-9]*}/forced/", apperror.Middleware(apiHandlers.Patient.PostForcedByPatient)).Methods(http.MethodPost, http.MethodOptions)
	apiRouter.HandleFunc("/patient/{id:[0-9]*}/forced/new/", apperror.Middleware(apiHandlers.Patient.PostNewForcedByPatient)).Methods(http.MethodPost, http.MethodOptions)
	apiRouter.HandleFunc("/patient/{id:[0-9]*}/forced/end/", apperror.Middleware(apiHandlers.Patient.EndForcedByPatient)).Methods(http.MethodPost, http.MethodOptions)
	apiRouter.HandleFunc("/patient/forced/", apperror.Middleware(apiHandlers.Patient.GetForced)).Methods(http.MethodGet, http.MethodOptions)
	apiRouter.HandleFunc("/patient/{id:[0-9]*}/viewed/", apperror.Middleware(apiHandlers.Patient.GetViewed)).Methods(http.MethodGet, http.MethodOptions)
	apiRouter.HandleFunc("/patient/{id:[0-9]*}/policy/", apperror.Middleware(apiHandlers.Patient.GetPolicy)).Methods(http.MethodGet, http.MethodOptions)
	apiRouter.HandleFunc("/patient/{id:[0-9]*}/policy/", apperror.Middleware(apiHandlers.Patient.UpdatePolicy)).Methods(http.MethodPut, http.MethodOptions)
	apiRouter.HandleFunc("/patient/{id:[0-9]*}/", apperror.Middleware(apiHandlers.Patient.Get)).Methods(http.MethodGet, http.MethodOptions)
	apiRouter.HandleFunc("/patient/prof/", apperror.Middleware(apiHandlers.Patient.NewProf)).Methods(http.MethodPost, http.MethodOptions)
	apiRouter.HandleFunc("/administration/doctor/location/", apperror.Middleware(apiHandlers.Administration.DoctorLocation)).Methods(http.MethodPost, http.MethodOptions)
	apiRouter.HandleFunc("/administration/doctor/lead/", apperror.Middleware(apiHandlers.Administration.DoctorLeadSection)).Methods(http.MethodPost, http.MethodOptions)
	apiRouter.HandleFunc("/doctor/{id:[0-9]*}/rate", apperror.Middleware(apiHandlers.Doctor.GetRate)).Methods(http.MethodGet, http.MethodOptions)
	apiRouter.HandleFunc("/doctor/{id:[0-9]*}/rate", apperror.Middleware(apiHandlers.Doctor.UpdRate)).Methods(http.MethodPut, http.MethodOptions)
	apiRouter.HandleFunc("/doctor/{id:[0-9]*}/rate", apperror.Middleware(apiHandlers.Doctor.DelRate)).Methods(http.MethodDelete, http.MethodOptions)
	apiRouter.HandleFunc("/doctor/{id:[0-9]*}/visit/count", apperror.Middleware(apiHandlers.Doctor.VisitCountPlan)).Methods(http.MethodGet, http.MethodOptions)
	apiRouter.HandleFunc("/doctor/{id:[0-9]*}/units", apperror.Middleware(apiHandlers.Doctor.GetUnits)).Methods(http.MethodGet, http.MethodOptions)
	apiRouter.HandleFunc("/ukl/doctors/visit/", apperror.Middleware(apiHandlers.Patient.GetDoctorsVisitByPatient)).Methods(http.MethodGet, http.MethodOptions)
	apiRouter.HandleFunc("/ukl/visit/", apperror.Middleware(apiHandlers.Patient.GetLastUKLByVisitPatient)).Methods(http.MethodGet, http.MethodOptions)
	apiRouter.HandleFunc("/ukl/visit/", apperror.Middleware(apiHandlers.Patient.SetUKLByVisitPatient)).Methods(http.MethodPost, http.MethodOptions)
	apiRouter.HandleFunc("/ukl/suicide/", apperror.Middleware(apiHandlers.Patient.GetLastUKLBySuicide)).Methods(http.MethodGet, http.MethodOptions)
	apiRouter.HandleFunc("/ukl/suicide/", apperror.Middleware(apiHandlers.Patient.SetUKLBySuicide)).Methods(http.MethodPost, http.MethodOptions)
	apiRouter.HandleFunc("/ukl/psychotherapy/", apperror.Middleware(apiHandlers.Patient.GetLastUKLByPsychotherapy)).Methods(http.MethodGet, http.MethodOptions)
	apiRouter.HandleFunc("/ukl/psychotherapy/", apperror.Middleware(apiHandlers.Patient.SetUKLByPsychotherapy)).Methods(http.MethodPost, http.MethodOptions)
	apiRouter.HandleFunc("/ukl/", apperror.Middleware(apiHandlers.Patient.GetListUKLByPatient)).Methods(http.MethodGet, http.MethodOptions)

	apiRouterNonAuth := server.Router.PathPrefix("/api/v0").Subrouter()
	apiRouterNonAuth.Use(middleware.CORS)
	apiRouterNonAuth.Use(middleware.JsonHeader)
	apiRouterNonAuth.Use(middleware.Logging)
	apiRouterNonAuth.HandleFunc("/auth/login/", apperror.Middleware(apiHandlers.User.Login)).Methods(http.MethodPost, http.MethodOptions)
	apiRouterNonAuth.HandleFunc("/report/", apperror.Middleware(apiHandlers.Report.HandleList)).Methods(http.MethodGet, http.MethodOptions)
	apiRouterNonAuth.HandleFunc("/report/download/", apperror.Middleware(apiHandlers.Report.HandleDownload)).Methods(http.MethodGet, http.MethodOptions)
	apiRouterNonAuth.HandleFunc("/report/order/", apperror.Middleware(apiHandlers.Report.HandleOrder)).Methods(http.MethodPost, http.MethodOptions)
	apiRouterNonAuth.HandleFunc("/spr/podr/", apperror.Middleware(apiHandlers.Spr.GetPodr)).Methods(http.MethodGet, http.MethodOptions)
	apiRouterNonAuth.HandleFunc("/spr/prava/", apperror.Middleware(apiHandlers.Spr.GetPrava)).Methods(http.MethodGet, http.MethodOptions)
	apiRouterNonAuth.HandleFunc("/spr/visit/", apperror.Middleware(apiHandlers.Spr.GetSprVisit)).Methods(http.MethodGet, http.MethodOptions)
	apiRouterNonAuth.HandleFunc("/spr/visit/code/", apperror.Middleware(apiHandlers.Spr.GetSprVisitByCode)).Methods(http.MethodGet, http.MethodOptions)
	apiRouterNonAuth.HandleFunc("/spr/diag/", apperror.Middleware(apiHandlers.Spr.GetSprDiags)).Methods(http.MethodGet, http.MethodOptions)
	apiRouterNonAuth.HandleFunc("/spr/reason/", apperror.Middleware(apiHandlers.Spr.GetSprReasons)).Methods(http.MethodGet, http.MethodOptions)
	apiRouterNonAuth.HandleFunc("/spr/invalid/kind/", apperror.Middleware(apiHandlers.Spr.GetSprInvalidKind)).Methods(http.MethodGet, http.MethodOptions)
	apiRouterNonAuth.HandleFunc("/spr/invalid/anomaly/", apperror.Middleware(apiHandlers.Spr.GetSprInvalidChildAnomaly)).Methods(http.MethodGet, http.MethodOptions)
	apiRouterNonAuth.HandleFunc("/spr/invalid/limit/", apperror.Middleware(apiHandlers.Spr.GetSprInvalidChildLimit)).Methods(http.MethodGet, http.MethodOptions)
	apiRouterNonAuth.HandleFunc("/spr/invalid/reason/", apperror.Middleware(apiHandlers.Spr.GetSprInvalidReason)).Methods(http.MethodGet, http.MethodOptions)
	apiRouterNonAuth.HandleFunc("/spr/custody/who/", apperror.Middleware(apiHandlers.Spr.GetSprCustodyWho)).Methods(http.MethodGet, http.MethodOptions)
	apiRouterNonAuth.HandleFunc("/spr/address/republic/", apperror.Middleware(apiHandlers.Spr.FindRepublic)).Methods(http.MethodGet, http.MethodOptions)
	apiRouterNonAuth.HandleFunc("/spr/address/region/", apperror.Middleware(apiHandlers.Spr.FindRegion)).Methods(http.MethodGet, http.MethodOptions)
	apiRouterNonAuth.HandleFunc("/spr/address/district/", apperror.Middleware(apiHandlers.Spr.FindDistrict)).Methods(http.MethodGet, http.MethodOptions)
	apiRouterNonAuth.HandleFunc("/spr/address/area/", apperror.Middleware(apiHandlers.Spr.FindArea)).Methods(http.MethodGet, http.MethodOptions)
	apiRouterNonAuth.HandleFunc("/spr/address/street/", apperror.Middleware(apiHandlers.Spr.FindStreet)).Methods(http.MethodGet, http.MethodOptions)
	apiRouterNonAuth.HandleFunc("/spr/section/", apperror.Middleware(apiHandlers.Spr.FindSections)).Methods(http.MethodGet, http.MethodOptions)
	apiRouterNonAuth.HandleFunc("/spr/doctor/section/", apperror.Middleware(apiHandlers.Spr.FindSectionDoctor)).Methods(http.MethodGet, http.MethodOptions)
	apiRouterNonAuth.HandleFunc("/spr/doctor/lead/", apperror.Middleware(apiHandlers.Spr.FindSectionLead)).Methods(http.MethodGet, http.MethodOptions)
	apiRouterNonAuth.HandleFunc("/spr/doctors/", apperror.Middleware(apiHandlers.Spr.GetDoctors)).Methods(http.MethodGet, http.MethodOptions)
	apiRouterNonAuth.HandleFunc("/service/", apperror.Middleware(apiHandlers.Spr.GetParams)).Methods(http.MethodGet, http.MethodOptions)

	webRouter := server.Router.PathPrefix("/").Subrouter()
	webRouter.PathPrefix("/static/").Handler(http.FileServer(http.Dir("./")))
	webRouter.Use(middleware.Logging)

	web.RoutesFrontend(webRouter, web.IndexServe)

}
