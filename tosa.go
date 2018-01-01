package tosa

import (
	"flag"
	"fmt"
	"os"

	"github.com/skratchdot/open-golang/open"
)

const version = "0.0.1"

type ErrorMessage struct {
	message string
}

func NewErrorMessage(message string) *ErrorMessage {
	var m string
	if message == "" {
		m = ""
	} else {
		m = fmt.Sprintf("%s\n%s", "", message)
	}
	return &ErrorMessage{
		message: m,
	}
}

func (e *ErrorMessage) Error() string {
	return e.message
}

func New() *Tosa {
	return &Tosa{
		Argv:   os.Args,
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
	}
}

func (o *Tosa) Err() error {
	return o.err
}

func (o *Tosa) Run() (err error) {
	if err := tosa(); err != nil {
		return err
	}
	return nil
}

func tosa() error {
	flag.Parse()
	sha := flag.Args()[0]
	client, err := NewClient()
	if err != nil {
		return err
	}

	openPr(client, sha)
	return nil
}

func openPr(client *APIClient, sha string) error {
	pr, err := client.PullRequest(sha)
	if err != nil {
		return err
	}

	return open.Run(*pr.HTMLURL)
}
