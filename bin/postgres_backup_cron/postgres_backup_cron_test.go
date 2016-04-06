package main

import (
	"testing"

	"bytes"

	"time"

	. "github.com/bborbe/assert"
)

func TestDoFail(t *testing.T) {
	writer := bytes.NewBufferString("")
	err := do(writer, func(host string, port int, user string, pass string, database string, targetDir string) error {
		return nil
	}, "", 0, "", "", "", "", time.Minute, false, "/tmp/lock")
	if err = AssertThat(err, NotNilValue()); err != nil {
		t.Fatal(err)
	}
}

func TestDoSuccess(t *testing.T) {
	writer := bytes.NewBufferString("")
	err := do(writer, func(host string, port int, user string, pass string, database string, targetDir string) error {
		return nil
	}, "host", 5432, "user", "pass", "db", "/tmp", time.Minute, true, "/tmp/lock")
	if err = AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
}