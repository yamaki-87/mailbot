package mail

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"
	"time"

	"github.com/jordan-wright/email"
	"github.com/yamaki-87/mailbot/src/config"
	"github.com/yamaki-87/mailbot/src/consts"
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
	PAIRCOUNT       = 2
	DEFAULT_REASON  = "私用のため"
	DEFALUT_HALF    = "全休"
)

func SendMail(mail *mailtmpl.Mail) error {
	pass := os.Getenv("GMAIL_PASS")

	auth := smtp.PlainAuth("", mail.GetFrom(), pass, GMAIL_SMTP)
	return smtp.SendMail(GMAIL_SMTP_PORT, auth, mail.GetFrom(), []string{mail.GetTo()}, []byte(mail.CreateMail()))
}

func SendMailWithAttachments(mail *mailtmpl.Mail) error {
	pass := os.Getenv("GMAIL_PASS")
	email := email.NewEmail()
	email.From = mail.GetFrom()
	email.To = []string{mail.GetTo()}
	email.Subject = mail.GetSubject()
	email.Text = []byte(mail.GetBody())
	for _, v := range mail.GetAttachments() {
		if _, err := email.AttachFile(v); err != nil {
			return fmt.Errorf("failed to attach files: %w", err)
		}
	}

	auth := smtp.PlainAuth("", mail.GetFrom(), pass, GMAIL_SMTP)
	return email.Send(GMAIL_SMTP_PORT, auth)
}

type BindFunc func(req *domain.MailSendType) map[string]string

var bindFuncMap = map[domain.MailRequestType]BindFunc{
	domain.PaidLeave: func(req *domain.MailSendType) map[string]string {
		name := os.Getenv("NAME")
		mmdd := req.Args.Date.Format(MMDD_LAYOUT)
		week := utils.GetWeekDayFromDate(mmdd)

		var paidLeave = ""
		if utils.IsStrEmpty(req.Args.Half) {
			paidLeave = consts.PAIDLEAVE
		} else {
			paidLeave = consts.HALF_PAIDLEAVE
		}

		return map[string]string{
			"NAME":      name,
			"DATE":      fmt.Sprintf("%s(%s)", mmdd, week),
			"REASON":    req.Args.Reason,
			"HALF":      req.Args.Half,
			"PAIDLEAVE": paidLeave,
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
	case strings.HasPrefix(input, consts.PAIDLEAVECOMMAND):
		return parsePaidLeaveCommand(input)
	case strings.HasPrefix(input, consts.LATECOMMAND):
		mail = domain.MailSendType{Type: domain.LateArrival}

	case strings.HasPrefix(input, consts.ABSENTCOMMAND):
		mail = domain.MailSendType{Type: domain.Absence}

	default:
		return &domain.MailSendType{Type: domain.Unknown}, nil
	}

	return &mail, nil
}
func parseKeyValueArgs(args []string) map[string]string {
	result := make(map[string]string)

	for _, arg := range args {
		parts := strings.SplitN(arg, ":", PAIRCOUNT)
		if len(parts) == PAIRCOUNT {
			key := strings.ToLower(parts[0])
			value := parts[1]
			result[key] = value
		}
	}

	return result
}

var halfs = []string{"午後", "午前"}

func parsePaidLeaveCommand(input string) (*domain.MailSendType, error) {
	fields := strings.Fields(input)
	args := fields[1:] //!有給を除く

	dataMap := parseKeyValueArgs(args)

	// 日付処理
	date := dataMap["date"]
	if utils.IsStrEmpty(date) {
		return nil, fmt.Errorf("有給日付が指定されていません")
	}
	t, err := time.Parse(MMDD_LAYOUT, strings.TrimSpace(date))
	if err != nil {
		return nil, fmt.Errorf("日付の形式が不正です: %v", err)
	}

	reason := dataMap["reason"]
	if utils.IsStrEmpty(reason) {
		reason = DEFAULT_REASON
	}

	half := dataMap["half"]

	if utils.IsNotStrEmpty(half) && !utils.Contains(halfs, half) {
		return nil, fmt.Errorf("半休の形式が不正です。 午後、午前どちらかを指定してください")
	}

	mail := domain.MailSendType{
		Type: domain.PaidLeave,
		Args: domain.MailArgs{Date: t, Reason: reason, Half: half},
	}

	return &mail, nil
}
