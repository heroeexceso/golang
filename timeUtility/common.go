package timeutility

import (
	"time"

	"github.com/vigneshuvi/GoDateFormat"
)

// GetToday ... obtener la fecha del día
func GetToday() string {
	today := time.Now()
	return today.Format(GoDateFormat.ConvertFormat("dd-MM-yyyy HH:mm:ss"))
}

// GetYear ... obtener el año actual
func GetYear() string {
	today := time.Now()
	return today.Format(GoDateFormat.ConvertFormat("yyyy"))
}
