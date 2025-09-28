# Go Email Library

一个轻量的Go语言邮件发送库，支持普通文本邮件、附件发送。

## 功能特性

- 无外部依赖，仅封装标准库
- 构建邮件内容（发件人、收件人、抄送、密送、主题、正文）
- 支持添加附件
- 支持TLS加密连接
- 支持批量发送邮件

## 安装

```bash
go get github.com/Meha555/go-email
```

## 快速开始

### 基本用法

```go
package main

import (
    "log"
    "github.com/Meha555/go-email"
)

func main() {
    // 配置SMTP服务器信息
    conf := email.ServerConf{
        SmtpServer: "smtp.example.com",
        SmtpPort:   587,
        Username:   "your-username",
        Password:   "your-password",
    }
    
    // 创建邮件构建器
    eb := email.NewBuilder()
    
    // 构建邮件
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
    
    // 发送邮件
    err := email.Send(conf, e)
    if err != nil {
        log.Fatal(err)
    }
}
```

### 使用Sender进行连接复用

对于需要发送多封邮件的场景，可以复用SMTP和TLS连接：

```go
sender := email.NewSender(conf)
if err := sender.Connect(); err != nil {
    log.Fatal(err)
}
defer sender.Disconnect()

// 发送多封邮件
for i := 0; i < 5; i++ {
    e := eb.Subject(fmt.Sprintf("Email #%d", i)).Build()
    err := sender.Send(e)
    if err != nil {
        log.Printf("Failed to send email #%d: %v", i, err)
    }
}
```

## 运行示例

项目包含一个完整的示例程序，可以通过以下方式运行：

```bash
./run_example.sh <smtpServer> <smtpPort> <userName> <password> [fromEmail] [toEmail] [ccEmail] [bccEmail]
```

例如：
```bash
./run_example.sh smtp.gmail.com 587 your-email@gmail.com your-password
```

## API文档

### ServerConf
SMTP服务器配置信息：
- `SmtpServer`: SMTP服务器地址
- `SmtpPort`: SMTP服务器端口
- `Username`: 用户名
- `Password`: 密码

### Email Builder
邮件构建器提供了链式调用的方法：
- `From(addr string)`: 设置发件人地址
- `To(addr ...string)`: 设置收件人地址列表
- `Cc(addr ...string)`: 设置抄送地址列表
- `Bcc(addr ...string)`: 设置密送地址列表
- `Subject(subject string)`: 设置邮件主题
- `Body(body string)`: 设置邮件正文
- `Attachment(attachment Attachment)`: 添加附件
- `Build()`: 构建最终的邮件对象

### Sender
用于复用SMTP连接：
- `NewSender(conf ServerConf)`: 创建新的Sender实例
- `Connect()`: 建立到SMTP服务器的连接
- `Disconnect()`: 断开与SMTP服务器的连接
- `Send(email *Email)`: 通过已建立的连接发送邮件
