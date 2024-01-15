//go:build darwin && cgo

package keystore

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestExists(t *testing.T) {
	t.Parallel()
	assert.True(t, Exists())
}

func TestName(t *testing.T) {
	t.Parallel()
	assert.True(t, strings.Contains(Name(), "keychain"))
}

func TestSecret(t *testing.T) {
	t.Parallel()

	// Use a different service name for testing
	origServiceStringRef := serviceStringRef
	serviceStringRef = stringToCFString("com.fleetdm.fleetd.enroll.secret.test")

	t.Cleanup(
		func() {
			// Delete test secret from keychain, and deallocate memory.
			_ = deleteSecret()
			releaseCFString(serviceStringRef)
			serviceStringRef = origServiceStringRef
		},
	)

	// Make sure secret doesn't exist
	_ = deleteSecret()

	// Get secret -- should be empty
	result, err := GetSecret()
	require.NoError(t, err)
	assert.Equal(t, "", result)

	// Add secret
	secret := "testSecret"
	require.NoError(t, AddSecret(secret))
	result, err = GetSecret()
	require.NoError(t, err)
	assert.Equal(t, secret, result)

	// Update secret
	secret = "updatedSecret"
	require.NoError(t, UpdateSecret(secret))
	result, err = GetSecret()
	require.NoError(t, err)
	assert.Equal(t, secret, result)
}
