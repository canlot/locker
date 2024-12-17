package internals

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"main/cryptography"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func GetMagicString() []byte {
	return []byte{76, 111, 99, 107, 101, 114, 58}
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func isFile(filePath string) (bool, error) {
	stats, err := os.Stat(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, err
		}
		if err != nil {
			return false, err
		}
	}
	if stats.IsDir() {
		return false, errors.New("path is a directory")
	}
	return true, nil
}

func isDir(dirPath string) (bool, error) {
	stats, err := os.Stat(dirPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, err
		}
		if err != nil {
			return false, err
		}
	}
	if !stats.IsDir() {
		return false, errors.New("path is not a directory")
	}
	return true, nil
}

func GetPathsForEncryption(sourcePath, destinationPath string) (sPath, dPath string, err error) {

	// source path checks

	if sourcePath == "" {
		return "", "", errors.New("Source path is empty")
	}

	srcPathExists, err := pathExists(sourcePath)
	if err != nil {
		return "", "", err
	}
	if !srcPathExists {
		return "", "", errors.New("Source path does not exist")
	}
	srcIsFile, err := isFile(sourcePath)
	if err != nil {
		return "", "", err
	}
	if !srcIsFile {
		return "", "", errors.New("Source path is not a file")
	}

	// destination path checks

	if destinationPath == "" {
		destinationPath = sourcePath + ".lock"
		return sourcePath, destinationPath, nil
	}

	dstPathExists, err := pathExists(destinationPath)
	if err != nil {
		return "", "", err
	}

	if dstPathExists {
		dstIsDir, err := isDir(destinationPath)
		if err != nil {
			return "", "", err
		}
		filename := filepath.Base(sourcePath)

		if dstIsDir {
			destinationPath = filepath.Join(destinationPath, (filename + ".lock"))
		}

	}

	dstPathExists, err = pathExists(destinationPath)
	if err != nil {
		return "", "", err
	}

	if dstPathExists {
		dstIsFile, err := isFile(destinationPath)
		if err != nil {
			return "", "", err
		}
		if dstIsFile {
			return "", "", errors.New("Destination file already exists")
		}
		dstIsDir, err := isDir(destinationPath)
		if err != nil {
			return "", "", err
		}
		if dstIsDir {
			return "", "", errors.New("Destination is a directory")
		}
	}

	return sourcePath, destinationPath, nil

}
func GetPathsForDecryption(sourcePath, destinationPath string) (sPath, dPath string, err error) {
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

func fileIsEncrypted(path string) (bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer file.Close()

	marker := make([]byte, len(GetMagicString()))
	uid := make([]byte, 16)

	// magic string at start comparison
	n, err := file.Read(marker)
	if err != nil {
		return false, err
	}
	if n < len(GetMagicString()) {
		return false, nil
	}
	if !bytes.Equal(GetMagicString(), marker) {
		return false, nil
	}
	////

	n, err = file.Read(uid)
	if err != nil {
		return false, err
	}
	if n < len(uid) {
		return false, nil
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
	sourcePath, destinationPath, err := GetPathsForEncryption(sourcePath, destinationPath)
	if err != nil {
		return err
	}

	alreadyEncrypted, err := fileIsEncrypted(sourcePath)
	if err != nil {
		return err
	}
	if alreadyEncrypted {
		return errors.New("File already encrypted")
	}

	sourceFile, err := os.Open(sourcePath)
	defer sourceFile.Close()
	if err != nil {
		return err
	}

	tx, err := Database.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	uid, err := uuid.New().MarshalText()
	if err != nil {
		return err
	}
	if len(uid) != 36 {
		return errors.New("uid has not the expected lenght")
	}
	err = ThrowErrorIfUidAlreadyExist(tx, uid, BucketFileHash, BucketFileInformation, BucketFilePasswordEncrypted)
	if err != nil {
		return err
	}
	publicKey, err := getPublicKey(tx)
	if err != nil {
		return err
	}

	destinationFile, err := os.Create(destinationPath)
	if err != nil {
		return err
	}

	defer destinationFile.Close()

	n, err := destinationFile.Write(GetMagicString())
	if err != nil {
		return err
	}
	if n != len(GetMagicString()) {
		return errors.New("Magix string has not been written")
	}
	n, err = destinationFile.Write(uid)
	if err != nil {
		return err
	}
	if n != 36 {
		return errors.New("UUID has not been written")
	}
	randomPassword := cryptography.GenerateRandomBytes()
	hashBytes, err := cryptography.EncryptFileSymmetricWithHash(randomPassword, sourceFile, destinationFile)
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

	err = saveValue(tx, uid, hashBytes, BucketFileHash)
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
	sourcePath, destinationPath, err := GetPathsForDecryption(sourcePath, destinationPath)
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

	marker := make([]byte, len(GetMagicString()))
	byteCount, err := encryptedFile.Read(marker)
	if err != nil {
		return err
	}
	if byteCount != len(GetMagicString()) {
		return errors.New("Not enough bytes read")
	}
	if !bytes.Equal(marker, GetMagicString()) {
		return errors.New("File has no marker")
	}
	uid := make([]byte, 36)
	byteCount, err = encryptedFile.Read(uid)
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
	hashBytes, err := cryptography.DecryptFileSymmetricWithHash(filePasswordDecrypted, encryptedFile, decryptedFile)
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
	if !bytes.Equal(fileHash, hashBytes) {
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
