package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSucGetWeekDay(t *testing.T) {
	weekDay := GetWeekDayFromDate("6/14")

	assert := assert.New(t)
	assert.Equal("土", weekDay, "2025 06.14 is 土")
}
