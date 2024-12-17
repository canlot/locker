package internals

import (
	"errors"
	"log"
	"strconv"
	"strings"
)

const AppVersion = "1.0.0.1"

const MinDBVersion = "1.0.0"
const MaxDBVersion = "1.999.999.999"

const DBSchemaVersion = "1.0.0"

type Version struct {
	Version     string
	Name        string
	Description string
}

func GetAllVersions() ([5]Version, error) {
	versions := [5]Version{}

	versions[0] = Version{
		Version:     AppVersion,
		Name:        "AppVersion",
		Description: "Version of the application",
	}

	tx, err := Database.Begin(false)
	if err != nil {
		return versions, errors.New("Database is not accessible")
	}
	defer tx.Rollback()

	dbVersionBytes, err := getValue(tx, []byte(DBVersionName), BucketVersion)
	if err != nil {
		return versions, errors.New("Could not read version from database")
	}
	dbVersionString := string(dbVersionBytes)

	versions[1] = Version{
		Version:     dbVersionString,
		Name:        "DBVersion",
		Description: "Database version of current database db_locker.db in this directory",
	}
	versions[2] = Version{
		Version:     MinDBVersion,
		Name:        "MinDBVersion",
		Description: "Minimum supported database version",
	}
	versions[3] = Version{
		Version:     MaxDBVersion,
		Name:        "MaxDBVersion",
		Description: "Maximum supported database version",
	}
	versions[4] = Version{
		Version:     DBSchemaVersion,
		Name:        "DBSchemaVersion",
		Description: "Database version if new database would be created",
	}
	return versions, nil
}

func CompareVersions() {
	tx, err := Database.Begin(false)
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	dbVersionBytes, err := getValue(tx, []byte(DBVersionName), BucketVersion)
	if err != nil {
		log.Fatal(err)
	}
	dbVersionString := string(dbVersionBytes)
	dbVersion, err := GetVersionInInt(dbVersionString)
	if err != nil {
		log.Fatal(err)
	}
	var minVersion int
	if MinDBVersion != "" {
		minVersion, err = GetVersionInInt(MinDBVersion)
		if err != nil {
			log.Fatalln("minVersion err:", err)
		}
	}
	if MinDBVersion != "" && dbVersion < minVersion {
		log.Fatalf("DB version %s is less than min allowed version %s", dbVersionString, MinDBVersion)
	}

	var maxVersion int
	if MaxDBVersion != "" {
		maxVersion, err = GetVersionInInt(MaxDBVersion)
		if err != nil {
			log.Fatalln("maxVersion err:", err)
		}
	}
	if MaxDBVersion != "" && dbVersion > maxVersion {
		log.Fatalf("DB version %s is greater than max allowed version %s", dbVersionString, MaxDBVersion)
	}

}

func GetVersionInInt(version string) (int, error) {
	subVersions := strings.Split(version, ".")
	factor := 1
	versionNumber := 0
	if len(subVersions) < 1 || len(subVersions) > 4 {
		return 0, errors.New("Version string is invalid")
	}
	switch len(subVersions) {
	case 1:
		factor = 1000 * 1000 * 1000
	case 2:
		factor = 1000 * 1000
	case 3:
		factor = 1000
	}
	for i := len(subVersions) - 1; i >= 0; i-- {
		number, err := strconv.Atoi(subVersions[i])
		if err != nil {
			return 0, err
		}
		if number > 999 {
			return 0, errors.New("Version could not be higher then 999")
		}
		versionNumber += number * factor
		factor *= 1000
	}
	if versionNumber < 0 {
		return 0, errors.New("Version less than 0")
	}
	return versionNumber, nil
}
