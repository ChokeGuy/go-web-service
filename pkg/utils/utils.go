package utils

import (
	"errors"
	"strconv"
)

func GetId(id string) (int, error) {
	value, err := strconv.Atoi(id)

	if err != nil {
		return 0, errors.New("invalid id")
	}
	return value, nil
}
