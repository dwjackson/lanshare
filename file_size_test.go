package main

import (
	"testing"
)

func TestParseFileSizeInBytes(t *testing.T) {
	testParse(t, "1024", 1024)
}

func TestParseFileSizeInMebibytes(t *testing.T) {
	testParse(t, "1MiB", 1048576)
}

func TestParseFileSizeInGibibytes(t *testing.T) {
	testParse(t, "2GiB", 2*1073741824)
}

func TestParseFileSizeInKibibytes(t *testing.T) {
	testParse(t, "4KiB", 4*1024)
}

func TestFractionalSizesNotAllowed(t *testing.T) {
	_, err := parseFileSize("2.5MiB")
	if err == nil {
		t.Fatalf("This should fail")
	}
	errorMessage := err.Error()
	expectedErrorMessage := "Fractional sizes not allowed"
	if errorMessage != expectedErrorMessage {
		t.Fatalf("Wrong error message: %s", errorMessage)
	}
}

func testParse(t *testing.T, input string, expected int64) {
	actual, err := parseFileSize(input)
	if err != nil {
		t.Fatalf("Error during parse: %v", err)
	} else if actual != expected {
		t.Fatalf("Expected %d, got %d", expected, actual)
	}
}

func TestFormatHumanReadableBytes(t *testing.T) {
	tests := []struct {
		byteCount int64
		expected  string
		name      string
	}{
		{123, "123.00B", "in bytes"},
		{1024, "1.00KiB", "kibibyte"},
		{1024 + 512, "1.50KiB", "fractional kibibytes"},
		{1024*1024 + 1024*512, "1.50MiB", "mebibytes"},
		{1342177280, "1.25GiB", "gibibytes"},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			actual := formatHumanReadableBytes(testCase.byteCount)
			if actual != testCase.expected {
				t.Fatalf("Expected %s, got %s", testCase.expected, actual)
			}
		})
	}
}
