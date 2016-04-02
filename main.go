package main

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	ses "github.com/tamcgoey/go-ses"
	"io/ioutil"
	"net/mail"
	"os"
)

type CliConfig struct {
	Sender      string
	ExtractTo   bool
	AllowPeriod bool
	Debug       bool
	Version     bool
	Config      string
}

type TOMLConfig struct {
	AccessKey string
	SecretKey string
	Region    string
}

func (t *TOMLConfig) Endpoint() string {
	return fmt.Sprintf("https://email.%s.amazonaws.com", t.Region)
}

var config CliConfig

func init() {
	flag.StringVar(&config.Sender, "f", "", "Set the FROM address")
	flag.StringVar(&config.Config, "config", "/etc/sesmail/sesmail.toml", "Change the config file")
	flag.BoolVar(&config.ExtractTo, "t", false, "Extract the recipients from the body headers")
	flag.BoolVar(&config.Debug, "debug", false, "Print too much output")
	flag.BoolVar(&config.AllowPeriod, "i", true, "Allow a single period on a line by itself without terminating output (IGNORED)")
	flag.BoolVar(&config.Version, "version", false, "Print version and exit")
	flag.Parse()
}

func main() {
	version := "0.2.0"
	var to string
	var sender string

	var tomlConfig TOMLConfig
	if _, err := toml.DecodeFile(config.Config, &tomlConfig); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

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
	reply_to := m.Header.Get("Reply-To")
	if reply_to == "" {
		reply_to = sender
	}

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

	s := ses.Config{
		Endpoint:        tomlConfig.Endpoint(),
		AccessKeyID:     tomlConfig.AccessKey,
		SecretAccessKey: tomlConfig.SecretKey,
	}
	s.SendEmail(sender, reply_to, to, subject, string(body))
}
