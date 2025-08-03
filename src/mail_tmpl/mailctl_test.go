package mailtmpl

import (
	"log"
	"os"
	"strings"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func test_init() {
	// .env 読み込み（テスト開始前に一度だけ実行される）
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("⚠️ .envファイルが読み込めませんでした:", err)
	}

}

func TestSuccessPaidLeaveCreateMailTmplWithMaxArgs(t *testing.T) {
	test_init()

	bind := map[string]string{
		"NAME":         "test alice",
		"DATE_HALF":    "・12/31(日) 午前\n・1/1(月) 午前",
		"SUBJECT_DATE": "12/31(日) 1/1(月)",
		"REASON":       "私用のため",
		"HALF":         "午前",
		"PAIDLEAVE":    "半休",
	}

	tmplPath := "../../tmpl/paidLeave.txt"
	mail, err := CreateMailTmpl(bind, tmplPath)

	assert.Nil(t, err)
	assert.Equal(t, "【勤怠連絡】12/31(日) 1/1(月) 半休 test alice", strings.TrimSpace(mail.subject))
	expected, _ := os.ReadFile("../../test_tmpl/paidLeaveMaxArgs.txt")
	assert.Equal(t, string(expected), strings.TrimSpace(mail.body))
}

func TestSuccessPaidLeaveCreateMailTmplWithNoHalf(t *testing.T) {
	test_init()

	bind := map[string]string{
		"NAME":         "test alice",
		"DATE_HALF":    "・1/1(日) ",
		"SUBJECT_DATE": "1/1(日)",
		"REASON":       "通院のため",
		"HALF":         "",
		"PAIDLEAVE":    "全休",
	}

	tmplPath := "../../tmpl/paidLeave.txt"
	mail, err := CreateMailTmpl(bind, tmplPath)

	assert.Nil(t, err)
	assert.Equal(t, "【勤怠連絡】1/1(日) 全休 test alice", strings.TrimSpace(mail.subject))
	expected, _ := os.ReadFile("../../test_tmpl/paidLeaveNoHalf.txt")
	assert.Equal(t, string(expected), strings.TrimSpace(mail.body))
}

func TestSuccessSpecialLeaveCreateMailWithMaxArgs(t *testing.T) {
	test_init()

	bind := map[string]string{
		"NAME":         "test tarou",
		"DATE":         "・1/1(日) 冬期休暇\n・1/2(月) 冬期休暇",
		"SUBJECT_DATE": "1/1(日) 1/2(月)",
		"REASON":       "現場が年末のため",
		"DETAIL_TYPE":  "冬期休暇",
	}

	tmplPath := "../../tmpl/specialLeave.txt"
	mail, err := CreateMailTmpl(bind, tmplPath)
	assert.Nil(t, err)
	assert.Equal(t, "【勤怠連絡】1/1(日) 1/2(月) 冬期休暇 test tarou", strings.TrimSpace(mail.subject))
	expected, _ := os.ReadFile("../../test_tmpl/specialLeaveMaxArgs.txt")
	assert.Equal(t, string(expected), strings.TrimSpace(mail.body))
}
