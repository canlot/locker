package internals

import (
	"errors"
)

type DBStore struct {
}

var Store DBStore

func (s DBStore) DatabaseEmpty() (bool, error) {
	db, err := Database.Begin(false)
	if err != nil {
		return false, err
	}
	bucket := db.Bucket([]byte(BucketLoginInformation))
	if bucket == nil {
		return errors.New("Bucket not created")
	}
	c := bucket.Cursor()
	first, _ := c.First()
	if first != nil {
		return false
	}
}
