package tests

import (
	"errors"
	"os"
	"path/filepath"
)

var BaseFolderAbsolute string
var TestFolderAbsolute string
var ArtifactsFolderAbsolute string

var currentFolder string

var encryptedFilePathRelative string
var decryptedFilePathAbsolute string

func SetUpTestFolders() error {
	path, err := os.Getwd()
	if err != nil {
		return err
	}
	if !filepath.IsAbs(path) {
		return err
	}
	BaseFolderAbsolute = path

	ArtifactsFolderAbsolute = filepath.Join(path, "artifacts")

	TestFolderAbsolute = filepath.Join(path, "running_testcases")
	err = os.Mkdir(TestFolderAbsolute, 0777)
	if err != nil && (!errors.Is(err, os.ErrExist)) {
		return err
	}
	currentFolder = TestFolderAbsolute
	err = os.Chdir(currentFolder)
	if err != nil {
		return err
	}

	return nil
}
