package tests

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"main/internals"
	"os"
	"path/filepath"
	"testing"
)

var baseFolder string
var testFolder string
var currentFolder string
var folderTestPathAbsolute string
var folderTestPathRelative string
var encryptedFilePathAbsolute string
var encryptedFilePathRelative string
var decryptedFilePathAbsolute string
var decryptedFilePathRelative string

func init() {

}

func createFoldersAndFilesForTesting() {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	baseFolder = path
	err = os.Mkdir(filepath.Join(path, "test"), 0777)
	if err != nil {
		panic(err)
	}
	testFolder = filepath.Join(path, "test")
	err = os.Chdir(testFolder)
	if err != nil {
		panic(err)
	}
	currentFolder = testFolder

	folderName := "Test"

	folderTestPathAbsolute = filepath.Join(currentFolder, folderName)
	folderTestPathRelative = folderName
	err = os.Mkdir(folderTestPathAbsolute, 0777)
	if err != nil {
		panic(err)
	}

	decryptedFilePathAbsolute = filepath.Join(folderTestPathAbsolute, "decrypted.txt")
	decryptedFilePathRelative = filepath.Join(folderName, "decrypted.txt")
	encryptedFilePathAbsolute = filepath.Join(folderTestPathAbsolute, "encrypted.lock")
	encryptedFilePathRelative = filepath.Join(folderName, "encrypted.lock")

	file, err := os.Create(encryptedFilePathAbsolute)
	defer file.Close()
	if err != nil {
		panic(err)
	}
	file, err = os.Create(decryptedFilePathAbsolute)
	defer file.Close()
	if err != nil {
		panic(err)
	}
}
func TestMain(m *testing.M) {
	createFoldersAndFilesForTesting()
	defer deleteFoldersAndFilesAfterTesting()
	m.Run()
}
func deleteFoldersAndFilesAfterTesting() {
	err := os.Chdir(baseFolder)
	if err != nil {
		panic(err)
	}
	err = os.RemoveAll(testFolder)
	if err != nil {
		panic(err)
	}
}
func pathAndPrint(path string) string {
	fmt.Println(path)
	return path
}
func Test_getPathsForEncryption(t *testing.T) {

	//invalid cases, no source file provided
	_, _, err := internals.GetPathsForEncryption("", encryptedFilePathAbsolute)
	assert.NotNil(t, err)
	_, _, err = internals.GetPathsForEncryption("", encryptedFilePathRelative)
	assert.NotNil(t, err)

	// invalid case, source file already encrypted, destination fine
	//_, _, err = internals.GetPathsForEncryption(encryptedFilePathRelative, filepath.Join(currentFolder, "test.lock"))
	//assert.NotNil(t, err)

	//invalid case, dest file already exists
	_, _, err = internals.GetPathsForEncryption(decryptedFilePathRelative, encryptedFilePathAbsolute)
	assert.NotNil(t, err)

	//valid case
	dPath, sPath, err := internals.GetPathsForEncryption(decryptedFilePathAbsolute, filepath.Join(currentFolder, "test.lock"))
	assert.Nil(t, err)
	assert.Equal(t, dPath, decryptedFilePathAbsolute)
	assert.Equal(t, sPath, filepath.Join(currentFolder, "test.lock"))

	//valid case, providing source file and directory as destination, dest file should be dir + filename + .lock
	sPath, dPath, err = internals.GetPathsForEncryption(decryptedFilePathRelative, folderTestPathRelative)
	assert.Nil(t, err)
	assert.Equal(t, sPath, decryptedFilePathRelative)
	assert.Equal(t, dPath, filepath.Join(folderTestPathRelative, "decrypted.txt.lock"))
}
