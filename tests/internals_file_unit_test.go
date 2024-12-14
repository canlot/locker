package tests

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/google/uuid"
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
var artifactsFolder string
var currentFolder string

var encryptedFilePathRelative string
var decryptedFilePathAbsolute string

func init() {

}
func getSha256HashFile(filePath string) (hash []byte, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		derr := file.Close()
		err = errors.Join(err, derr)
	}(file)
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
func createPseudoEncryptedFile(filePath string) (err error) {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		derr := file.Close()
		err = errors.Join(err, derr)
	}(file)
	magic := []byte{76, 111, 99, 107, 101, 114, 58}
	uid, err := uuid.New().MarshalBinary()
	fmt.Println(len(uid))
	if err != nil {
		return err
	}
	file.Write(magic)
	file.Write(uid)
	return nil
}
func setUpTest() error {
	path, err := os.Getwd()
	if err != nil {
		return err
	}
	baseFolder = path

	artifactsFolder = filepath.Join(path, "artifacts")

	testFolder = filepath.Join(path, "running_testcases")
	err = os.Mkdir(testFolder, 0777)
	if err != nil && (!errors.Is(err, os.ErrExist)) {
		return err
	}
	currentFolder = testFolder
	err = os.Chdir(currentFolder)
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

	// copy valid files for testing
	plainFile := "testfile.txt"
	encryptedFile := "encrypted_file.txt.lock"

	err := copyFile(filepath.Join(artifactsFolder, plainFile), filepath.Join(currentFolder, plainFile))
	assert.Nil(t, err)
	err = copyFile(filepath.Join(artifactsFolder, encryptedFile), filepath.Join(currentFolder, encryptedFile))
	assert.Nil(t, err)

	//invalid case, no source file provided
	_, _, err = internals.GetPathsForEncryption("", "non_existent_file.lock")
	assert.NotNil(t, err)

	//invalid case, dest file already exists

	_, _, err = internals.GetPathsForEncryption(plainFile, encryptedFile)
	assert.NotNil(t, err)
	////

	//valid case
	sPath, dPath, err := internals.GetPathsForEncryption(plainFile, "test.lock")
	assert.Nil(t, err)
	assert.Equal(t, sPath, plainFile)
	assert.Equal(t, dPath, "test.lock")
	////

	//valid case
	testDirRelative := "testDir"
	err = os.Mkdir(testDirRelative, 0777)
	assert.Nil(t, err)
	sPath, dPath, err = internals.GetPathsForEncryption(plainFile, testDirRelative)
	assert.Nil(t, err)
	assert.Equal(t, sPath, plainFile)
	assert.Equal(t, dPath, filepath.Join(testDirRelative, "testfile.txt.lock"))
	////

	//invalid case, file already in destination, and no destination filename given, it should generate the same name as the file that already exists
	err = copyFile(filepath.Join(artifactsFolder, encryptedFile), filepath.Join(testDirRelative, "testfile.txt.lock"))
	assert.Nil(t, err)

	_, _, err = internals.GetPathsForEncryption(plainFile, testDirRelative)
	assert.NotNil(t, err)
	///

	unecryptedFilePath := filepath.Join(currentFolder, "testfile.txt")
	encryptedFilePath := filepath.Join(currentFolder, "testfile.txt.lock")

	//valid case, providing source file and destination
	sPath, dPath, err = internals.GetPathsForEncryption(unecryptedFilePath, encryptedFilePath)
	assert.Nil(t, err)
	assert.Equal(t, sPath, unecryptedFilePath)
	assert.Equal(t, dPath, encryptedFilePath)

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

	unencryptedFileName := "testfile.txt"
	encryptedFileName := "testfile.txt.lock"
	decryptedFileName := "testfile_decrypted.txt"
	err = copyFile(filepath.Join(artifactsFolder, unencryptedFileName), filepath.Join(currentFolder, unencryptedFileName))

	unencryptedFileHash, err := getSha256HashFile(unencryptedFileName)
	if err != nil {
		t.Fail()
		return
	}

	fileSrcUnencrypted, err := os.Open(unencryptedFileName)
	defer fileSrcUnencrypted.Close()
	if err != nil {
		t.Fatal(err)
	}

	fileDstEncrypted, err := os.Create(encryptedFileName)
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

	fileSrcEncrypted, err := os.Open(encryptedFileName)
	defer fileSrcEncrypted.Close()
	if err != nil {
		t.Fatal(err)
	}
	fileDstDecrypted, err := os.Create(decryptedFileName)
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

	decryptedFileHash, err := getSha256HashFile(decryptedFileName)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, unencryptedFileHash, decryptedFileHash)

}
