package server

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/fleetdm/fleet/v4/pkg/fleethttp"
)

// GenerateRandomText return a string generated by filling in keySize bytes with
// random data and then base64 encoding those bytes
func GenerateRandomText(keySize int) (string, error) {
	key := make([]byte, keySize)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(key), nil
}

func httpSuccessStatus(statusCode int) bool {
	return statusCode >= 200 && statusCode <= 299
}

func PostJSONWithTimeout(ctx context.Context, url string, v interface{}) error {
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		return err
	}

	client := fleethttp.NewClient(fleethttp.WithTimeout(30 * time.Second))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to POST to %s: %s, request-size=%d", maskSecretURLParams(url), err, len(jsonBytes))
	}
	defer resp.Body.Close()

	if !httpSuccessStatus(resp.StatusCode) {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("error posting to %s: %d. %s", maskSecretURLParams(url), resp.StatusCode, string(body))
	}

	return nil
}

// maskSecretURLParams masks URL parameters that match a limited set of keywords (comprising "secret", "token", "key", "password"). It accepts a raw
// string and returns a redacted string if the raw string is URL-parseable. If it is not
// URL-parseable, the raw string is returned unchanged.
func maskSecretURLParams(rawURL string) string {
	matchKeyword := func(k string) bool {
		kws := []string{"secret", "token", "key", "password"}
		for _, kw := range kws {
			if strings.Contains(strings.ToLower(k), kw) {
				return true
			}
		}
		return false
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}

	q := u.Query()
	for k := range q {
		if matchKeyword(k) {
			q[k] = []string{"MASKED"}
		}
	}
	u.RawQuery = q.Encode()

	return u.Redacted()
}

// TODO: Consider moving other crypto functions from server/mdm/apple/util to here

// DecodePrivateKeyPEM decodes PEM-encoded private key data.
func DecodePrivateKeyPEM(encoded []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(encoded)
	if block == nil {
		return nil, errors.New("no PEM-encoded data found")
	}
	if block.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("unexpected block type %s", block.Type)
	}

	return x509.ParsePKCS1PrivateKey(block.Bytes)
}
