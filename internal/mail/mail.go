package mail

import (
	"bytes"
	"fmt"
	"html/template"
)

type (
	Mail struct {
		to      []string
		from    string
		subject string
		body    interface{}
	}

	Template struct {
		Mail
		template *template.Template
	}
)

func (m *Mail) Subject() string {
	return m.subject
}

func (m *Mail) From() string {
	return m.from
}

func (m *Mail) To() []string {
	return m.to
}

func (m *Mail) headers() *bytes.Buffer {
	buf := new(bytes.Buffer)
	buf.WriteString("Subject: ")
	buf.WriteString(m.subject)
	buf.WriteString("\n")
	buf.WriteString("MIME-version: 1.0;\n")
	buf.WriteString("Content-Type: text/html; charset=\"UTF-8\";\n")
	if m.from != "" {
		buf.WriteString("From: " + m.from + "\n")
		buf.WriteString(m.from)
		buf.WriteString("\n")
	}
	buf.WriteString("\n\n")
	return buf
}

func (m *Mail) Message() ([]byte, error) {
	buf := m.headers()
	buf.WriteString(fmt.Sprintf("%v", m.body))
	return buf.Bytes(), nil
}

func (t *Template) Message() ([]byte, error) {
	buf := t.headers()
	if err := t.template.ExecuteTemplate(buf, "base", t.body); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func NewTemplate(from, subject string, data interface{}, tpl *template.Template, to ...string) Mailer {
	return &Template{
		Mail: Mail{
			subject: subject,
			from:    from,
			body:    data,
			to:      to,
		},
		template: tpl,
	}
}
