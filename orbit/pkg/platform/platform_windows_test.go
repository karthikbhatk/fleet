//go:build windows
// +build windows

package platform

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUUIDNotPresent(t *testing.T) {
	uuidBytes := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	_, err := isValidUUID(uuidBytes)

	assert.NotNil(t, err, "UUID not present: ")
}

func TestUUIDNotSet(t *testing.T) {
	uuidBytes := []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
	_, err := isValidUUID(uuidBytes)

	assert.NotNil(t, err, "UUID not set: ")
}

func TestUUIDInvalidSize(t *testing.T) {
	uuidBytes := []byte{0x11, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
	_, err := isValidUUID(uuidBytes)

	assert.NotNil(t, err, "UUID validation size: ")
}

func TestUUIDValid(t *testing.T) {
	uuidBytes := []byte{0x11, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x00, 0x00, 0x00, 0x00}
	_, err := isValidUUID(uuidBytes)

	assert.Nil(t, err, "UUID validation error: ")
}
