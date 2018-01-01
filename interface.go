package tosa

import (
	"io"

	api "github.com/google/go-github/github"
)

type Tosa struct {
	Argv   []string
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer

	args []string

	err error
}

type APIClient struct {
	client     *api.Client
	repository *api.Repository
}
