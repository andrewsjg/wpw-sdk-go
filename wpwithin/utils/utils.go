package utils

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/wptechinnovation/wpw-sdk-go/wpwithin/wpwerrors"
	"github.com/nu7hatch/gouuid"
)

// NewUUID Generate a new UUID
// This function outputs a version 4 (Random) UUID as per https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_.28random.29
func NewUUID() (string, error) {

	uuid, err := uuid.NewV4()

	if err != nil {

		return "", wpwerrors.GetError(wpwerrors.UUIDGEN, err)
	}

	return uuid.String(), nil
}

// ReadLocalUUID Read a local UUID value
// Input a file path that contains the UUID string, only on one line with no line break
func ReadLocalUUID(path string) (string, error) {

	f, err := os.Open(path)

	if err != nil {

		return "", wpwerrors.GetError(wpwerrors.UUID_FILE_READ, fmt.Sprintf("path=%s", path), err)
	}

	defer f.Close()

	reader := bufio.NewReader(f)

	line, _, err := reader.ReadLine()

	if err != nil {

		return "", wpwerrors.GetError(wpwerrors.UUID_FILE_READ, err)
	}

	return string(line), nil
}

// FileExists Checks if a path exists (Directory or File)
// Returns true if exists, false if not
func FileExists(path string) (bool, error) {

	if _, err := os.Stat(path); os.IsNotExist(err) {

		return false, nil
	}

	return true, nil
}

// WriteString Write a string to a file
// Input the path to file and input string to be written
// Set truncate to true if you would like to truncate the file before writing or false if appending is desired
func WriteString(path, input string, truncate bool) error {

	flags := os.O_RDWR | os.O_CREATE

	if truncate {
		flags |= os.O_TRUNC
	}

	file, err := os.OpenFile(path, flags, 0666)
	if err != nil {
		return wpwerrors.GetError(wpwerrors.IO_OPEN, fmt.Sprintf("path=%s", path), err)
	}
	defer file.Close()

	_, err = file.Write([]byte(input))

	if err != nil {

		return wpwerrors.GetError(wpwerrors.IO_WRITE, err)
	}

	err = file.Sync()

	if err != nil {

		return wpwerrors.GetError(wpwerrors.IO_SYNC, err)
	}

	return nil
}

// TimeFormatISO convert an instance of time.time into an ISO8601 formatted timestamp
func TimeFormatISO(input time.Time) string {

	return input.Format("2006-01-02T15:04:05Z")
}

// ParseISOTime parse an ISO8601 time string into a time object
func ParseISOTime(input string) (time.Time, error) {

	return time.Parse("2006-01-02T15:04:05Z", input)
}

// DoUnitConvertFormat similar to DoUnitConvert but a format string can be provided using "<amount>"
// to represent the actual amount. Also note that an empty format string can be provided
func DoUnitConvertFormat(amount, exponent int, format string) string {

	var strAmount string

	strAmount = fmt.Sprintf("%.2f", float32(amount)/100)

	if !strings.EqualFold(format, "") {

		return strings.Replace(format, "<amount>", strAmount, -1)
	}

	return strAmount
}

// DoUnitConvert convert an integer amount that represents minor currency units into a String formatted to represent
// major units
func DoUnitConvert(amount, exponent int) string {

	return DoUnitConvertFormat(amount, exponent, "")
}

// ToFixed ..
func ToFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(Round(num*output)) / output
}

// Round ..
func Round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}
