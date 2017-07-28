package api

import "strconv"

func Atoi(str string, defaultInt int) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		return defaultInt
	}
	return i
}
