package gistdb

import (
	"context"

	"github.com/google/go-github/github"
)

type GistDb interface {
	Insert(string, string) error
	UpdateOne(string, string) error
	UpdateMany(map[string]string) error
	FindOne(string) (string, error)
}

type GithubGists interface {
	Get(context.Context, string) (*github.Gist, *github.Response, error)
	List(context.Context, string, *github.GistListOptions) ([]*github.Gist, *github.Response, error)

	Edit(context.Context, string, *github.Gist) (*github.Gist, *github.Response, error)
}
