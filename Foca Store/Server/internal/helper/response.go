package helper

import (
	"reflect"

	"github.com/gin-gonic/gin"
)

func WrapListIfNeeded(data interface{}) interface{} {
	v := reflect.ValueOf(data)

	// Jika nil, return langsung
	if !v.IsValid() {
		return data
	}

	// Cek apakah slice atau array → bungkus sebagai entries
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		return gin.H{
			"entries": data,
		}
	}

	// Bukan slice/array → return apa adanya
	return data
}