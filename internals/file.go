package internals

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io"
	"main/cryptography"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func getPathsForEncryption(sourcePath, destinationPath string) (sPath, dPath string, err error) {
	if sourcePath == "" {
		return "", "", errors.New("Source path is empty")
	}
	sfile, err := os.Open(sourcePath)
	if err != nil {
		return "", "", err
	}
	defer sfile.Close()
	sfileInfo, err := sfile.Stat()
	if err != nil {
		return "", "", err
	}
	if sfileInfo.IsDir() {
		return "", "", errors.New("Path is not a file")
	}

	if destinationPath != "" {
		dfile, err := os.Open(destinationPath)
		if err != nil {
		}
		defer dfile.Close()
		dfileInfo, err := dfile.Stat()
		if err != nil {
			return "", "", err
		}
		if !(dfileInfo.IsDir()) {
			return sourcePath, destinationPath, nil
		} else {
			sourceFileName := filepath.Base(sourcePath)
			dPath = filepath.Join(destinationPath, (sourceFileName + ".lock"))
			return sourcePath, dPath, nil
		}
	} else {
		dPath = sourcePath + ".lock"
		return sourcePath, dPath, nil
	}
}
func getPathsForDecryption(sourcePath, destinationPath string) (sPath, dPath string, err error) {
	if sourcePath == "" {
		return "", "", errors.New("Source path is empty")
	}
	sfile, err := os.Open(sourcePath)
	if err != nil {
		return "", "", err
	}
	defer sfile.Close()
	sfileInfo, err := sfile.Stat()
	if err != nil {
		return "", "", err
	}
	if sfileInfo.IsDir() {
		return "", "", errors.New("Path is not a file")
	}

	fileName := filepath.Base(sourcePath)

	var destinationDirectory string
	if destinationPath != "" {
		dfile, err := os.Open(destinationPath)
		if err != nil {
		}
		defer dfile.Close()
		dfileInfo, err := dfile.Stat()
		if err != nil {
			return "", "", err
		}
		if !(dfileInfo.IsDir()) {
			destinationDirectory = filepath.Dir(destinationPath)
			fileName = filepath.Base(destinationPath)
		}
	} else {
		destinationDirectory = filepath.Dir(sourcePath)
	}
	destinationFileName, found := strings.CutSuffix(fileName, ".lock")
	if found != true {
		return "", "", errors.New("Source file has no ending .lock")
	}

	dPath = filepath.Join(destinationDirectory, destinationFileName)
	newFile, err := os.Open(dPath)
	defer newFile.Close()
	if err == nil { // file already exist
		dPath = filepath.Join(destinationDirectory, ("unlocked_" + destinationFileName))
		return sourcePath, dPath, nil
	}
	return sourcePath, dPath, nil
}
func fileAlreadyEncrypted(path string) (bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer file.Close()
	uid := make([]byte, 36)
	_, err = file.Read(uid)
	if err != nil {
		return false, err
	}
	_, err = uuid.ParseBytes(uid)
	if err == nil { //if file has uuid at start
		_, exist := strings.CutSuffix(path, FileEncryptionEnding)
		if exist { // and has .lock file ending
			return true, nil // then file is encrypted
		}
	}
	return false, nil // otherwise file is not encrypted

}
func ListAllFiles() (ids, hashes []string, fileInfo []FileInformation, err error) {
	tx, err := Database.Begin(false)
	if err != nil {
		return nil, nil, nil, err
	}
	defer tx.Rollback()
	fileInfoBucket := tx.Bucket([]byte(BucketFileInformation))
	if fileInfoBucket == nil {
		return nil, nil, nil, errors.New("No bucket " + BucketFileInformation + "found")
	}
	c := fileInfoBucket.Cursor()
	for k, v := c.First(); k != nil; k, v = c.Next() {
		var file FileInformation
		err = json.Unmarshal(v, &file)
		if err != nil {
			return nil, nil, nil, err
		}
		fileInfo = append(fileInfo, file)
		ids = append(ids, string(k))
		hash, err := getValue(tx, k, BucketFileHash)
		if err != nil {
			return nil, nil, nil, err
		}
		hashes = append(hashes, hex.EncodeToString(hash))
	}
	return ids, hashes, fileInfo, nil
}
func EncryptFile(sourcePath, destinationPath string) error {
	sourcePath, destinationPath, err := getPathsForEncryption(sourcePath, destinationPath)
	if err != nil {
		return err
	}

	alreadyEncrypted, err := fileAlreadyEncrypted(sourcePath)
	if err != nil {
		return err
	}
	if alreadyEncrypted {
		return errors.New("File already encrypted")
	}

	sourceFile, err := os.Open(sourcePath)

	if err != nil {
		return err
	}

	defer sourceFile.Close()
	destinationFile, err := os.Create(destinationPath)
	if err != nil {
		return err
	}

	defer destinationFile.Close()
	uid, err := uuid.New().MarshalText()
	if err != nil {
		return err
	}
	fileReader := io.Reader(sourceFile)
	fileWriter := io.Writer(destinationFile)
	if len(uid) != 36 {
		return errors.New("uid has not the expected lenght")
	}
	fileWriter.Write(uid)
	tx, err := Database.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	err = ThrowErrorIfUidAlreadyExist(tx, uid, BucketFileHash, BucketFileInformation, BucketFilePasswordEncrypted)
	if err != nil {
		return err
	}
	publicKey, err := getPublicKey(tx)
	if err != nil {
		return err
	}
	randomPassword := cryptography.GenerateRandomBytes()
	fmt.Println("Start encryption: " + time.Now().String())
	err = cryptography.EncryptFileSymmetric(randomPassword, fileReader, fileWriter)
	if err != nil {
		return err
	}
	randomPasswordEncrypted, err := cryptography.EncryptDataAsymmetric(publicKey, randomPassword)
	if err != nil {
		return err
	}
	err = saveValue(tx, uid, randomPasswordEncrypted, BucketFilePasswordEncrypted)
	if err != nil {
		return err
	}
	dataInfo := FileInformation{Path: sourcePath, CreateTime: time.Now()}
	dataInfoBytes, err := json.Marshal(&dataInfo)
	if err != nil {
		return err
	}
	err = saveValue(tx, uid, dataInfoBytes, BucketFileInformation)
	if err != nil {
		return err
	}
	fileHash, err := cryptography.GetSha256HashFile(sourcePath)
	if err != nil {
		return err
	}
	err = saveValue(tx, uid, fileHash, BucketFileHash)
	if err != nil {
		return err
	}
	sourceFile.Close()
	destinationFile.Close()
	err = os.Remove(sourcePath)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}
