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

// Exit codes are in value that represnet an exit code for a paticular error
const (
	ExitCodeOK int = 0

	// Errors start at 10
	ExitCodeError = 10 + iota
	ExitCodeParseFlagsError
	ExitCodeBadArgs
	ExitCodePullRequestNotFound
	ExitCodeOpenPageError
)

// Debugf prints debug output when EnvDebug is given
func Debugf(format string, args ...interface{}) {
	if env := os.Getenv(EnvDebug); len(env) != 0 {
		log.Printf("[DEBUG] "+format+"\n", args...)
	}
}

// PrintErrorf prints error message on console
func PrintErrorf(format string, args ...interface{}) {
	format = fmt.Sprintf("[red]%s[reset]\n", format)
	fmt.Fprint(os.Stderr,
		colorstring.Color(fmt.Sprintf(format, args...)))
}

// CLI is the command line object
type CLI struct {
	outStream, errStream io.Writer
}

// APIClient is access/operate Github object
type APIClient struct {
	client     *api.Client
	repository *api.Repository
}

// Run invokes the CLI with the given arguments
func (c *CLI) Run(args []string) int {
	var (
		debug   bool
		url     bool
		apiurl  bool
		newline bool
		version bool
	)
	flags := flag.NewFlagSet(args[0], flag.ContinueOnError)
	flags.Usage = func() {
		fmt.Fprint(c.errStream, helpText)
	}
	flags.BoolVar(&debug, "debug", false, "")
	flags.BoolVar(&debug, "d", false, "")
	flags.BoolVar(&url, "url", false, "")
	flags.BoolVar(&url, "u", false, "")
	flags.BoolVar(&apiurl, "apiurl", false, "")
	flags.BoolVar(&apiurl, "a", false, "")
	flags.BoolVar(&newline, "newline", true, "")
	flags.BoolVar(&newline, "n", true, "")
	flags.BoolVar(&version, "version", false, "")
	flags.BoolVar(&version, "v", false, "")

	// Parse flag
	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeParseFlagsError
	}

	if debug {
		os.Setenv(EnvDebug, "1")
		Debugf("Run as DEBUG mode")
	}

	if version {
		fmt.Fprintf(c.outStream, fmt.Sprintf("%s\n", Version))
		return ExitCodeOK
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

	var status int
	if url {
		status = printUrl(client, sha, newline)
	} else if apiurl {
		status = printAPIUrl(client, sha, newline)
	} else {
		status = openPr(client, sha)
	}

	return status
}

func printUrl(client *APIClient, sha string, newline bool) int {
	Debugf("Print PullRequest URL")

	pr, err := client.PullRequest(sha)
	if err != nil || pr == nil {
		return ExitCodePullRequestNotFound
	}

	lastc := ""
	if newline {
		lastc = "\n"
	}
	format := fmt.Sprintf("%s%s", *pr.HTMLURL, lastc)
	fmt.Fprint(os.Stdout, format)

	return ExitCodeOK
}

func printAPIUrl(client *APIClient, sha string, newline bool) int {
	Debugf("Print API URL")

	pr, err := client.PullRequest(sha)
	if err != nil || pr == nil {
		return ExitCodePullRequestNotFound
	}

	lastc := ""
	if newline {
		lastc = "\n"
	}
	format := fmt.Sprintf("%s%s", *pr.URL, lastc)
	fmt.Fprint(os.Stdout, format)

	return ExitCodeOK
}

func openPr(client *APIClient, sha string) int {
	pr, err := client.PullRequest(sha)
	if err != nil || pr == nil {
		return ExitCodePullRequestNotFound
	}

	browser, err := GetBrowser()
	if err != nil {
		return ExitCodeOpenPageError
	}
	Debugf("browser: %s", browser)

	url := *pr.HTMLURL
	Debugf("URL: %s", url)

	var openErr error
	if browser == "" {
		openErr = open.Run(url)
	} else {
		openErr = open.RunWith(url, browser)
	}
	if openErr != nil {
		PrintErrorf("Could not open page. Check your browser and URL.")
		return ExitCodeError
	}

	return ExitCodeOK
}

// NewClient creates APIClient
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

// PullRequest gets PullRequest object
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

// Repository returns api.Repository
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

var helpText = `Usage: tosa [options...] sha

tosa is a tool to open the PullRequest page.

You must specify commit sha what you want to know PullRequest.

Options:

  -u, --url      Print the PullRequest url.

  -a, --apiurl   Print the Issue API url.

  -n, --newline  If -u(--url) or --apiurl option is specified, print
                 the url with newline character at last.

  -d, --debug    Enable debug mode.
                 Print debug log.

  -h, --help     Show this help message and exit.

  -v, --version  Print current version.
`
