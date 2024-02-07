package internals

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"main/cryptography"
	"time"
)

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