func DecryptFile(sourcePath, destinationPath, login, password string) error {
	sourcePath, destinationPath, err := getPathsForDecryption(sourcePath, destinationPath)
	if err != nil {
		return err
	}
	encryptedFile, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer encryptedFile.Close()
	decryptedFile, _ := os.Create(destinationPath)
	if err != nil {
		return err
	}
	defer encryptedFile.Close()
	fileReader := io.Reader(encryptedFile)
	fileWriter := io.Writer(decryptedFile)
	uid := make([]byte, 36)
	byteCount, err := fileReader.Read(uid)
	if err != nil {
		return err
	}
	if byteCount != 36 {
		return errors.New("Not an exact amount of byted read")
	}
	_, err = uuid.ParseBytes(uid)
	if err != nil {
		return err
	}
	tx, err := Database.Begin(true)
	if err != nil {
		return err
	}
	loginId, err := getLoginId(login, tx)
	if err != nil {
		return err
	}
	passwordHash, err := cryptography.GenerateUserHash([]byte(password))
	if err != nil {
		return err
	}
	privateKey, err := getAndDecryptPrivateKey(loginId, passwordHash, tx)
	if err != nil {
		return err
	}
	privateKeyHash, err := getPrivateKeyHash(tx)
	if err != nil {
		return err
	}
	privateKeyHashGenerated, err := cryptography.GetSha256Hash(privateKey)
	if !bytes.Equal(privateKeyHash, privateKeyHashGenerated) {
		return errors.New("Private key hash is different, password maybe incorrect")
	}
	filePasswordEncrypted, err := getValue(tx, uid, BucketFilePasswordEncrypted)
	if err != nil {
		return err
	}
	filePasswordDecrypted, err := cryptography.DecryptDataAsymmetric(privateKey, filePasswordEncrypted)
	if err != nil {
		return err
	}
	err = cryptography.DecryptFileSymmetric(filePasswordDecrypted, fileReader, fileWriter)
	if err != nil {
		return err
	}
	err = decryptedFile.Close()
	if err != nil {
		return err
	}
	fileHash, err := getValue(tx, uid, BucketFileHash)
	if err != nil {
		return err
	}
	generatedFileHash, err := cryptography.GetSha256HashFile(destinationPath)
	if err != nil {
		return err
	}
	if !bytes.Equal(fileHash, generatedFileHash) {
		return errors.New("File hash is not equal, file has been changed")
	}
	err = deleteValues(tx, uid, BucketFileHash, BucketFileInformation, BucketFilePasswordEncrypted)
	if err != nil {
		return err
	}
	encryptedFile.Close()
	decryptedFile.Close()
	err = os.Remove(sourcePath)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}
