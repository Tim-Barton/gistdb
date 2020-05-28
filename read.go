package gistdb

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
)

func getUserGists(ctx context.Context, client *github.Client) ([]*github.Gist, error) {
	gists, _, err := client.Gists.List(ctx, "", nil)
	return gists, err
}

func SelectStar(c Connection, filename string) (*string, error) {
	if !isFilenameInGistFileList(filename, c.gistFiles) {
		return nil, fmt.Errorf("Filename %s not found in Gist", filename)
	}
	return c.gistFiles[github.GistFilename(filename)].Content, nil
}
