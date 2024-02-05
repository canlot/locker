package internals

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	bolt "go.etcd.io/bbolt"
	"io"
	"main/cryptography"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type DBStore struct {
}

type DataInformation struct {
	Label      string
	CreateTime time.Time
}

type FileInformation struct {
	Path       string
	CreateTime time.Time
}

var Store DBStore

func IsDatabaseEmpty() (bool, error) {
	fmt.Println("IsDatabaseEmpty")
	var err error
	if Database == nil {
		fmt.Println("Database is nil")
	}
	tx, err := Database.Begin(false)
	defer fmt.Println("IsDatabaseEmpty Commited")
	defer tx.Rollback()
	if err != nil {
		return false, err
	}
	var empty bool
	empty, err = IsBucketEmpty(tx, BucketLoginInformation)
	if err != nil {
		return false, err
	}
	if !empty {
		return false, nil
	}
	empty, err = IsBucketEmpty(tx, BucketLoginPrivateKeyEncrypted)
	if err != nil {
		return false, err
	}
	if !empty {
		return false, nil
	}
	empty, err = IsBucketEmpty(tx, BucketLoginPrivateKeyEncrypted)
	if err != nil {
		return false, err
	}
	if !empty {
		return false, nil
	}

	return true, nil
}
func IsBucketEmpty(tx *bolt.Tx, bucketName string) (bool, error) {
	bucket := tx.Bucket([]byte(bucketName))
	if bucket == nil {
		return false, errors.New("Bucket not created: " + bucketName)
	}
	c := bucket.Cursor()
	first, _ := c.First()
	if first != nil {
		return false, nil
	}
	return true, nil
}
func CreateFirstLoginWithRSAKeys(username, password string) error {
	//fmt.Println("CreateFirstLoginWithRSAKeys")
	createTime := time.Now()
	passwordHash, err := cryptography.GenerateUserHash([]byte(password))
	if err != nil {
		return err
	}
	u := uuid.New()
	uuidStringBytes, err := u.MarshalText()
	if err != nil {
		return err
	}
	login := Login{Login: username, CreateTime: createTime}
	loginBytes, err := json.Marshal(login)
	if err != nil {
		return err
	}
	privateKey, publicKey, err := cryptography.GenerateRSAKeys()
	//fmt.Printf("private key: %x\n", privateKey)
	//fmt.Printf("public key: %x\n", publicKey)
	if err != nil {
		return err
	}
	privateKeyHash, err := cryptography.GetSha256Hash(privateKey)
	//fmt.Printf("private key hash: %x\n", privateKeyHash)
	if err != nil {
		return err
	}
	privateKeyEncrypted, err := cryptography.EncryptDataSymmetric(passwordHash, privateKey)
	if err != nil {
		return err
	}
	tx, err := Database.Begin(true)
	defer tx.Rollback()
	if err != nil {
		return err
	}
	userInfoBucket := tx.Bucket([]byte(BucketLoginInformation))
	err = userInfoBucket.Put(uuidStringBytes, loginBytes)
	if err != nil {
		return err
	}
	userPrivateKeyEncryptedBucket := tx.Bucket([]byte(BucketLoginPrivateKeyEncrypted))
	err = userPrivateKeyEncryptedBucket.Put(uuidStringBytes, privateKeyEncrypted)
	if err != nil {
		return err
	}
	privateKeyHashBucket := tx.Bucket([]byte(BucketPrivateKeyHash))
	err = privateKeyHashBucket.Put([]byte(PrivateKeyHashKeyName), privateKeyHash)
	if err != nil {
		return err
	}
	publicKeyBucket := tx.Bucket([]byte(BucketPublicKey))
	err = publicKeyBucket.Put([]byte(PublicKeyKeyName), publicKey)
	if err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}
