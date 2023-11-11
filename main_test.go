package main

import (
	"testing"
)

type tc struct {
	ip  string
	err error
}

var TestCases []tc = []tc{
	{"219.100.37.141", nil}, // valid IP address
	{"219.100.37.145", nil}, // invalid IP address
	{"127.0.0.1", nil},      // local IP address
	{"nonexistent_ip", nil}, // nonexistent IP address
}

func TestPingVpnServer(t *testing.T) {
	f := NewFastVpn()
	for _, testCase := range TestCases {
		t.Run(testCase.ip, func(t *testing.T) {
			pingTime, err := f.pingVpnServer(testCase.ip)

			// Check if the error matches the expected error
			if err != nil && testCase.err == nil {
				t.Errorf("Expected no error, but got: %v", err)
			} else if err == nil && testCase.err != nil {
				t.Errorf("Expected error: %v, but got no error", testCase.err)
			}

			// Check if pingTime is non-zero if there is no error
			if err == nil && pingTime == 0 {
				t.Error("Expected non-zero ping time, got zero")
			}
		})
	}
}
