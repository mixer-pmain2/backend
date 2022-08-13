package types

type Doctor struct {
	Id     int    `json:"id"`
	Lname  string `json:"lname"`
	Fname  string `json:"fname"`
	Sname  string `json:"sname"`
	Access int    `json:"access"`
	Z152   int    `json:"z152"`
}

type DoctorFindParams struct {
	DoctorId int
	Month    int
	Year     int
	Unit     int
}

type DoctorRate struct {
	Unit     int     `json:"unit"`
	DoctorId int     `json:"doctorId"`
	Rate     float64 `json:"rate"`
	Id       int64   `json:"id"`
}

type DoctorVisitCountPlan struct {
	Unit  int     `json:"unit"`
	Visit int64   `json:"visit"`
	Plan  float64 `json:"plan"`
}

type DoctorQueryUpdRate struct {
	Year     int    `json:"year"`
	Month    int    `json:"month"`
	Unit     int    `json:"unit"`
	Rate     string `json:"rate"`
	DoctorId int    `json:"doctorId"`
}
