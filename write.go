package gistdb

import (
	"fmt"

	"github.com/google/go-github/github"
)

//Insert adds a new file to the gist. Returns an error if the file already exists
func (c *Connection) Insert(filename string, contents string) error {
	if isFilenameInGistFileList(filename, c.gist.Files) {
		return fmt.Errorf("Attempting to insert file %s which already exists in the Gist", filename)
	}
	c.gist.Files[github.GistFilename(filename)] = github.GistFile{Content: github.String(contents)}
	c.client.Edit(c.ctx, c.gistID, &c.gist)

	return nil
}

//UpdateOne updates a single existing file. Returns an error if the file does not currently exist
func (c *Connection) UpdateOne(filename string, contents string) error {
	if !isFilenameInGistFileList(filename, c.gist.Files) {
		return fmt.Errorf("Attempting to update file %s which does not exist ing the Gist", filename)
	}
	c.gist.Files[github.GistFilename(filename)] = github.GistFile{Content: github.String(contents)}
	c.client.Edit(c.ctx, c.gistID, &c.gist)
	return nil
}

//UpdateMany updates many files in the Gist at once. Warning this calls UpdateOne multiple times so is NOT idempotent
func (c *Connection) UpdateMany(files map[string]string) error {
	for filename, contents := range files {
		err := c.UpdateOne(filename, contents)
		if err != nil {
			return err
		}
	}
	return nil
}
