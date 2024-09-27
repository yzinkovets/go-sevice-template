package utils

import (
	"testing"
)

func TestGetNumRate(t *testing.T) {
	cases := []struct {
		rate string
		num  int
		err  bool
	}{
		{"100Mbps", 100, false},
		{"1Gbps", 1000, false},
		{"1.5Gbps", 1500, false},
		{"1.5Mbps", 1, false},
		{"1.5", 1, false},
		{"10", 10, false},
		{"", 0, true},
		{"1B", 0, true},
		{"-10", -10, false},
	}

	for _, c := range cases {
		num, err := GetNumRate(c.rate)
		if err != nil && !c.err {
			t.Errorf("Unexpected error: %v", err)
		} else if err == nil && c.err {
			t.Errorf("Expected error, got nil")
		} else if num != c.num {
			t.Errorf("Expected %d, got %d", c.num, num)
		}
	}
}

func TestGetNum(t *testing.T) {
	cases := []struct {
		rate string
		num  int
		err  bool
	}{
		{"100Mbps", 100, false},
		{"-1Gbps", -1, false},
		{"1.5Gbps", 1, false},
		{"-1.5", -1, false},
		{"10", 10, false},
		{"", 0, true},
		{"B", 0, true},
		{"s", 0, true},
	}

	for _, c := range cases {
		t.Log(c.rate, c.num, c.err)
		num, err := GetNum(c.rate)
		if err != nil && !c.err {
			t.Errorf("Unexpected error: %v", err)
		} else if err == nil && c.err {
			t.Errorf("Expected error, got nil")
		} else if num != c.num {
			t.Errorf("Expected %d, got %d", c.num, num)
		}
	}
}

func TestGetRateRssi(t *testing.T) {
	cases := []struct {
		rate string
		rssi string
		ret  string
		out  string
	}{
		{"100Mbps", "10", "rate", "100Mbps"},
		{"100Mbps,1Gbps", "10,20", "rate", "1Gbps"},
		{"100Mbps,1Gbps", "10,20", "rssi", "20"},
		{"100Mbps", "10", "rssi", "10"},
		{"1Gbps,100Mbps", "10,20", "rssi", "10"},
		{"1Gbps,100Mbps", "10,20", "rate", "1Gbps"},
		{"", "10", "rssi", "10"},
		{"", "10", "rate", ""},
	}

	for _, c := range cases {
		out := GetRateRssi(c.rate, c.rssi, c.ret)
		if out != c.out {
			t.Log(c.rate, c.rssi, c.ret, c.out)
			t.Errorf("Expected %s, got %s", c.out, out)
		}
	}
}

func TestGetBoolFromString(t *testing.T) {
	cases := []struct {
		s   string
		out bool
	}{
		{"1", true},
		{"true", true},
		{"True", true},
		{"TRUE", true},
		{"0", false},
		{"false", false},
		{"False", false},
		{"FALSE", false},
		{"", false},
		{"Truee", false},
	}

	for _, c := range cases {
		out := GetBoolFromString(c.s)
		if out != c.out {
			t.Log(c.s, c.out)
			t.Errorf("Expected %v, got %v", c.out, out)
		}
	}
}

func TestCompactString(t *testing.T) {
	cases := []struct {
		s   string
		out string
	}{
		{"", ""},
		{"a", "a"},
		{"aa\r\nb\tc", "aabc"},
		{"aa\r\nb\tc\n", "aabc"},
	}

	for _, c := range cases {
		out := CompactString(c.s)
		if out != c.out {
			t.Log(c.s, c.out)
			t.Errorf("Expected %v, got %v", c.out, out)
		}
	}
}
