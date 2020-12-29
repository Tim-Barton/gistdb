package gistdb

import (
	"context"
	"gistdb/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-github/github"
)

func TestFindOne(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := mocks.NewMockGithubGists(ctrl)
	ctx := context.Background()

	gistID := "1"

	conn := Connection{ctx: ctx,
		client: mock,
		gistID: gistID}

	//load known data into the connection struct
	content := "This is the file's contents"
	files := map[github.GistFilename]github.GistFile{github.GistFilename("name"): github.GistFile{Content: &content}}
	mock.EXPECT().Get(ctx, gistID).Return(&github.Gist{Files: files}, nil, nil).Times(1)
	gistList := []*github.Gist{&github.Gist{ID: &gistID}}
	mock.EXPECT().List(ctx, "", nil).Return(gistList, nil, nil).Times(1)
	_ = conn.loadGistFiles()

	_, err := conn.FindOne("bad name")
	if err == nil {
		t.Errorf("Should not find 'bad name' in list of files")
	}

	_, err = conn.FindOne("name")
	if err != nil {
		t.Errorf("Should find good name in list of files")
	}
}