func CreateLoginWithExistingRSAKeys(existingLogin, existingLoginPassword, newLogin, newLoginPassword string) error {
	//fmt.Println("CreateLoginWithExistingRSAKeys")
	tx, err := Database.Begin(true)
	defer tx.Rollback()
	if err != nil {
		return err
	}
	passwordHashExistingLogin, err := cryptography.GenerateUserHash([]byte(existingLoginPassword))
	if err != nil {
		return err
	}
	passwordHashNewLogin, err := cryptography.GenerateUserHash([]byte(newLoginPassword))
	if err != nil {
		return err
	}

	userid, err := getLoginId(existingLogin, tx)
	if err != nil {
		return err
	}

	decryptedPrivateKey, err := getAndDecryptPrivateKey(userid, []byte(passwordHashExistingLogin), tx)
	if err != nil {
		return err
	}
	privateKeyHash, err := getPrivateKeyHash(tx)
	if err != nil {
		return err
	}
	//fmt.Printf("private key: %x\n", decryptedPrivateKey)

	if err != nil {
		return err
	}
	generatedPrivateKeyHash, err := cryptography.GetSha256Hash(decryptedPrivateKey)
	if err != nil {
		return err
	}
	//fmt.Printf("private key hash: %x\n", generatedPrivateKeyHash)
	if bytes.Equal(generatedPrivateKeyHash, privateKeyHash) != true {
		return errors.New("Private keys do not match, wrong password for existing login: " + existingLogin)
	}
	encryptedPrivateKey, err := cryptography.EncryptDataSymmetric([]byte(passwordHashNewLogin), decryptedPrivateKey)
	if err != nil {
		return err
	}
	u := uuid.New()
	uuidStringBytes, err := u.MarshalText()
	if err != nil {
		return err
	}
	createTime := time.Now()

	login := Login{Login: newLogin, CreateTime: createTime}
	loginBytes, err := json.Marshal(login)
	if err != nil {
		return err
	}
	loginInformationBucket := tx.Bucket([]byte(BucketLoginInformation))
	err = loginInformationBucket.Put(uuidStringBytes, loginBytes)
	if err != nil {
		return err
	}
	loginPrivateKeyBucket := tx.Bucket([]byte(BucketLoginPrivateKeyEncrypted))
	err = loginPrivateKeyBucket.Put(uuidStringBytes, encryptedPrivateKey)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
func getLoginId(loginString string, tx *bolt.Tx) (uid []byte, err error) { // return error if no login found
	loginInformationBucket := tx.Bucket([]byte(BucketLoginInformation))
	c := loginInformationBucket.Cursor()

	var userid []byte
	for k, v := c.First(); k != nil; k, v = c.Next() {
		var login Login
		err = json.Unmarshal(v, &login)
		if err != nil {
			return nil, err
		}
		if login.Login == loginString {
			userid = k
		}
	}
	if userid == nil {
		return nil, errors.New("No existing login found: " + loginString)
	}
	return userid, nil
}
func getAndDecryptPrivateKey(loginId, passwordHash []byte, tx *bolt.Tx) (privateKey []byte, err error) {
	loginPrivateKeyBucket := tx.Bucket([]byte(BucketLoginPrivateKeyEncrypted))
	privateKeyEncrypted := loginPrivateKeyBucket.Get(loginId)
	if privateKeyEncrypted == nil {
		return nil, errors.New("No private key exists for this login: " + string(loginId))
	}
	decryptedPrivateKey, err := cryptography.DecryptDataSymmetric(passwordHash, privateKeyEncrypted)
	if err != nil {
		return nil, err
	}
	//fmt.Printf("private key: %x\n", decryptedPrivateKey)
	return decryptedPrivateKey, nil
}

func getPrivateKeyHash(tx *bolt.Tx) (hash []byte, err error) {
	privateKeyHashBucket := tx.Bucket([]byte(BucketPrivateKeyHash))
	privateKeyHash := privateKeyHashBucket.Get([]byte(PrivateKeyHashKeyName))
	if privateKeyHash == nil {
		return nil, errors.New("Private key hash is not set")
	}
	return privateKeyHash, nil
}

func getPublicKey(tx *bolt.Tx) (publicKey []byte, err error) {
	publicKeyBucket := tx.Bucket([]byte(BucketPublicKey))
	publicKey = publicKeyBucket.Get([]byte(PublicKeyKeyName))
	if publicKey == nil {
		return nil, errors.New("Public key is empty")
	}
	return publicKey, nil
}

func ListAllUsers() (logins []Login, err error) {
	tx, err := Database.Begin(false)
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}
	loginInfoBucket := tx.Bucket([]byte(BucketLoginInformation))
	if loginInfoBucket == nil {
		return nil, err
	}
	c := loginInfoBucket.Cursor()
	for k, v := c.First(); k != nil; k, v = c.Next() {
		var login Login
		err = json.Unmarshal(v, &login)
		if err != nil {
			return nil, err
		}
		logins = append(logins, login)

	}
	return logins, nil
}
func EncryptData(label, plainData string) error {
	tx, err := Database.Begin(true)
	defer tx.Rollback()
	if err != nil {
		return err
	}
	data := DataInformation{Label: label, CreateTime: time.Now()}
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	uid, err := uuid.New().MarshalText()
	if err != nil {
		return err
	}
	randomPasswordByte := cryptography.GenerateRandomBytes()
	//fmt.Printf("Random password: %x\n", randomPasswordByte)
	dataInfoBucket := tx.Bucket([]byte(BucketDataInformation))
	if dataInfoBucket == nil {
		return errors.New("Bucket: " + BucketDataInformation + " doesn't exist")
	}
	err = dataInfoBucket.Put(uid, dataBytes)
	if err != nil {
		return err
	}

	encryptedData, err := cryptography.EncryptDataSymmetric(randomPasswordByte, []byte(plainData))
	if err != nil {
		return err
	}
	publicKey, err := getPublicKey(tx)
	if err != nil {
		return err
	}
	encryptedRandomPassword, err := cryptography.EncryptDataAsymmetric(publicKey, randomPasswordByte)
	if err != nil {
		return err
	}
	//fmt.Printf("Encrypted random password: %x\n", encryptedRandomPassword)

	encrypteDataBucket := tx.Bucket([]byte(BucketDataEncrypted))
	if encrypteDataBucket == nil {
		return errors.New("Bucket: " + BucketDataEncrypted + " doesn't exist")
	}
	err = encrypteDataBucket.Put(uid, encryptedData)
	if err != nil {
		return err
	}
	encryptedDataPasswordBucket := tx.Bucket([]byte(BucketDataPasswordEncrypted))
	if encryptedDataPasswordBucket == nil {
		return errors.New("Bucket: " + BucketDataPasswordEncrypted + " doesn't exist")
	}
	err = encryptedDataPasswordBucket.Put(uid, encryptedRandomPassword)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}
