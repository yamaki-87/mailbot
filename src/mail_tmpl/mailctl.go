package mailtmpl

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"
)

const (
	SEPERATE_SUBJECT_BODY = `/\-*-/\`
	SEPERATE_COUNT        = 2
	SEPERATE_SUBJECT      = "件名：\n"
	SEPERATE_BODY         = "本文：\n"
)

type Mail struct {
	from    string
	to      string
	subject string
	body    string
}

func NewMail(subject, body string) *Mail {
	return &Mail{
		from:    os.Getenv("GMAIL_USER"),
		to:      os.Getenv("GMAIL_TO"),
		subject: subject,
		body:    body,
	}
}
func (m *Mail) GetSubject() string {
	return m.subject
}

func (m *Mail) GetBody() string {
	return m.body
}

func (m *Mail) GetFrom() string {
	return m.from
}

func (m *Mail) GetTo() string {
	return m.to
}

func (m *Mail) String() string {
	return fmt.Sprintf("subject:%s \nbody:%s", m.subject, m.body)
}

func (m *Mail) CreateMail() string {
	var b strings.Builder
	b.WriteString("From: ")
	b.WriteString(m.from)
	b.WriteString("\n")
	b.WriteString("To: ")
	b.WriteString(m.to)
	b.WriteString("\n")
	b.WriteString("Subject: ")
	b.WriteString(m.subject)
	b.WriteString(m.body)

	return b.String()
}

func CreateMailTmpl(bind map[string]string, tmplFile string) (*Mail, error) {
	tmpl, err := template.ParseFiles(tmplFile)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, bind)
	if err != nil {
		return nil, err
	}

	mail, err := seperateBodySujbect(buf.String())
	if err != nil {
		return nil, err
	}

	return mail, err
}

func seperateBodySujbect(tmpl string) (*Mail, error) {
	parts := strings.SplitN(tmpl, SEPERATE_SUBJECT_BODY, SEPERATE_COUNT)
	if len(parts) != SEPERATE_COUNT {
		return nil, fmt.Errorf("TMPLの形が不正です。")
	}
	return NewMail(strings.Replace(parts[0], SEPERATE_SUBJECT, "", 1), strings.Replace(parts[1], SEPERATE_BODY, "", 1)), nil
}
