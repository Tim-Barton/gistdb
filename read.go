package gistdb

import (
	"fmt"

	"github.com/google/go-github/github"
)

func (c Connection) FindOne(filename string) (*string, error) {
	if !isFilenameInGistFileList(filename, c.gistFiles) {
		return nil, fmt.Errorf("Filename %s not found in Gist", filename)
	}
	return c.gistFiles[github.GistFilename(filename)].Content, nil
}
