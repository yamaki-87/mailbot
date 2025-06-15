package mailtmpl

import "testing"

func TestSuccessCreateMailTmpl(t *testing.T) {
	bind := make(map[string]string)
	bind["DATE"] = ""
	bind["NAME"] = "foobar"

	mail, err := CreateMailTmpl(bind, "../../tmpl/yuukyuu.txt")
	if err != nil {
		t.Errorf("%v", err)
	}

	t.Log(mail.String())
}
