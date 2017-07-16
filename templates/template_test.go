package templates

import (
	"bytes"
	"testing"
)

func TestMaster(t *testing.T) {
	err := CreateConfig("test")
	if err != nil {
		t.Fatal("CreateConfig:", err)
	}

	err = SaveMaster("test", "", "default", "Test Master")
	if err != nil {
		t.Fatal("SaveMaster:", err)
	}

	err = SaveHost("test", "test-host", nil, nil)
	if err != nil {
		t.Fatal("SaveHost:", err)
	}

	b := &bytes.Buffer{}
	err = Execute("test-host", "default", b, nil)
	if err != nil {
		t.Fatal("Execute:", err)
	}
	if b.String() != "Test Master" {
		t.Fatal("Expected:", "Test Master", "Found:", b.String())
	}
}

func TestErrors(t *testing.T) {
}
