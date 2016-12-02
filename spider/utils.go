package spider

import (
	"strconv"
)

func FindMaxFromSliceString(min int, sofs []string) int {
	max := min
	for _, v := range sofs {
		n, _ := strconv.Atoi(v)
		if n > max {
			max = n
		}
	}
	return max
}
