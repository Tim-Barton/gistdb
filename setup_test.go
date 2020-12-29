package gistdb

import (
	"context"
	"fmt"
	"gistdb/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-github/github"
)

func TestLoadGistFiles(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := mocks.NewMockGithubGists(ctrl)
	ctx := context.Background()

	gistID := "1"

	conn := Connection{ctx: ctx,
		client: mock,
		gistID: gistID}

	//error getting the list of gists
	mock.EXPECT().Get(ctx, gistID).Return(nil, nil, nil).Times(0)
	mock.EXPECT().List(ctx, "", nil).Return(make([]*github.Gist, 0), nil, fmt.Errorf("Error")).Times(1)
	_ = conn.loadGistFiles()

	//handle getting an empty list /  list with without provided gist ID
	mock.EXPECT().Get(ctx, gistID).Return(nil, nil, nil).Times(0)
	mock.EXPECT().List(ctx, "", nil).Return(make([]*github.Gist, 0), nil, nil).Times(1)
	_ = conn.loadGistFiles()

	//check when a correct list is returned
	files := map[github.GistFilename]github.GistFile{github.GistFilename("name"): github.GistFile{}}
	mock.EXPECT().Get(ctx, gistID).Return(&github.Gist{Files: files}, nil, nil).Times(1)
	gistList := []*github.Gist{&github.Gist{ID: &gistID}}
	mock.EXPECT().List(ctx, "", nil).Return(gistList, nil, nil).Times(1)
	_ = conn.loadGistFiles()

	//throw an error
	files = map[github.GistFilename]github.GistFile{github.GistFilename("name"): github.GistFile{}}
	mock.EXPECT().Get(ctx, gistID).Return(&github.Gist{Files: files}, nil, fmt.Errorf("Error")).Times(1)
	gistList = []*github.Gist{&github.Gist{ID: &gistID}}
	mock.EXPECT().List(ctx, "", nil).Return(gistList, nil, nil).Times(1)
	err := conn.loadGistFiles()
	if err == nil {
		t.Fail()
	}

}
