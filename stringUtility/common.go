package stringutility

import (
	"strconv"
)

// IntToString ... devolver un numero entero en formato string
func IntToString(value int) string {
	return strconv.Itoa(value)
}
