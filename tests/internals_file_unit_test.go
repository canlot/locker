package tests

import (
	"crypto/sha256"
	"fmt"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/scrypt"
	"io"
	"main/cryptography"
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
func getSha256HashFile(filePath string) (hash []byte, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

func copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	err = os.WriteFile(dst, input, 0644)
	if err != nil {
		return err
	}
	return nil
}
func setUpTest() error {
	path, err := os.Getwd()
	if err != nil {
		return err
	}
	baseFolder = path
	err = os.Mkdir(filepath.Join(path, "test"), 0777)
	if err != nil {
		return err
	}
	currentFolder = filepath.Join(path, "test")

	err = copyFile(filepath.Join("artifacts", "testfile.txt"), filepath.Join(currentFolder, "testfile.txt"))
	if err != nil {
		return err
	}

	err = os.Chdir(currentFolder)
	if err != nil {
		return err
	}

	testFolder = "Test"

	folderTestPathAbsolute = filepath.Join(currentFolder, testFolder)
	folderTestPathRelative = testFolder
	err = os.Mkdir(folderTestPathAbsolute, 0777)
	if err != nil {
		return err
	}

	decryptedFilePathAbsolute = filepath.Join(folderTestPathAbsolute, "test.txt")
	decryptedFilePathRelative = filepath.Join(testFolder, "test.txt")
	encryptedFilePathAbsolute = filepath.Join(folderTestPathAbsolute, "test.txt.lock")
	encryptedFilePathRelative = filepath.Join(testFolder, "test.txt.lock")

	file, err := os.Create(encryptedFilePathAbsolute)
	defer file.Close()
	if err != nil {
		return err
	}
	file, err = os.Create(decryptedFilePathAbsolute)
	defer file.Close()
	if err != nil {
		return err
	}
	return nil
}

func TestMain(m *testing.M) {
	m.Run()
}
func teardownTest() {
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
func Test_GetPathsForEncryption(t *testing.T) {
	setUpTest()
	defer teardownTest()
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
	dPath, sPath, err := internals.GetPathsForEncryption(decryptedFilePathAbsolute, filepath.Join(testFolder, "test.lock"))
	assert.Nil(t, err)
	assert.Equal(t, dPath, decryptedFilePathAbsolute)
	assert.Equal(t, sPath, filepath.Join(testFolder, "test.lock"))

	//invalid case, file already exist
	sPath, dPath, err = internals.GetPathsForEncryption(decryptedFilePathRelative, folderTestPathRelative)
	assert.NotNil(t, err)

	unecryptedFilePath := filepath.Join(currentFolder, "testfile.txt")
	ecryptedFilePath := filepath.Join(currentFolder, "testfile.txt.lock")

	//valid case, providing source file and destination
	sPath, dPath, err = internals.GetPathsForEncryption(unecryptedFilePath, ecryptedFilePath)
	assert.Nil(t, err)
	assert.Equal(t, sPath, unecryptedFilePath)
	assert.Equal(t, dPath, ecryptedFilePath)

	//valid case, providing source file and directory as destination, dest file should be dir + filename + .lock
	sPath, dPath, err = internals.GetPathsForEncryption(unecryptedFilePath, testFolder)
	assert.Nil(t, err)
	assert.Equal(t, sPath, unecryptedFilePath)
	assert.Equal(t, dPath, filepath.Join(testFolder, "testfile.txt.lock"))

}

func Test_GetPathsForDecryption(t *testing.T) {
	setUpTest()
	defer teardownTest()
	// invalid case, no source Path
	_, _, err := internals.GetPathsForDecryption("", decryptedFilePathAbsolute)
	assert.NotNil(t, err)

	_, _, err = internals.GetPathsForDecryption(encryptedFilePathRelative, "")
	assert.NotNil(t, err)

}

func Test_EnsureEncryptionAndDecryptionHaveSameResult(t *testing.T) {
	setUpTest()
	defer teardownTest()
	bytePassword, err := scrypt.Key([]byte("test"), nil, 32768, 8, 2, 32)
	if err != nil {
		t.Fail()
		return
	}
	unencryptedFilePath := filepath.Join(testFolder, "testfile.txt")
	encryptedFilePath := filepath.Join(testFolder, "testfile.txt.lock")
	decryptedFilePath := filepath.Join(testFolder, "testfile_decrypted.txt")

	unencryptedFileHash, err := getSha256HashFile(unencryptedFilePath)
	if err != nil {
		t.Fail()
		return
	}

	fileSrcUnencrypted, err := os.Open(unencryptedFilePath)
	defer fileSrcUnencrypted.Close()
	if err != nil {
		t.Fatal(err)
	}

	fileDstEncrypted, err := os.Create(encryptedFilePath)
	defer fileDstEncrypted.Close()
	if err != nil {
		t.Fatal(err)
	}

	byteHash, err := cryptography.EncryptFileSymmetricWithHash(bytePassword, fileSrcUnencrypted, fileDstEncrypted)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, unencryptedFileHash, byteHash)

	fileSrcUnencrypted.Close()
	fileDstEncrypted.Close()

	fileSrcEncrypted, err := os.Open(encryptedFilePath)
	defer fileSrcEncrypted.Close()
	if err != nil {
		t.Fatal(err)
	}
	fileDstDecrypted, err := os.Create(decryptedFilePath)
	defer fileDstDecrypted.Close()
	if err != nil {
		t.Fatal(err)
	}

	byteHashDecrypted, err := cryptography.DecryptFileSymmetricWithHash(bytePassword, fileSrcEncrypted, fileDstDecrypted)
	if err != nil {
		t.Fatal(err)
	}
	fileSrcEncrypted.Close()
	fileDstDecrypted.Close()

	assert.Equal(t, unencryptedFileHash, byteHashDecrypted)

	decryptedFileHash, err := getSha256HashFile(decryptedFilePath)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, unencryptedFileHash, decryptedFileHash)

}
