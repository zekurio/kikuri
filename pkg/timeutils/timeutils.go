package timeutils

import (
	"strconv"
	"strings"
	"time"
)

func ParseDuration(s string) (time.Duration, error) {
	if strings.HasSuffix(s, "d") {
		daysStr := strings.TrimSuffix(s, "d")
		days, err := strconv.Atoi(daysStr)
		if err != nil {
			return 0, err
		}
		hours := days * 24
		return time.Duration(hours) * time.Hour, nil
	}

	d, err := time.ParseDuration(s)
	if err != nil {
		return 0, err
	}
	return d, nil
}
