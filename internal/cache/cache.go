package cache

import "fmt"

var (
	DoctorRate = func(DoctorId, Year, Month, Unit int) string {
		return fmt.Sprintf("doctor_rate_%v_%v_%v_%v", DoctorId, Year, Month, Unit)
	}
)
