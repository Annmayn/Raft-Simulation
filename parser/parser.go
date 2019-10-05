package parser

import (
	"strconv"
	"strings"
)

func parseMessage(s string) int {
	i := 0
	for _, alpha := range strings.Split(s, "") {
		t, _ := strconv.Atoi(alpha)
		if t != 0 {
			i++
		} else {
			if alpha == "0" {
				i++
			} else {
				break
			}
		}
	}
	arg1, _ := strconv.Atoi(s[:i])
	arg2, _ := strconv.Atoi(s[i+1:])
	var res int
	switch string(s[i]) {
	case "+":
		res = arg1 + arg2
	case "-":
		res = arg1 - arg2
	case "*":
		res = arg1 * arg2
	case "/":
		res = arg1 / arg2
	default:
		res = 0
	}
	return res
}
