package common

import (
	"strconv"
	"strings"
)

func String_to_int(s []string) ([]int, error) {
	var num []int
	for _, char := range s {
		d, err := strconv.Atoi(strings.TrimSpace(char))
		if err != nil {
			return []int{}, err
		}
		num = append(num, d)
	}
	return num, nil
}
