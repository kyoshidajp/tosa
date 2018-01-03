package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/oauth2"

	"github.com/github/hub/github"
	api "github.com/google/go-github/github"
	"github.com/mitchellh/colorstring"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/skratchdot/open-golang/open"
)

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

// CLI is the command line object
type CLI struct {
	outStream, errStream io.Writer
}

type APIClient struct {
	client     *api.Client
	repository *api.Repository
}

// Run invokes the CLI with the given arguments.
func (c *CLI) Run(args []string) int {
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

func NewClient() (*APIClient, error) {
	homeDir, err := homedir.Dir()
	if err != nil {
		return nil, err
	}

	confPath := filepath.Join(homeDir, ".config", "tosa")
	err = os.Setenv("HUB_CONFIG", confPath)
	if err != nil {
		return nil, err
	}

	c := github.CurrentConfig()
	host, err := c.DefaultHost()
	if err != nil {
		return nil, err
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: host.AccessToken},
	)
	tc := oauth2.NewClient(context.Background(), ts)

	client := api.NewClient(tc)
	repo, err := Repository(client)
	if err != nil {
		return nil, err
	}

	return &APIClient{
		client:     client,
		repository: repo,
	}, nil
}

func (a *APIClient) PullRequest(sha string) (*api.Issue, error) {
	res, _, err := a.client.Search.Issues(context.Background(),
		fmt.Sprintf("%s is:merged repo:%v", sha, *a.repository.FullName), nil)
	if err != nil {
		return nil, err
	}

	if len(res.Issues) == 0 {
		a.repository = a.repository.Parent
		if a.repository == nil {
			PrintErrorf("Pull Request is not found")
			return nil, nil

		}
		Debugf("Searching parent repository: %s", *a.repository.FullName)
		return a.PullRequest(sha)
	}

	return &res.Issues[0], nil
}

func Repository(client *api.Client) (*api.Repository, error) {
	localRepo, err := github.LocalRepo()
	if err != nil {
		return nil, err
	}
	prj, err := localRepo.MainProject()
	if err != nil {
		return nil, err
	}

	repo, _, err := client.Repositories.Get(context.Background(), prj.Owner, prj.Name)
	if err != nil {
		PrintErrorf("Repository not found.\n%s", err)
		return nil, err
	}
	return repo, err
}
