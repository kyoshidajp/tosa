package tosa

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/oauth2"

	"github.com/github/hub/github"
	api "github.com/google/go-github/github"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
)

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
			return nil, errors.New("Pull Request is not found")
		}
		return a.PullRequest(sha)
	}

	return &res.Issues[0], nil
}

func Repository(client *api.Client) (*api.Repository, error) {
	repo, err := github.LocalRepo()
	if err != nil {
		return nil, err
	}
	prj, err := repo.MainProject()
	if err != nil {
		return nil, err
	}

	res, _, err := client.Repositories.Get(context.Background(), prj.Owner, prj.Name)
	return res, err
}
