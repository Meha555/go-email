package main

import (
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/Meha555/go-email"
)

var (
	smtpServer string
	smtpPort   string
	userName   string
	password   string

	fromEmail string
	toEmail   string
	ccEmail   string
	bccEmail  string
)

func printArgs() {
	fmt.Printf("smtpServer: %s\n", smtpServer)
	fmt.Printf("smtpPort: %s\n", smtpPort)
	fmt.Printf("userName: %s\n", userName)
	fmt.Printf("password: %s\n", password)
	fmt.Printf("fromEmail: %s\n", fromEmail)
	fmt.Printf("toEmail: %s\n", toEmail)
	fmt.Printf("ccEmail: %s\n", ccEmail)
	fmt.Printf("bccEmail: %s\n", bccEmail)
}

func main() {
	printArgs()
	port, _ := strconv.Atoi(smtpPort)
	conf := email.ServerConf{
		SmtpServer: smtpServer,
		SmtpPort:   port,
		Username:   userName,
		Password:   password,
	}
	eb := email.NewBuilder()
	e := eb.
		From(fromEmail).
		To(toEmail).
		Cc(ccEmail).
		Bcc(bccEmail).
		Subject("subject").
		Body("body").
		Attachment(email.Attachment{
			Name:        "example/example.jpeg",
			ContentType: "image/jpeg",
			WithFile:    true,
		}).
		Build()

	log.Println(email.Send(conf, e))

	sender := email.NewSender(conf)
	if err := sender.Connect(); err != nil {
		log.Fatal(err)
	}
	defer sender.Disconnect()
	wg := sync.WaitGroup{}
	for i := 0; i < 5; i++ {
		// 可能频率过快被服务器拒绝
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_ = eb.Body("body@" + strconv.Itoa(i))
			log.Println(sender.Send(e))
		}(i)
	}
	wg.Wait()
}
