package gistdb

import (
	"fmt"

	"github.com/google/go-github/github"
)

func (c Connection) getUserGists() ([]*github.Gist, error) {
	gists, _, err := c.client.Gists.List(c.ctx, "", nil)
	return gists, err
}
func (c *Connection) loadGistFiles() error {
	gists, err := c.getUserGists()
	if err != nil {
		fmt.Printf("Error retrieving Gists: %s", err)
		return fmt.Errorf("Could not connect to Github. Err: %s", err)
	}
	if !isGistIdinGistList(c.gistID, gists) {
		return fmt.Errorf("Gist of id %s not found for this user. Check details and try again", c.gistID)
	}

	gist, _, err := c.client.Gists.Get(c.ctx, c.gistID)
	if err != nil {
		fmt.Printf("Error retrieving Gist %s. Err: %s", c.gistID, err)
		return fmt.Errorf("Gist of id %s not retreivable. Check details and try again", c.gistID)
	}
	c.gistFiles = gist.Files
	return nil
}
