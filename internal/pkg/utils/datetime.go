package utils

import "time"

func TimeParse(t string) time.Duration {
	tt, err := time.ParseDuration(t)
	if err != nil {
		return 0
	}
	return tt
}