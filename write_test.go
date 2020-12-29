package gistdb

import (
	"context"
	"fmt"
	"gistdb/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-github/github"
)

func TestInsert(t *testing.T) {
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

	oldName := "name"
	newName := "newName"
	newContents := "This is some new data"

	err := conn.Insert(oldName, newContents)
	if err == nil {
		t.Error("Should return an error when inserting file of the same name")
	}

	mock.EXPECT().Edit(ctx, gistID, gomock.Any()).Return(nil, nil, fmt.Errorf("Throw an error")).Times(1)
	err = conn.Insert(newName, newContents)
	if err == nil {
		t.Error("Expecting the Edit to throw an error")
	}
	if isFilenameInGistFileList(newName, conn.gist.Files) {
		t.Error("If Edit threw an error then our local cache shouldn't update")
	}

	mock.EXPECT().Edit(ctx, gistID, gomock.Any()).Return(nil, nil, nil).Times(1)
	err = conn.Insert(newName, newContents)
	if err != nil {
		t.Error("Expecting edit to succeed and no other errors")
	}
	if !isFilenameInGistFileList(newName, conn.gist.Files) {
		t.Error("If Edit succeeds then our local cache should be updated")
	}

}

func TestUpdateOne(t *testing.T) {
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

	oldName := "name"
	newName := "newName"
	newContents := "This is some new data"

	err := conn.UpdateOne(newName, newContents)
	if err == nil {
		t.Error("Should return an error when updating file with a new name")
	}

	mock.EXPECT().Edit(ctx, gistID, gomock.Any()).Return(nil, nil, fmt.Errorf("Throw an error")).Times(1)
	err = conn.UpdateOne(oldName, newContents)
	if err == nil {
		t.Error("Expecting the Edit to throw an error")
	}
	if *conn.gist.Files[github.GistFilename(oldName)].Content == newContents {
		t.Error("If Edit threw an error then the local cache shouldn't be the new contents")
	}

	mock.EXPECT().Edit(ctx, gistID, gomock.Any()).Return(nil, nil, nil).Times(1)
	err = conn.UpdateOne(oldName, newContents)
	if err != nil {
		t.Error("Not expecting any other errors")
	}
	if *conn.gist.Files[github.GistFilename(oldName)].Content != newContents {
		t.Error("With no errors local cache should be updated")
	}

}

func TestUpdateMany(t *testing.T) {
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

	oldName := "name"
	newName := "newName"
	newContents := "This is some new data"

	newFiles := map[string]string{newName: newContents}
	updatedFiles := map[string]string{oldName: newContents}

	err := conn.UpdateMany(newFiles)
	if err == nil {
		t.Error("Should return an error when updating file with a new name")
	}

	mock.EXPECT().Edit(ctx, gistID, gomock.Any()).Return(nil, nil, fmt.Errorf("Throw an error")).Times(1)
	err = conn.UpdateMany(updatedFiles)
	if err == nil {
		t.Error("Expecting the Edit to throw an error")
	}
	if *conn.gist.Files[github.GistFilename(oldName)].Content == newContents {
		t.Error("If Edit threw an error then the local cache shouldn't be the new contents")
	}

	mock.EXPECT().Edit(ctx, gistID, gomock.Any()).Return(nil, nil, nil).Times(1)
	err = conn.UpdateMany(updatedFiles)
	if err != nil {
		t.Error("Not expecting any other errors")
	}
	if *conn.gist.Files[github.GistFilename(oldName)].Content != newContents {
		t.Error("With no errors local cache should be updated")
	}

}
