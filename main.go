package main

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	log "github.com/Sirupsen/logrus"
	logrus_syslog "github.com/Sirupsen/logrus/hooks/syslog"
	ses "github.com/tamcgoey/go-ses"
	"io/ioutil"
	"log/syslog"
	"net/mail"
	"os"
)

type CliConfig struct {
	Sender      string
	ExtractTo   bool
	AllowPeriod bool
	Debug       bool
	Syslog      bool
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
	flag.BoolVar(&config.Syslog, "syslog", false, "Log to syslog")
	flag.BoolVar(&config.Version, "version", false, "Print version and exit")
	flag.Parse()
}

func main() {
	version := "0.3.0"
	codename := "Melted Chihuahas"
	var to string
	var sender string

	if config.Syslog {
		hook, err := logrus_syslog.NewSyslogHook("", "", syslog.LOG_LOCAL0, "sesmail")
		if err != nil {
			log.Fatal(err)
		}
		log.AddHook(hook)
		log.SetFormatter(&log.TextFormatter{})
	}

	var tomlConfig TOMLConfig
	if _, err := toml.DecodeFile(config.Config, &tomlConfig); err != nil {
		log.Fatal(err)
	}

	if config.Version == true {
		fmt.Printf("sesmail %s, '%s'", version, codename)
		os.Exit(0)
	}

	m, err := mail.ReadMessage(os.Stdin)
	if err != nil {
		log.Fatalf("Error reading from stdin '%s'\n", err)
	}

	to = m.Header.Get("To")

	if config.Sender == "" {
		sender = m.Header.Get("From")
	} else {
		sender = config.Sender
	}

	body, err := ioutil.ReadAll(m.Body)
	if err != nil {
		log.Fatalf("Error Reading message body: '%s'", err)
	}

	subject := m.Header.Get("Subject")
	reply_to := m.Header.Get("Reply-To")
	if reply_to == "" {
		reply_to = sender
	}

	switch {
	case to == "":
		log.Fatal("To is empty!")
	case sender == "":
		log.Fatal("sender is empty!")
	case subject == "":
		log.Fatal("Subject is empty!")
	}

	s := ses.Config{
		Endpoint:        tomlConfig.Endpoint(),
		AccessKeyID:     tomlConfig.AccessKey,
		SecretAccessKey: tomlConfig.SecretKey,
	}
	_, err = s.SendEmail(sender, reply_to, to, subject, string(body))
	if err != nil {
		log.WithFields(log.Fields{
			"to":       to,
			"sender":   sender,
			"reply_to": reply_to,
			"subject":  subject,
		}).Fatalf("Error Sending Email: '%s'", err)
	} else {
		log.WithFields(log.Fields{
			"to":       to,
			"sender":   sender,
			"reply_to": reply_to,
			"subject":  subject,
		}).Info("Message sent successfully")
	}
}
