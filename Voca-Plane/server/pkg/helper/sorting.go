package helper

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

func ApplySorting(query *gorm.DB, sortBy, order string, allowed map[string]bool, defaultOrder string) *gorm.DB {
	if sortBy != "" && allowed[sortBy] {
		order = strings.ToLower(order)

		if order != "asc" && order != "desc" {
			order = "asc"
		}

		return query.Order(fmt.Sprintf("%s %s", sortBy, order))
	}

	return query.Order(defaultOrder)
}