func DecryptData(dataid, login, password string) (dataInfo DataInformation, plainData string, err error) {
	dataIdBytes := []byte(dataid)
	tx, err := Database.Begin(false)
	defer tx.Rollback()
	if err != nil {
		return dataInfo, "", err
	}
	loginId, err := getLoginId(login, tx)
	if err != nil {
		return dataInfo, "", err
	}
	passwordHash, err := cryptography.GenerateUserHash([]byte(password))
	privateKey, err := getAndDecryptPrivateKey(loginId, passwordHash, tx)
	if err != nil {
		return dataInfo, "", err
	}
	privateKeyHash, err := getPrivateKeyHash(tx)
	if err != nil {
		return dataInfo, "", err
	}
	generatedPrivateKeyHash, err := cryptography.GetSha256Hash(privateKey)
	if err != nil {
		return dataInfo, "", err
	}
	if !bytes.Equal(privateKeyHash, generatedPrivateKeyHash) {
		return dataInfo, "", errors.New("Hashes are different")
	}
	passwordDataBucket := tx.Bucket([]byte(BucketDataPasswordEncrypted))
	if passwordDataBucket == nil {
		return dataInfo, "", errors.New("No bucket found")
	}
	encryptedRandomPasswordForData := passwordDataBucket.Get(dataIdBytes)
	//fmt.Printf("encrypted data password: %x\n", encryptedRandomPasswordForData)
	if encryptedRandomPasswordForData == nil {
		return dataInfo, "", errors.New("No encrypted password found")
	}
	decryptedRandomPasswordForData, err := cryptography.DecryptDataAsymmetric(privateKey, encryptedRandomPasswordForData)
	if err != nil {
		return dataInfo, "", err
	}
	encryptedDataBucket := tx.Bucket([]byte(BucketDataEncrypted))
	if encryptedDataBucket == nil {
		return dataInfo, "", errors.New("No bucket found")
	}
	encryptedData := encryptedDataBucket.Get(dataIdBytes)
	if encryptedData == nil {
		return dataInfo, "", errors.New("No data found")
	}
	//fmt.Printf("Decrypted password: %x\n", decryptedRandomPasswordForData)
	decryptedData, err := cryptography.DecryptDataSymmetric(decryptedRandomPasswordForData, encryptedData)
	if err != nil {
		return dataInfo, "", err
	}
	dataInformationBucket := tx.Bucket([]byte(BucketDataInformation))
	if dataInformationBucket == nil {
		return dataInfo, "", errors.New("No bucket found")
	}
	dataInformationBytes := dataInformationBucket.Get(dataIdBytes)
	if dataInformationBytes == nil {
		return dataInfo, "", errors.New("No data found")
	}
	var dataInformation DataInformation
	err = json.Unmarshal(dataInformationBytes, &dataInformation)
	if err != nil {
		return dataInfo, "", err
	}
	return dataInformation, string(decryptedData), nil
}
func ListAllData() (keys []string, dataInfo []DataInformation, err error) {
	tx, err := Database.Begin(false)
	defer tx.Rollback()
	if err != nil {
		return nil, nil, err
	}
	dataInfoBucket := tx.Bucket([]byte(BucketDataInformation))
	if dataInfoBucket == nil {
		return nil, nil, errors.New("No bucket found")
	}
	c := dataInfoBucket.Cursor()
	for k, v := c.First(); k != nil; k, v = c.Next() {
		var data DataInformation
		err = json.Unmarshal(v, &data)
		if err != nil {
			return nil, nil, err
		}
		dataInfo = append(dataInfo, data)
		keys = append(keys, string(k))

	}
	return keys, dataInfo, nil
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
func ThrowErrorIfUidAlreadyExist(tx *bolt.Tx, uid []byte, buckets ...string) error {
	for i := range buckets {
		bucket := tx.Bucket([]byte(buckets[i]))
		if bucket == nil {
			return errors.New("Bucket: " + buckets[i] + " is empty")
		}
		value := bucket.Get(uid)
		if value != nil {
			return errors.New("A glitch in the universe has been happen, uuid already exist, please retry")
		}
	}
	return nil

}
func getValue(tx *bolt.Tx, uid []byte, bucketName string) (value []byte, err error) {
	bucket := tx.Bucket([]byte(bucketName))
	if bucket == nil {
		return nil, errors.New("Bucket: " + bucketName + " doesn't exist")
	}
	value = bucket.Get(uid)
	if value == nil {
		return nil, errors.New("No value for id: " + string(uid) + " in bucket: " + bucketName + " found")
	}
	return value, nil
}
func saveValue(tx *bolt.Tx, uid, value []byte, bucketName string) error {
	bucket := tx.Bucket([]byte(bucketName))
	if bucket == nil {
		return errors.New("Bucket: " + bucketName + " doesn't exist")
	}
	err := bucket.Put(uid, value)
	if err != nil {
		return err
	}
	return nil
}
func deleteValue(tx *bolt.Tx, uid []byte, bucketName string) error {
	bucket := tx.Bucket([]byte(bucketName))
	if bucket == nil {
		return errors.New("Bucket: " + bucketName + " doesn't exist")
	}
	err := bucket.Delete(uid)
	if err != nil {
		return err
	}
	return nil
}
func deleteValues(tx *bolt.Tx, uid []byte, bucketNames ...string) error {
	for i := range bucketNames {
		err := deleteValue(tx, uid, bucketNames[i])
		if err != nil {
			return err
		}
	}
	return nil
}
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
	//fmt.Println(destinationDirectory)
	//fmt.Println(destinationFileName)
	dPath = filepath.Join(destinationDirectory, destinationFileName)
	newFile, err := os.Open(dPath)
	defer newFile.Close()
	if err == nil { // file already exist
		dPath = filepath.Join(destinationDirectory, ("unlocked_" + destinationFileName))
		return sourcePath, dPath, nil
	}
	return sourcePath, dPath, nil
}
func EncryptFile(sourcePath, destinationPath string) error {
	sourcePath, destinationPath, err := getPathsForEncryption(sourcePath, destinationPath)
	if err != nil {
		return err
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
	fmt.Println("Done encryption: " + time.Now().String())
	fileHash, err := cryptography.GetSha256HashFile(sourcePath)
	fmt.Println("Done Hashing " + time.Now().String())
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
	fmt.Println(sourcePath)
	fmt.Println(destinationPath)
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
