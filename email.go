package email

import "strings"

type Address = string

type AddressLists []Address

func (a AddressLists) String() Address {
	return strings.Join(a, ",")
}

type Email struct {
	from    Address
	to      AddressLists
	cc      AddressLists
	bcc     AddressLists
	subject string
	body    string
	// contentType string
	attachment Attachment

	// b      strings.Builder
	// header string
	// builed bool
}

type Attachment struct {
	Name        string
	ContentType string
	WithFile    bool
}

// func (e *Email) Reset() {
// 	e.from = ""
// 	e.to = nil
// 	e.cc = nil
// 	e.bcc = nil
// 	e.subject = ""
// 	e.body = ""
// 	e.attachment = Attachment{}
// 	e.b.Reset()
// 	e.header = ""
// 	e.builed = false
// }

// func (e *Email) Header() string {
// 	return e.header
// }

// func (e *Email) Message() string {
// 	e.b.WriteString(e.header)
// 	e.b.WriteString("\r\n")
// 	e.b.WriteString(e.body)
// 	return e.b.String()
// }

func (e *Email) CcRecipients() AddressLists {
	return e.cc
}

func (e *Email) BccRecipients() AddressLists {
	return e.bcc
}

func (e *Email) Recipients() AddressLists {
	return e.to
}

// func (e *Email) AllRecipients() AddressLists {
// 	allRecipients := append(e.to, e.cc...)
// 	allRecipients = append(allRecipients, e.bcc...)
// 	return allRecipients
// }

func (e *Email) AllRecipients() AddressLists {
	allRecipients := make(AddressLists, 0, len(e.to)+len(e.cc)+len(e.bcc))
	allRecipients = append(allRecipients, e.to...)
	allRecipients = append(allRecipients, e.cc...)
	allRecipients = append(allRecipients, e.bcc...)
	return allRecipients
}

type Builder struct {
	email *Email
}

func NewBuilder() *Builder {
	return &Builder{
		email: &Email{},
	}
}

func (e *Builder) From(addr string) *Builder {
	// if e.email.builed {
	// 	e.email.Reset()
	// }
	e.email.from = addr
	// e.email.b.WriteString("From: " + e.email.from + "\r\n")
	return e
}

func (e *Builder) To(addr ...string) *Builder {
	// if e.email.builed {
	// 	e.email.Reset()
	// }
	e.email.to = AddressLists(addr)
	// e.email.b.WriteString("To: " + e.email.to.String() + "\r\n")
	return e
}

func (e *Builder) Cc(addr ...string) *Builder {
	// if e.email.builed {
	// 	e.email.Reset()
	// }
	e.email.cc = AddressLists(addr)
	// e.email.b.WriteString("CC: " + e.email.cc.String() + "\r\n")
	return e
}

func (e *Builder) Bcc(addr ...string) *Builder {
	// if e.email.builed {
	// 	e.email.Reset()
	// }
	e.email.bcc = AddressLists(addr)
	return e
}

func (e *Builder) Subject(subject string) *Builder {
	// if e.email.builed {
	// 	e.email.Reset()
	// }
	e.email.subject = subject
	// e.email.b.WriteString("Subject: " + e.email.subject + "\r\n")
	return e
}

func (e *Builder) Body(body string) *Builder {
	// if e.email.builed {
	// 	e.email.Reset()
	// }
	e.email.body = body
	return e
}

func (e *Builder) Attachment(attachment Attachment) *Builder {
	// if e.email.builed {
	// 	e.email.Reset()
	// }
	e.email.attachment = attachment
	// if attachment.WithFile {
	// 	e.email.b.WriteString("Content-Transfer-Encoding:base64\r\n")
	// 	e.email.b.WriteString("Content-Disposition:attachment\r\n")
	// 	e.email.b.WriteString("Content-Type:" + e.email.attachment.ContentType + ";name=\"" + e.email.attachment.Name + "\"\r\n")
	// }
	return e
}

func (e *Builder) Build() *Email {
	// if e.email.builed {
	// 	return e.email
	// }
	// e.email.header = e.email.b.String()
	// e.email.b.Reset()
	// e.email.builed = true
	return e.email
}
