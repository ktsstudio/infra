package encrypt

import (
	"path/filepath"
	"testing"

	"gotest.tools/v3/assert"
)

func TestSealAndUnseal(t *testing.T) {
	tmp := t.TempDir()
	rootKeyPath := filepath.Join(tmp, "root-key")
	assert.NilError(t, CreateRootKey(rootKeyPath))

	key, err := CreateDataKey(rootKeyPath)
	assert.NilError(t, err)

	secretMessage := "This is the message"
	encrypted, err := Seal(key, []byte(secretMessage))
	assert.NilError(t, err)

	unsealed, err := Unseal(key, encrypted)
	assert.NilError(t, err)
	assert.Equal(t, secretMessage, string(unsealed))
}

func TestDecryptDataKey(t *testing.T) {
	tmp := t.TempDir()
	rootKeyPath := filepath.Join(tmp, "root-key")
	assert.NilError(t, CreateRootKey(rootKeyPath))

	dataKey, err := CreateDataKey(rootKeyPath)
	assert.NilError(t, err)

	actual, err := DecryptDataKey(rootKeyPath, dataKey.Encrypted)
	assert.NilError(t, err)

	assert.DeepEqual(t, actual.unencrypted, dataKey.unencrypted)
}

func TestPayloadMetadataLengthLimit(t *testing.T) {
	for _, tc := range []struct {
		name string
		set  func(*encryptedPayload, int)
	}{
		{"algorithm", func(p *encryptedPayload, n int) { p.Algorithm = string(make([]byte, n)) }},
		{"key ID", func(p *encryptedPayload, n int) { p.KeyID = make([]byte, n) }},
		{"root key ID", func(p *encryptedPayload, n int) { p.RootKeyID = string(make([]byte, n)) }},
		{"nonce", func(p *encryptedPayload, n int) { p.Nonce = make([]byte, n) }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			payload := &encryptedPayload{Ciphertext: []byte("message"), Algorithm: "aesgcm", KeyID: []byte("key"), RootKeyID: "root", Nonce: []byte("nonce")}
			tc.set(payload, 255)
			encoded, err := marshalPayload(payload)
			assert.NilError(t, err)
			var decoded encryptedPayload
			assert.NilError(t, unmarshalPayload(encoded, &decoded))
			assert.DeepEqual(t, payload, &decoded)
			tc.set(payload, 256)
			_, err = marshalPayload(payload)
			assert.ErrorContains(t, err, "wire-format length limit")
		})
	}
}
