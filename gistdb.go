package gistdb

import (
	"context"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

//Connection is the handle for all access to the Gists
type Connection struct {
	ctx    context.Context
	client *github.Client
	gistID string
	gist   github.Gist
}

//NewConnection creates a new Connection to the listed gist, returning error if the id cannot be found
func NewConnection(pat string, gistid string) (*Connection, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: pat},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	conn := Connection{ctx: ctx,
		client: client,
		gistID: gistid}

	err := conn.loadGistFiles()
	if err != nil {
		return nil, err
	}

	return &conn, nil
}
