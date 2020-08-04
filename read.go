package gistdb

import (
	"fmt"

	"github.com/google/go-github/github"
)

//FindOne returns the single named file from the gist if it exists. Returns an error file doesn't exist
func (c Connection) FindOne(filename string) (*string, error) {
	if !isFilenameInGistFileList(filename, c.gist.Files) {
		return nil, fmt.Errorf("Filename %s not found in Gist", filename)
	}
	return c.gist.Files[github.GistFilename(filename)].Content, nil
}
