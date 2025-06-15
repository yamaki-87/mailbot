package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

var (
	WEEKDAYJP = [7]string{"日", "月", "火", "水", "木", "金", "土"}
)

const (
	SEPERATE_SLASH = "/"
)

func GetWeekDayFromDate(mmdd string) string {
	t, err := SepearteMMDD(mmdd)
	if err != nil {
		log.Warn().Err(err)
		return mmdd
	}
	return WEEKDAYJP[t.Weekday()]
}

func SepearteMMDD(mmdd string) (*time.Time, error) {
	parts := strings.Split(mmdd, SEPERATE_SLASH)
	if len(parts) != 2 {
		return nil, fmt.Errorf("mmdd:%sの形が不正です", mmdd)
	}
	month, _ := strconv.Atoi(parts[0])
	day, _ := strconv.Atoi(parts[1])

	year := time.Now().Year()

	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	return &t, nil
}
