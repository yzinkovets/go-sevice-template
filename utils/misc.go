package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/guregu/null/v5"
)

func GetNumRate(rate string) (int, error) {
	pow := 1
	if strings.HasSuffix(rate, "Gbps") {
		pow = 1000
	}

	numStr := strings.TrimRight(rate, "GMbps")

	num, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0, err
	}

	return int(num * float64(pow)), nil
}

func GetNumRateIntOrNull(s string) null.Int {
	if v, err := GetNumRate(s); err == nil {
		return null.IntFrom(int64(v))
	}
	return null.Int{}
}

func GetNumOrNull(s string) null.Int {
	if v, err := GetNum(s); err == nil {
		return null.IntFrom(int64(v))
	}
	return null.Int{}
}

func GetNum(s string) (int, error) {
	v := 0

	n, err := fmt.Sscanf(s, "%d", &v)
	if err != nil {
		return 0, err
	}
	if n != 1 {
		return 0, fmt.Errorf("expected 1 value, got %d", n)
	}
	return v, nil
}

func GetIntFromString(s string) null.Int {
	if v, err := GetNum(s); err == nil {
		return null.IntFrom(int64(v))
	}
	return null.Int{}
}

func GetInt16FromString(s string) null.Int16 {
	if v, err := GetNum(s); err == nil {
		return null.Int16From(int16(v))
	}
	return null.Int16{}
}

func GetRateRssi(rate string, rssi string, ret string) string {
	rates := strings.Split(rate, ",")
	if len(rates) <= 1 {
		if ret == "rate" {
			return rate
		}
		return rssi
	}

	r1, _ := GetNumRate(rates[0])
	r2, _ := GetNumRate(rates[1])

	if ret == "rate" {
		if r1 > r2 {
			return rates[0]
		}
		return rates[1]
	}

	rssis := strings.Split(rssi, ",")
	if len(rssis) <= 1 {
		return rssi
	}

	if r1 > r2 {
		return rssis[0]
	}
	return rssis[1]
}

func GetBoolFromString(s string) bool {
	if s == "1" || strings.ToLower(s) == "true" {
		return true
	}
	return false
}

func SetMaxMapValue(m map[string]string, key string, value string) {
	if value == "" {
		return
	}

	if current, ok := m[key]; !ok || len(current) < len(value) {
		m[key] = value
	}
}

// CompactString removes all tabs, newlines and carriage returns from the input string
func CompactString(input string) string {
	re := regexp.MustCompile(`[\t\n\r]+`)
	return re.ReplaceAllString(input, "")
}
