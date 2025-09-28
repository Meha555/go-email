# Go Email Library

[中文版](README_zh.md)

[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/Meha555/go-email?tab=doc)

A lightweight Go library for sending emails with support for plain text emails and attachments.

## Features

- No external dependencies, only wraps the standard library
- Build email content (sender, recipients, CC, BCC, subject, body)
- Support for adding attachments
- TLS encrypted connections
- Bulk email sending capability

## Installation

```bash
go get github.com/Meha555/go-email
```

## Quick Start

### Basic Usage

```go
package main

import (
    "log"
    "github.com/Meha555/go-email"
)

func main() {
    // Configure SMTP server information
    conf := email.ServerConf{
        SmtpServer: "smtp.example.com",
        SmtpPort:   587,
        Username:   "your-username",
        Password:   "your-password",
    }
    
    // Create email builder
    eb := email.NewBuilder()
    
    // Build email
    e := eb.
        From("sender@example.com").
        To("recipient@example.com").
        Cc("cc@example.com").
        Bcc("bcc@example.com").
        Subject("Hello World").
        Body("This is a test email.").
        Attachment(email.Attachment{
            Name:        "document.pdf",
            ContentType: "application/pdf",
            WithFile:    true,
        }).
        Build()
    
    // Send email
    err := email.Send(conf, e)
    if err != nil {
        log.Fatal(err)
    }
}
```

### Reusing Connections with Sender

For scenarios requiring sending multiple emails, you can reuse SMTP and TLS connections:

```go
sender := email.NewSender(conf)
if err := sender.Connect(); err != nil {
    log.Fatal(err)
}
defer sender.Disconnect()

// Send multiple emails
for i := 0; i < 5; i++ {
    e := eb.Subject(fmt.Sprintf("Email #%d", i)).Build()
    err := sender.Send(e)
    if err != nil {
        log.Printf("Failed to send email #%d: %v", i, err)
    }
}
```

## Running Examples

The project includes a complete example program that can be run with:

```bash
./run_example.sh <smtpServer> <smtpPort> <userName> <password> [fromEmail] [toEmail] [ccEmail] [bccEmail]
```

For example:
```bash
./run_example.sh smtp.gmail.com 587 your-email@gmail.com your-password
```

## API Documentation

### ServerConf
SMTP server configuration:
- `SmtpServer`: SMTP server address
- `SmtpPort`: SMTP server port
- `Username`: Username
- `Password`: Password

### Email Builder
The email builder provides chainable methods:
- `From(addr string)`: Set sender address
- `To(addr ...string)`: Set recipient address list
- `Cc(addr ...string)`: Set CC address list
- `Bcc(addr ...string)`: Set BCC address list
- `Subject(subject string)`: Set email subject
- `Body(body string)`: Set email body
- `Attachment(attachment Attachment)`: Add attachment
- `Build()`: Build the final email object

### Sender
Used for reusing SMTP connections:
- `NewSender(conf ServerConf)`: Create a new Sender instance
- `Connect()`: Establish connection to SMTP server
- `Disconnect()`: Disconnect from SMTP server
- `Send(email *Email)`: Send email through established connection