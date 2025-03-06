package internals

import (
	"github.com/stretchr/testify/assert"
	bolt "go.etcd.io/bbolt"
	"path/filepath"
	"testing"
)
import "main/tests"

func Test_DatabaseCreation(t *testing.T) {
	tests.SetUpTestFolders()
	DatabasePath = tests.TestFolderAbsolute

	CreateDatabaseIfNotExists()
	Database.Close()

	dbPath := filepath.Join(DatabasePath, "db_locker.db")

	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	tx, err := db.Begin(true)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
	bucket := tx.Bucket([]byte("BucketLoginInformation"))
	assert.NotNil(t, bucket)

}
