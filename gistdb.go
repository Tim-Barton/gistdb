package gistdb

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

//Connection is the handle for all access to the Gists
type Connection struct {
	ctx       context.Context
	client    *github.Client
	gistID    string
	gistFiles map[github.GistFilename]github.GistFile
}

//NewConnection creates a new Connection to the listed gist, returning error if the id cannot be found
func NewConnection(pat string, gistid string) (*Connection, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: pat},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	gists, err := getUserGists(ctx, client)
	if err != nil {
		fmt.Printf("Error retrieving Gists: %s", err)
		return nil, fmt.Errorf("Could not connect to Github. Err: %s", err)
	}
	if !isGistIdinGistList(gistid, gists) {
		return nil, fmt.Errorf("Gist of id %s not found for this user. Check details and try again", gistid)
	}

	gist, _, err := client.Gists.Get(ctx, gistid)
	if err != nil {
		fmt.Printf("Error retrieving Gist %s. Err: %s", gistid, err)
		return nil, fmt.Errorf("Gist of id %s not retreivable. Check details and try again", gistid)
	}
	files := gist.Files

	return &Connection{ctx, client, gistid, files}, nil
}
