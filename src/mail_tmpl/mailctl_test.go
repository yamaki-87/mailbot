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
		"NAME":      "test alice",
		"DATE":      "12/31(日)",
		"REASON":    "私用のため",
		"HALF":      "午前",
		"PAIDLEAVE": "半休",
	}

	tmplPath := "../../tmpl/paidLeave.txt"
	mail, err := CreateMailTmpl(bind, tmplPath)

	assert.Nil(t, err)
	assert.Equal(t, "【勤怠連絡】12/31(日) 半休 test alice", strings.TrimSpace(mail.subject))
	expected, _ := os.ReadFile("../../test_tmpl/paidLeaveMaxArgs.txt")
	assert.Equal(t, string(expected), strings.TrimSpace(mail.body))
}

func TestSuccessPaidLeaveCreateMailTmplWithNoHalf(t *testing.T) {
	test_init()

	bind := map[string]string{
		"NAME":      "test alice",
		"DATE":      "1/1(日)",
		"REASON":    "通院のため",
		"HALF":      "",
		"PAIDLEAVE": "全休",
	}

	tmplPath := "../../tmpl/paidLeave.txt"
	mail, err := CreateMailTmpl(bind, tmplPath)

	assert.Nil(t, err)
	assert.Equal(t, "【勤怠連絡】1/1(日) 全休 test alice", strings.TrimSpace(mail.subject))
	expected, _ := os.ReadFile("../../test_tmpl/paidLeaveNoHalf.txt")
	assert.Equal(t, string(expected), strings.TrimSpace(mail.body))
}
