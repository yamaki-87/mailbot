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

func GetWeekDayFromDateStr(mmdd string) string {
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

// 日付配列からlayout形式の文字列に結合(複数ある場合は改行で結合)
//
// input: times=[12/21,1/1l]a,layout="12/12" -> output: "12/21(木)\n1/11(月)"
func FormatDates(times []time.Time, layout string) string {
	return FormatDatesWithSeparator(times, layout, "\n")
}

// 日付配列からlayout形式の文字列に結合(複数ある場合はセパレータで結合)
//
// input: times=[12/21,1/1l]a,layout="12/12",separator=" " -> output: "12/21(木) 1/11(月)"
func FormatDatesWithSeparator(times []time.Time, layout string, separator string) string {
	var sb strings.Builder
	for _, time := range times {
		sb.WriteString(fmt.Sprintf("%s(%s)%s", time.Format(layout), WEEKDAYJP[time.Weekday()], separator))
	}
	return sb.String()
}

func GetJpnWeek(date time.Time) string {
	return WEEKDAYJP[date.Weekday()]
}
