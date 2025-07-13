package domain

import (
	"time"
)

type MailRequestType int

type MailSendType struct {
	Type MailRequestType
	Args MailArgs
}

type MailArgs struct {
	Date   time.Time
	Reason string
	Half   string
}

const (
	SEPERATE        = " "
	YYYYMMDD_LAYOUT = "2025/01/02"
	MMDD_LAYOUT     = "1/2"
)

const (
	Unknown = iota
	// 有給
	PaidLeave
	// 遅延
	LateArrival
	// 欠勤
	Absence
)
