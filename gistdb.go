package gistdb

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type GistDBConnection struct {
	ctx       context.Context
	client    *github.Client
	gistID    string
	gistFiles map[github.GistFilename]github.GistFile
}

func isGistIdinGistList(searchId string, gistlist []*github.Gist) bool {
	for _, v := range gistlist {
		if *(v.ID) == searchId {
			return true
		}
	}
	return false
}

func NewConnection(pat string, gistid string) (*GistDBConnection, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: pat},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	gists, err := getUserGists(ctx, client)
	if err != nil {
		fmt.Printf("Error retrieving Gists: %s", err)
		return nil, fmt.Errorf("Could not connect to Gist %s. Err: %s", gistid, err)
	}
	if !isGistIdinGistList(gistid, gists) {
		return nil, fmt.Errorf("Gist of id %s not found for this user. Check details and try again", gistid)
	}

	gist, _, err := client.Gists.Get(ctx, gistid)
	if err != nil {
		fmt.Printf("Error retrieving Gist %s. Err: %s", gistid, err)
		return nil, fmt.Errorf("Gist if id %s not retreivable. Check details and try again", gistid)
	}
	files := gist.Files

	return &GistDBConnection{ctx, client, gistid, files}, nil
}

func getUserGists(ctx context.Context, client *github.Client) ([]*github.Gist, error) {
	gists, _, err := client.Gists.List(ctx, "", nil)
	return gists, err
}
