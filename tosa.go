package tosa

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/skratchdot/open-golang/open"
)

const version = "0.0.2"

const (
	// EnvDebug is environmental var to handle debug mode
	EnvDebug = "TOSA_DEBUG"
)

// Debugf prints debug output when EnvDebug is given
func Debugf(format string, args ...interface{}) {
	if env := os.Getenv(EnvDebug); len(env) != 0 {
		log.Printf("[DEBUG] "+format+"\n", args...)
	}
}

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

func (t *Tosa) Run(args []string) (err error) {
	var debug bool
	flags := flag.NewFlagSet(args[0], flag.ContinueOnError)
	flags.BoolVar(&debug, "debug", false, "")

	// Parse flag
	if err := flags.Parse(args[1:]); err != nil {
		return err
	}

	if debug {
		os.Setenv(EnvDebug, "1")
		Debugf("Run as DEBUG mode")
	}

	sha := flags.Args()[0]
	if debug {
		Debugf("sha: %s", sha)
	}

	if err := tosa(sha); err != nil {
		return err
	}

	return nil
}

func tosa(sha string) error {
	client, err := NewClient()
	if err != nil {
		return err
	}

	if err := openPr(client, sha); err != nil {
		return err
	}

	return nil
}

func openPr(client *APIClient, sha string) error {
	pr, err := client.PullRequest(sha)
	if err != nil {
		return err
	}

	return open.Run(*pr.HTMLURL)
}
