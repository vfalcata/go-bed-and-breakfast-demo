// recall any file with suffix "_test.go" is considered a test file by go
package main

import "testing"

func TestRun(t *testing.T) {
	_, err := run()
	if err != nil {
		t.Error("Failed run()")
	}
}
