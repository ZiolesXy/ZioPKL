package config

import (
	"os"
	"strconv"
)

func GetBcryptCost() int {
	cost := os.Getenv("BCRYPT_COST")
	if cost == "" {
		return 10
	}

	costInt, err := strconv.Atoi(cost)
	if err != nil || costInt < 4 || costInt > 31 {
		return 10
	}

	return costInt
}