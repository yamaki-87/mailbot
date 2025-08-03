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

	input := "!有給 date:7/15,7/16 reason:通院のため half:午後 -t"
	mail, err := ParseMailSendType(input)
	assert.Nil(t, err)

	assert.Equal(t, domain.PaidLeave, int(mail.Type), "有給になっているかenumなのでintの値が同じかどうか")
	assert.Equal(t, "7/15", mail.Args.Dates[0].Format(MMDD_LAYOUT), "7/15ではないとおかしい")
	assert.Equal(t, "7/16", mail.Args.Dates[1].Format(MMDD_LAYOUT), "7/15ではないとおかしい")
	assert.Equal(t, "通院のため", mail.Args.Reason, "同じ値になっているはず")
	assert.Equal(t, "午後", mail.Args.Half, "同じ値になっているはず")
	assert.True(t, mail.IsTest, "-tを指定したためTrueになってるはず")
}

func TestParsePaidLeaveCommandWithDate(t *testing.T) {
	test_init()

	input := "!有給 date:12/31"
	mail, err := parsePaidLeaveCommand(input)
	assert.Nil(t, err)

	assert.Equal(t, domain.PaidLeave, int(mail.Type), "有給になっているかenumなのでintの値が同じかどうか")
	assert.Equal(t, "12/31", mail.Args.Dates[0].Format(MMDD_LAYOUT), "7/15ではないとおかしい")
	assert.Equal(t, DEFAULT_REASON, mail.Args.Reason, "同じ値になっているはず")
	assert.Equal(t, "", mail.Args.Half, "空文字のはず")
}

func TestParsePaidLeaveCommandWithDateAndReason(t *testing.T) {
	test_init()

	input := "!有給 date:1/1 reason:代休のため"
	mail, err := parsePaidLeaveCommand(input)
	assert.Nil(t, err)

	assert.Equal(t, domain.PaidLeave, int(mail.Type), "有給になっているかenumなのでintの値が同じかどうか")
	assert.Equal(t, "1/1", mail.Args.Dates[0].Format(MMDD_LAYOUT), "7/15ではないとおかしい")
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
	assert.Equal(t, "1/1", mail.Args.Dates[0].Format(MMDD_LAYOUT), "7/15ではないとおかしい")
	assert.Equal(t, DEFAULT_REASON, mail.Args.Reason, "reasonTest:代休のため:テストのためは訳されずDEFAULT REASONが使われる")
	assert.Equal(t, "午前", mail.Args.Half, "halfには午前が入ってるはず")
}

func TestParseSpecialLeaveWithMaxArgs(t *testing.T) {
	test_init()

	input := "!特別休暇 date:1/1,12/21 reason:現場が年末のため type:冬期休暇"
	mail, err := parseSpecialLeaveCommand(input)
	assert.Nil(t, err)

	assert.Equal(t, domain.SpecialLeave, int(mail.Type), "特別休暇になっているかenumなのでintの値が同じかどうかあ")
	assert.Equal(t, "1/1", mail.SpecialLeaveArgs.Dates[0].Format(MMDD_LAYOUT), "1/1のはず")
	assert.Equal(t, "12/21", mail.SpecialLeaveArgs.Dates[1].Format(MMDD_LAYOUT), "12/21のはず")
	assert.Equal(t, 2, len(mail.SpecialLeaveArgs.Dates), "2個しないはず")
	assert.Equal(t, "現場が年末のため", mail.SpecialLeaveArgs.Reason, "現場が年末のためという値が入ってるはず")
	assert.Equal(t, "冬期休暇", mail.SpecialLeaveArgs.DetailType, "冬期休暇が入ってるはず")
}
func TestParseSpecialLeaveWithNoReason(t *testing.T) {
	test_init()

	input := "!特別休暇 date:1/1 type:夏季休暇 "
	mail, err := parseSpecialLeaveCommand(input)

	assert.Nil(t, err)

	assert.Equal(t, domain.SpecialLeave, int(mail.Type), "特別休暇になっているかenumなのでintの値が同じ")
	assert.Equal(t, "1/1", mail.SpecialLeaveArgs.Dates[0].Format(MMDD_LAYOUT), "1/1のはず")
	assert.Equal(t, 1, len(mail.SpecialLeaveArgs.Dates), "2個しないはず")
	assert.Equal(t, DEFAULT_REASON, mail.SpecialLeaveArgs.Reason, "私用のため")
	assert.Equal(t, "夏季休暇", mail.SpecialLeaveArgs.DetailType, "夏季休暇入ってるはず")
}

func TestFailParseSpecialLeaveWithNoDateArgs(t *testing.T) {
	test_init()

	input := "!特別休暇"
	mail, err := parseSpecialLeaveCommand(input)

	assert.Nil(t, mail, "nullのはず")
	assert.EqualError(t, err, "特別休暇日付が指定されていません")
}

func TestFailParseSpecialLeaveWithNoDetailTypeArgs(t *testing.T) {
	test_init()

	input := "!特別休暇 date:1/1"
	mail, err := parseSpecialLeaveCommand(input)

	assert.Nil(t, mail, "nullのはず")
	assert.EqualError(t, err, "詳細な休暇種別を指定してください ex:夏季休暇etc")
}
func TestParseLateCommandWithIsTestTrue(t *testing.T) {
	test_init()

	input := "!遅延 -t"
	mail, err := ParseMailSendType(input)

	assert.Nil(t, err)

	assert.Equal(t, domain.LateArrival, int(mail.Type))
	assert.True(t, mail.IsTest, "-tを指定したためTrueになってるはず")
}
