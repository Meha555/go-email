package email

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"net/smtp"
	"net/textproto"
	"os"
	"strings"
)

type ServerConf struct {
	SmtpServer string
	SmtpPort   int
	Username   string
	Password   string
}

// Send is a convenience method to send email. SMTP server's info should specified in s.
func Send(s ServerConf, email *Email) error {
	if err := validateEmail(email); err != nil {
		return fmt.Errorf("validate email error: %w", err)
	}
	auth := smtp.PlainAuth("", s.Username, s.Password, s.SmtpServer)
	buffer := bytes.NewBuffer(nil)
	if err := assambleMail(buffer, email); err != nil {
		return fmt.Errorf("assamble email error: %w", err)
	}
	// 这里的必须是不带<>的地址，加了<>的地址会被忽略
	err := smtp.SendMail(fmt.Sprintf("%s:%d", s.SmtpServer, s.SmtpPort), auth, email.from, email.AllRecipients(), buffer.Bytes())
	if err != nil {
		return fmt.Errorf("send email error: %w", err)
	}
	return nil
}

type Sender struct {
	ServerConf

	addr   string
	client *smtp.Client
	auth   smtp.Auth
	ready  bool
}

func NewSender(conf ServerConf) *Sender {
	return &Sender{ServerConf: conf}
}

func (s *Sender) Connect() (err error) {
	if s.ready {
		return
	}
	if s.addr == "" {
		s.addr = fmt.Sprintf("%s:%d", s.SmtpServer, s.SmtpPort)
	}
	s.client, err = smtp.Dial(s.addr)
	if err != nil {
		return fmt.Errorf("dial smtp server error: %w", err)
	}
	// (smtp.Client).Extension已经进行了hello操作，所以不需要先显式调用了
	// hostname, err := os.Hostname()
	// if err != nil {
	// 	return fmt.Errorf("get hostname error: %w", err)
	// }
	// if err = s.client.Hello(hostname); err != nil {
	// 	return fmt.Errorf("hello smtp server error: %w", err)
	// }
	tlsConf := &tls.Config{
		ServerName: s.SmtpServer,
	}
	if ok, _ := s.client.Extension("STARTTLS"); ok {
		if err = s.client.StartTLS(tlsConf); err != nil {
			return fmt.Errorf("start tls error: %w", err)
		}
	}
	if ok, _ := s.client.Extension("AUTH"); ok {
		if err = s.doAuth(); err != nil {
			return fmt.Errorf("auth error: %w", err)
		}
	}
	s.ready = true
	return nil
}

func (s *Sender) Disconnect() error {
	defer func() {
		s.ready = false
	}()
	if s.client == nil {
		return nil
	}
	// 忽略错误，因为可能Connect失败了
	_ = s.client.Quit()
	return s.client.Close()
}

func (s *Sender) doAuth() error {
	// 这里的必须是不带<>的地址，加了<>的地址会被忽略
	s.auth = smtp.PlainAuth("", s.Username, s.Password, s.SmtpServer)
	return s.client.Auth(s.auth)
}

func (s *Sender) Send(email *Email) error {
	if !s.ready {
		return fmt.Errorf("sender is not connected, please call Connect() first")
	}
	if err := validateEmail(email); err != nil {
		return fmt.Errorf("validate email error: %w", err)
	}

	if err := s.client.Mail(email.from); err != nil {
		return fmt.Errorf("MAIL error: %w", err)
	}
	for _, rcpt := range email.AllRecipients() {
		if err := s.client.Rcpt(rcpt); err != nil {
			return fmt.Errorf("RCPT error: %w", err)
		}
	}
	wc, err := s.client.Data()
	if err != nil {
		return fmt.Errorf("DATA error: %w", err)
	}

	buffer := bytes.NewBuffer(nil)
	if err := assambleMail(buffer, email); err != nil {
		return fmt.Errorf("assamble email error: %w", err)
	}

	if _, err := wc.Write(buffer.Bytes()); err != nil {
		return fmt.Errorf("write email error: %w", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("close email error: %w", err)
	}
	return nil
}

// validateLine checks to see if a line has CR or LF as per RFC 5321.
func validateLine(line string) error {
	if strings.ContainsAny(line, "\n\r") {
		return errors.New("smtp: A line must not contain CR or LF")
	}
	return nil
}

func validateEmail(email *Email) error {
	if email == nil {
		return errors.New("no email set")
	}
	if err := validateLine(email.from); err != nil {
		return fmt.Errorf("from address error: %w", err)
	}
	for _, recp := range email.to {
		if err := validateLine(recp); err != nil {
			return fmt.Errorf("to address error: %w", err)
		}
	}
	for _, recp := range email.cc {
		if err := validateLine(recp); err != nil {
			return fmt.Errorf("cc address error: %w", err)
		}
	}
	for _, recp := range email.bcc {
		if err := validateLine(recp); err != nil {
			return fmt.Errorf("bcc address error: %w", err)
		}
	}
	return nil
}

const (
	CRLF     = "\r\n"
	boundary = "----GoEmailBoundary7MA4YWxkTrZu0gW"
)

func assambleMail(buffer *bytes.Buffer, email *Email) error {
	var builder strings.Builder
	header := make(textproto.MIMEHeader)
	header.Set("From", email.from)
	header.Add("To", email.Recipients().String())
	header.Add("CC", email.CcRecipients().String())
	header.Set("Subject", email.subject)
	header.Set("Content-Type", "multipart/mixed;boundary="+boundary)
	header.Set("Mime-Version", "1.0")
	writeHeader(buffer, header)

	builder.WriteString("\r\n--" + boundary + "\r\n")
	// builder.WriteString("Content-Type:" + email.contentType + "\r\n")
	builder.WriteString("Content-Type:" + "text/plain;charset=utf-8" + "\r\n")
	builder.WriteString("\r\n" + email.body + "\r\n")
	buffer.WriteString(builder.String())
	// buffer.WriteString(email.header)

	if email.attachment.WithFile {
		builder.Reset()
		builder.WriteString("\r\n--" + boundary + "\r\n")
		builder.WriteString("Content-Transfer-Encoding:base64\r\n")
		builder.WriteString("Content-Disposition:attachment\r\n")
		builder.WriteString("Content-Type:" + email.attachment.ContentType + ";name=\"" + email.attachment.Name + "\"\r\n")
		buffer.WriteString(builder.String())
		if err := writeFile(buffer, email.attachment.Name); err != nil {
			return fmt.Errorf("write file error: %w", err)
		}
	}

	buffer.WriteString("\r\n--" + boundary + "--")
	return nil
}

func writeHeader(buffer *bytes.Buffer, header textproto.MIMEHeader) {
	for key, value := range header {
		buffer.WriteString(key + ":" + strings.Join(value, ";") + "\r\n")
	}
	buffer.WriteString("\r\n")
}

func writeFile(buffer *bytes.Buffer, fileName string) error {
	file, err := os.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("read file error: %w", err)
	}
	nbytes := base64.StdEncoding.EncodedLen(len(file))
	payload := make([]byte, nbytes)
	base64.StdEncoding.Encode(payload, file)
	buffer.Grow(nbytes + 2)
	buffer.WriteString("\r\n")
	// 每写入 76 个字节（MIME 标准中 Base64 编码的行长度限制），就插入一个换行符 \r\n
	nlines := nbytes / 76
	for i := 0; i < nlines; i++ {
		buffer.Write(payload[i*76 : (i+1)*76])
		buffer.WriteString("\r\n")
	}
	buffer.Write(payload[nlines*76:])
	return nil
}
