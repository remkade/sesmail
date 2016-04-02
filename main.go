package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/mail"
	"os"
	"stathat.com/c/amzses"
)

type Config struct {
	Sender       string
	ExtractTo    bool
	Allow_period bool
	Debug        bool
	Version      bool
}

var config Config

func init() {
	flag.StringVar(&config.Sender, "f", "", "Set the FROM address")
	flag.BoolVar(&config.ExtractTo, "t", false, "Extract the recipients from the body headers")
	flag.BoolVar(&config.Debug, "debug", false, "Print too much output")
	flag.BoolVar(&config.Allow_period, "i", true, "Allow a single period on a line by itself without terminating output (IGNORED)")
	flag.BoolVar(&config.Version, "version", false, "Print version and exit")
	flag.Parse()
}

func main() {
	version := "0.1.0"
	var to string
	var sender string

	if config.Version == true {
		fmt.Println(version)
		os.Exit(0)
	}

	m, err := mail.ReadMessage(os.Stdin)
	if err != nil {
		fmt.Printf("Error reading from stdin '%s'\n", err)
		os.Exit(1)
	}

	to = m.Header.Get("To")

	if config.Sender == "" {
		sender = m.Header.Get("From")
	} else {
		sender = config.Sender
	}

	body, err := ioutil.ReadAll(m.Body)
	if err != nil {
		fmt.Printf("Error Reading message body: '%s'", err)
		os.Exit(1)
	}

	subject := m.Header.Get("Subject")
	fmt.Printf("%+v\n", config)
	fmt.Printf("%+v\n", m.Header)

	switch {
	case to == "":
		fmt.Printf("To is empty!\n")
		os.Exit(1)
	case sender == "":
		fmt.Printf("sender is empty!\n")
		os.Exit(1)
	case subject == "":
		fmt.Printf("Subject is empty!\n")
		os.Exit(1)
	}

	fmt.Println("Sending message")
	fmt.Printf("sender='%s',to='%s',subject='%s'\n", sender, to, subject)

	amzses.SendMail(sender, to, subject, string(body))
}
