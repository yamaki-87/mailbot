package mail

import (
	"log"
	"testing"

	"github.com/joho/godotenv"
)

func TestSuccsessSendMail(t *testing.T) {
	// .env 読み込み（テスト開始前に一度だけ実行される）
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("⚠️ .envファイルが読み込めませんでした:", err)
	}
	//err := sendMail("Test", "テストです")
	// if err != nil {
	// 	t.Errorf("%v", err)
	// }
}
