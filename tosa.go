package tosa

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/mitchellh/colorstring"
	"github.com/skratchdot/open-golang/open"
)

const version = "0.0.2"

const (
	// EnvDebug is environmental var to handle debug mode
	EnvDebug = "TOSA_DEBUG"
)

// Exit codes are in value that represnet an exit code for a paticular error.
const (
	ExitCodeOK int = 0

	// Errors start at 10
	ExitCodeError = 10 + iota
	ExitCodeParseFlagsError
	ExitCodeBadArgs
	ExitCodePullRequestNotFound
)

// Debugf prints debug output when EnvDebug is given
func Debugf(format string, args ...interface{}) {
	if env := os.Getenv(EnvDebug); len(env) != 0 {
		log.Printf("[DEBUG] "+format+"\n", args...)
	}
}

// PrintErrorf prints error message on console.
func PrintErrorf(format string, args ...interface{}) {
	format = fmt.Sprintf("[red]%s[reset]\n", format)
	fmt.Fprint(os.Stderr,
		colorstring.Color(fmt.Sprintf(format, args...)))
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

func (t *Tosa) Run(args []string) int {
	var debug bool
	flags := flag.NewFlagSet(args[0], flag.ContinueOnError)
	flags.BoolVar(&debug, "debug", false, "")

	// Parse flag
	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeParseFlagsError
	}

	if debug {
		os.Setenv(EnvDebug, "1")
		Debugf("Run as DEBUG mode")
	}

	parsedArgs := flags.Args()
	if len(parsedArgs) != 1 {
		PrintErrorf("Invalid argument: you must set sha.")
		return ExitCodeBadArgs
	}

	sha := parsedArgs[0]
	Debugf("sha: %s", sha)

	return tosa(sha)
}

func tosa(sha string) int {
	client, err := NewClient()
	if err != nil {
		return ExitCodeError
	}

	return openPr(client, sha)
}

func openPr(client *APIClient, sha string) int {
	pr, err := client.PullRequest(sha)
	if err != nil || pr == nil {
		return ExitCodePullRequestNotFound
	}

	openErr := open.Run(*pr.HTMLURL)
	if openErr != nil {
		return ExitCodeError
	}
	return ExitCodeOK
}
