package app

import (
	"strconv"
)

func isNumber(id string) bool {
	if _, err := strconv.Atoi(id); err == nil {
		return true
	} else {
		return false
	}
}

func isEmpty(str string) bool {
	if str == "" || len(str) == 0 {
		return true
	} else {
		return false
	}
}

func isEAN(str string) bool {
	if isNumber(str) && len(str) == 13 {
		return true
	} else {
		return false
	}
}
