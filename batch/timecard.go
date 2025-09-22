package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"github.com/yamaki-87/mailbot/src/config"
	"github.com/yamaki-87/mailbot/src/mail"
	mailtmpl "github.com/yamaki-87/mailbot/src/mail_tmpl"
)

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func main() {
	log.Info().Msg("TimeCardBatch process started")
	config := config.GetConfig()

	timecardPath := config.TimeCard.Path
	if timecardPath == "" || !fileExists(timecardPath) {
		log.Error().Msgf("TimeCard path is not set or does not exist: %s", timecardPath)
		return
	}

	bind := createBind()
	tmplPath := config.MailTmpl.TimeCard
	if tmplPath == "" || !fileExists(tmplPath) {
		log.Error().Msgf("TimeCard template path is not set or does not exist: %s", tmplPath)
		return
	}

	// 添付ファイルの準備
	attachments := make(map[string]string)
	attachments[filepath.Base(timecardPath)] = timecardPath

	mailS, err := mailtmpl.CreateMailTmplWithAttachments(bind, tmplPath, os.Getenv("GMAIL_USER"), os.Getenv("MAIL_TIMECARD_TO"), attachments)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create TimeCard mail template")
		return
	}

	log.Info().Msgf("%+v", mailS)

	err = mail.SendMailWithAttachments(mailS)
	if err != nil {
		log.Error().Err(err).Msg("Failed to send TimeCard mail")
		return
	}

	log.Info().Msg("TimeCard mail sent successfully")
}
func createBind() map[string]string {
	t := time.Now().Local().AddDate(0, -1, 0)
	bind := make(map[string]string)
	bind["MONTH"] = t.Format("1月")
	bind["NAME"] = os.Getenv("NAME")
	return bind
}

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf(".envファイルの読み込みに失敗しました：%v", err)
	}
}
