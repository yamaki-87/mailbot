package mail

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"mime/quotedprintable"
	"net/smtp"
	"net/textproto"
	"os"
	"strings"
	"time"

	"github.com/yamaki-87/mailbot/src/config"
	"github.com/yamaki-87/mailbot/src/domain"
	mailtmpl "github.com/yamaki-87/mailbot/src/mail_tmpl"
	"github.com/yamaki-87/mailbot/src/utils"
)

const (
	GMAIL_SMTP      = "smtp.gmail.com"
	GMAIL_SMTP_PORT = "smtp.gmail.com:587"
	SEPERATE        = " "
	YYYYMMDD_LAYOUT = "2025/01/02"
	MMDD_LAYOUT     = "1/2"
)

func SendMail(mail *mailtmpl.Mail) error {
	pass := os.Getenv("GMAIL_PASS")

	auth := smtp.PlainAuth("", mail.GetFrom(), pass, GMAIL_SMTP)
	return smtp.SendMail(GMAIL_SMTP_PORT, auth, mail.GetFrom(), []string{mail.GetTo()}, []byte(mail.CreateMail()))
}

func SendMailWithAttachments(mail *mailtmpl.Mail) error {
	pass := os.Getenv("GMAIL_PASS")

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	boundary := writer.Boundary()

	header := map[string]string{
		"From":         mail.GetFrom(),
		"To":           mail.GetTo(),
		"Subject":      mail.GetSubject(),
		"MIME-Version": "1.0",
		"Content-Type": fmt.Sprintf("multipart/mixed; boundary=%s", boundary),
	}
	for key, value := range header {
		fmt.Fprintf(&buf, "%s: %s\r\n", key, value)
	}
	fmt.Fprintf(&buf, "\r\n")

	bodyWriter, _ := writer.CreatePart(map[string][]string{
		"Content-Type":              {"text/plain; charset=UTF-8"},
		"Content-Transfer-Encoding": {"quoted-printable"},
	})
	qp := quotedprintable.NewWriter(bodyWriter)
	qp.Write([]byte(mail.GetBody()))
	qp.Close()

	for displayName, filePath := range mail.GetAttachments() {

		data, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("添付ファイルの読み込みに失敗しました: %v", err)
		}
		attachHeader := make(textproto.MIMEHeader)
		attachHeader.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", displayName))
		attachHeader.Set("Content-Type", "application/octet-stream")
		attachHeader.Set("Content-Transfer-Encoding", "base64")

		part, _ := writer.CreatePart(attachHeader)
		encoded := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
		base64.StdEncoding.Encode(encoded, data)

		for i := 0; i < len(encoded); i += 76 {
			end := i + 76
			if end > len(encoded) {
				end = len(encoded)
			}
			part.Write(encoded[i:end])
			part.Write([]byte("\r\n"))
		}
	}
	writer.Close()

	auth := smtp.PlainAuth("", mail.GetFrom(), pass, GMAIL_SMTP)
	return smtp.SendMail(GMAIL_SMTP_PORT, auth, mail.GetFrom(), []string{mail.GetTo()}, buf.Bytes())
}

type BindFunc func(req *domain.MailSendType) map[string]string

var bindFuncMap = map[domain.MailRequestType]BindFunc{
	domain.PaidLeave: func(req *domain.MailSendType) map[string]string {
		name := os.Getenv("NAME")
		mmdd := req.Args.Date.Format(MMDD_LAYOUT)
		week := utils.GetWeekDayFromDate(mmdd)
		return map[string]string{
			"NAME": name,
			"DATE": fmt.Sprintf("%s(%s)", mmdd, week),
		}
	},
	domain.LateArrival: func(req *domain.MailSendType) map[string]string {
		name := os.Getenv("NAME")
		mmdd := time.Now().Format(MMDD_LAYOUT)
		week := utils.GetWeekDayFromDate(mmdd)
		return map[string]string{
			"NAME": name,
			"DATE": fmt.Sprintf("%s(%s)", mmdd, week),
		}
	},
	domain.Absence: func(req *domain.MailSendType) map[string]string {
		return map[string]string{
			"NAME": os.Getenv("NAME"),
		}
	},
}

type TemplatePathResolver func(cfg config.MailTmplConfig) string

var templatePathMap = map[domain.MailRequestType]TemplatePathResolver{
	domain.PaidLeave:   func(cfg config.MailTmplConfig) string { return cfg.PaidLeave },
	domain.LateArrival: func(cfg config.MailTmplConfig) string { return cfg.LateArrival },
	domain.Absence:     func(cfg config.MailTmplConfig) string { return cfg.Absence },
}

type MailService struct{}

func (m *MailService) BindTemplate(req *domain.MailSendType) map[string]string {
	if bind, ok := bindFuncMap[req.Type]; ok {
		return bind(req)
	}
	return nil
}

func (m *MailService) GetTmplPath(t domain.MailRequestType) string {
	config := config.GetConfig()
	if resolver, ok := templatePathMap[t]; ok {
		return resolver(config.MailTmpl)
	}
	return ""
}

func ParseMailSendType(input string) (*domain.MailSendType, error) {
	parts := strings.Split(input, domain.SEPERATE)
	if len(parts) < 1 {
		return nil, fmt.Errorf("メッセージ形式が不正です")
	}

	var mail domain.MailSendType

	switch {
	case strings.HasPrefix(input, "!有給"):
		if len(parts) < 2 {
			return nil, fmt.Errorf("有給日付が指定されていません")
		}
		t, err := time.Parse(MMDD_LAYOUT, strings.TrimSpace(parts[1]))
		if err != nil {
			return nil, fmt.Errorf("日付の形式が不正です: %v", err)
		}
		mail = domain.MailSendType{
			Type: domain.PaidLeave,
			Args: domain.MailArgs{Date: t},
		}

	case strings.HasPrefix(input, "!遅延"):
		mail = domain.MailSendType{Type: domain.LateArrival}

	case strings.HasPrefix(input, "!欠勤"):
		mail = domain.MailSendType{Type: domain.Absence}

	default:
		return &domain.MailSendType{Type: domain.Unknown}, nil
	}

	return &mail, nil
}
