package internals

import (
	"errors"
	"log"
	"strconv"
	"strings"
)

const AppVersion = "1.0.0.1"

const MinDBVersion = "1.0.0"
const MaxDBVersion = ""

const DBVersion = "1.0.0"

func CompareVersions() {
	dbVersion, err := GetVersionInInt(DBVersion)
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
		log.Fatalf("DB version %s is less than min version %s", dbVersion, MinDBVersion)
	}

	var maxVersion int
	if MaxDBVersion != "" {
		maxVersion, err = GetVersionInInt(MaxDBVersion)
		if err != nil {
			log.Fatalln("maxVersion err:", err)
		}
	}
	if MaxDBVersion != "" && dbVersion > maxVersion {
		log.Fatalf("DB version %s is greater than max version %s", dbVersion, MaxDBVersion)
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
