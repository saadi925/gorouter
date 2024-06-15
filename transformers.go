package flow

import "strconv"

func ToInt(value string) int {
	val, err := strconv.Atoi(value)
	if err != nil {
		return -1
	}
	return val
}
func ToString(value int) string {
	val := strconv.Itoa(value)
	return val
}
