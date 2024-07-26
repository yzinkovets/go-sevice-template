package utils

import (
	"fmt"
	"hash/fnv"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v5"
)

func GetUUID() string {
	return uuid.NewString()
}

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func IsValidMACAddress(mac string) bool {
	var macRegex = regexp.MustCompile(`^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$`)
	return macRegex.MatchString(mac)
}

func hash(s string) string {
	h := fnv.New32()
	h.Write([]byte(s))
	return strconv.FormatInt(int64(h.Sum32()), 16)
}

func GetReqNameForLog(reqName string) string {
	return fmt.Sprintf("[%s] (%s)", hash(GetUUID()), reqName)
}

func Nvl(value string, defaultValue string) string {
	if value == "" {
		return defaultValue
	} else {
		return value
	}
}

func CutGw(gw_uuid string) string {
	if strings.HasPrefix(gw_uuid, "gw") {
		return gw_uuid[2:]
	}
	return gw_uuid
}

func SetGw(gw_uuid string) string {
	gw_uuid = strings.ReplaceAll(gw_uuid, "-", "")

	if !strings.HasPrefix(gw_uuid, "gw") {
		return "gw" + gw_uuid
	}
	return gw_uuid
}

func IsDateValue(stringDate string) bool {
	_, err := time.Parse("2006-01-02", stringDate)
	return err == nil
}

func IsDateTimeValue(stringDateTime string) bool {
	_, err := time.Parse("2006-01-02T15:04:05", stringDateTime)
	return err == nil
}

func TryGetDateTime(stringDateTime string) (time.Time, error) {
	stringDateTime = strings.ReplaceAll(stringDateTime, "-", "")
	stringDateTime = strings.ReplaceAll(stringDateTime, ":", "")
	stringDateTime = strings.ReplaceAll(stringDateTime, "T", "")

	var value time.Time
	var err error
	// "2006-01-02T15:04:05" without - & : & T

	if len(stringDateTime) > 8 {
		value, err = time.Parse("20060102150405", stringDateTime)
	} else {
		value, err = time.Parse("20060102", stringDateTime)
	}

	return value, err
}

func TryGetNullDateTime(stringDateTime string) (null.Time, error) {
	if stringDateTime == "" {
		return null.Time{}, nil
	}

	t, err := TryGetDateTime(stringDateTime)
	if err != nil {
		return null.Time{}, err
	}

	return null.TimeFrom(t), nil
}
