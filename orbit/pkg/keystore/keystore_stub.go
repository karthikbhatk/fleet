//go:build !darwin || !cgo

package keystore

import "errors"

func Exists() bool {
	return false
}

func Name() string {
	return "not implemented"
}

// AddSecret will add a secret to the keychain. This secret can be retrieved by this application without any user authorization.
func AddSecret(secret string) error {
	return errors.New("not implemented")
}

// GetSecret will retrieve a secret from the keychain. If the secret was added by user or another application,
// then this application needs to be authorized to retrieve the secret.
func GetSecret() (string, error) {
	return "", errors.New("not implemented")
}
