package mail

import (
	"log"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/yamaki-87/mailbot/src/domain"
)

func test_init() {
	// .env 読み込み（テスト開始前に一度だけ実行される）
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("⚠️ .envファイルが読み込めませんでした:", err)
	}

}

func TestParsePaidLeaveCommandSuccessMaxArgs(t *testing.T) {

	test_init()

	input := "!有給 date:7/15 reason:通院のため half:午後"
	mail, err := parsePaidLeaveCommand(input)
	assert.Nil(t, err)

	assert.Equal(t, domain.PaidLeave, int(mail.Type), "有給になっているかenumなのでintの値が同じかどうか")
	assert.Equal(t, "7/15", mail.Args.Date.Format(MMDD_LAYOUT), "7/15ではないとおかしい")
	assert.Equal(t, "通院のため", mail.Args.Reason, "同じ値になっているはず")
	assert.Equal(t, "午後", mail.Args.Half, "同じ値になっているはず")
}

func TestParsePaidLeaveCommandWithDate(t *testing.T) {
	test_init()

	input := "!有給 date:12/31"
	mail, err := parsePaidLeaveCommand(input)
	assert.Nil(t, err)

	assert.Equal(t, domain.PaidLeave, int(mail.Type), "有給になっているかenumなのでintの値が同じかどうか")
	assert.Equal(t, "12/31", mail.Args.Date.Format(MMDD_LAYOUT), "7/15ではないとおかしい")
	assert.Equal(t, DEFAULT_REASON, mail.Args.Reason, "同じ値になっているはず")
	assert.Equal(t, "", mail.Args.Half, "空文字のはず")
}

func TestParsePaidLeaveCommandWithDateAndReason(t *testing.T) {
	test_init()

	input := "!有給 date:1/1 reason:代休のため"
	mail, err := parsePaidLeaveCommand(input)
	assert.Nil(t, err)

	assert.Equal(t, domain.PaidLeave, int(mail.Type), "有給になっているかenumなのでintの値が同じかどうか")
	assert.Equal(t, "1/1", mail.Args.Date.Format(MMDD_LAYOUT), "7/15ではないとおかしい")
	assert.Equal(t, "代休のため", mail.Args.Reason, "同じ値になっているはず")
	assert.Equal(t, "", mail.Args.Half, "空文字のはず")
}

func TestFailParsePaidLeaveCommandWithNoArgs(t *testing.T) {
	test_init()

	input := "!有給"
	mail, err := parsePaidLeaveCommand(input)

	assert.Nil(t, mail, "nullのはず")
	assert.EqualError(t, err, "有給日付が指定されていません")
}

func TestFailParsePaidLeaveCommandWithInvalidDate(t *testing.T) {
	test_init()

	input := "!有給 11111/1111"
	mail, err := parsePaidLeaveCommand(input)

	assert.Nil(t, mail, "nullのはず")
	assert.EqualError(t, err, "有給日付が指定されていません")
}

func TestFailParsePaidLeaveCommandWithInvalidHalf(t *testing.T) {
	test_init()

	input := "!有給 date:11/11 half:全"
	mail, err := parsePaidLeaveCommand(input)

	assert.Nil(t, mail, "nullのはず")
	assert.EqualError(t, err, "半休の形式が不正です。 午後、午前どちらかを指定してください")
}

func TestParsePaidLeaveCommandWithInvalidKey(t *testing.T) {
	test_init()

	input := "!有給 date:1/1 reasonTest:代休のため:テストのため half:午前"
	mail, err := parsePaidLeaveCommand(input)
	assert.Nil(t, err)

	assert.Equal(t, domain.PaidLeave, int(mail.Type), "有給になっているかenumなのでintの値が同じかどうか")
	assert.Equal(t, "1/1", mail.Args.Date.Format(MMDD_LAYOUT), "7/15ではないとおかしい")
	assert.Equal(t, DEFAULT_REASON, mail.Args.Reason, "reasonTest:代休のため:テストのためは訳されずDEFAULT REASONが使われる")
	assert.Equal(t, "午前", mail.Args.Half, "halfには午前が入ってるはず")
}
