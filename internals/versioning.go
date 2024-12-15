package internals

import (
	"errors"
	"strconv"
	"strings"
)

const AppVersion = "1.0.0.1"

const MinDBVersion = "1.0.0"
const MaxDBVersion = ""

const DBVersion = "1.0.0"

func GetVersionInInt(version string) (int, error) {
	subVersions := strings.Split(version, ".")
	versionNumber := 0
	factor := 1
	for i := len(subVersions) - 1; i < 0; i-- {
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
