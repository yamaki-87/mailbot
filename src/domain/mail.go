package domain

import (
	"time"
)

type MailRequestType int

type MailSendType struct {
	Type             MailRequestType
	Args             MailArgs
	SpecialLeaveArgs SpecialLeaveMailArgs
	IsTest           bool
}

func (m *MailSendType) SetIsTest(isTest bool) {
	m.IsTest = isTest
}

type MailArgs struct {
	Dates  []time.Time
	Reason string
	Half   string
}

type SpecialLeaveMailArgs struct {
	Dates []time.Time
	// 詳細な種別 ex:夏季休暇etc
	DetailType string
	Reason     string
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
	// 特別休暇
	SpecialLeave
)
