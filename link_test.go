package main

import (
	"testing"
)

func TestHumanSizeInBytes(t *testing.T) {
	var size int64 = 512
	testHumanReadable(t, size, "512B")
}

func TestHumanSizeInKilobnytes(t *testing.T) {
	var size int64 = 1024 * 1.5
	testHumanReadable(t, size, "1.50KiB")
}

func TestHumanSizeInMegabytes(t *testing.T) {
	var size int64 = 1024 * 1024 * 1.5
	testHumanReadable(t, size, "1.50MiB")
}

func TestHumanSizeInGigabytes(t *testing.T) {
	var size int64 = 1024 * 1024 * 1024 * 1.5
	testHumanReadable(t, size, "1.50GiB")
}

func testHumanReadable(t *testing.T, size int64, expected string) {
	result := HumanSize(size)
	if result != expected {
		t.Fatalf("Expected \"%s\", got \"%s\"", expected, result)
	}
}
