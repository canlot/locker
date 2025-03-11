package tests

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

var BaseFolderAbsolute string
var TestFolderAbsolute string
var ArtifactsFolderAbsolute string

var currentFolder string

var encryptedFilePathRelative string
var decryptedFilePathAbsolute string

func SetUpTestFolders() (err error) {

	_, file, _, ok := runtime.Caller(0)
	if !ok {
		err = errors.New("could not determine test folder")
	}

	path := filepath.Dir(file)
	fmt.Printf("Current root test path: %s\n", path)

	if !filepath.IsAbs(path) {
		return errors.New("TestFolderAbsolute is not an absolute path")
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
func TeardownTest() {
	err := os.Chdir(BaseFolderAbsolute)
	if err != nil {
		panic(err)
	}
	err = os.RemoveAll(TestFolderAbsolute)
	if err != nil {
		panic(err)
	}
}
