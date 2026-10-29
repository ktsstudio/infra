package data

import (
	"crypto/ed25519"
	"crypto/rand"
	"testing"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
	"gotest.tools/v3/assert"

	"github.com/infrahq/infra/internal/claims"
	"github.com/infrahq/infra/internal/server/models"
)

// Existing installations store ED25519 key metadata; upgrades must continue to
// issue tokens accepted by both old and new connectors without rotating keys.
func TestCreateJWTKeyAlgorithmCompatibility(t *testing.T) {
	public, private, err := ed25519.GenerateKey(rand.Reader)
	assert.NilError(t, err)
	for _, algorithm := range []string{"ED25519", "EdDSA"} {
		t.Run(algorithm, func(t *testing.T) {
			key, err := (&jose.JSONWebKey{Key: private, Algorithm: algorithm, Use: "sig"}).MarshalJSON()
			assert.NilError(t, err)
			token, err := createJWT(nil, &models.Organization{PrivateJWK: models.EncryptedAtRest(key)}, &models.Identity{Name: "user@example.com"}, []string{"edit"}, time.Now().Add(time.Minute))
			assert.NilError(t, err)
			parsed, err := jwt.ParseSigned(token, []jose.SignatureAlgorithm{jose.EdDSA})
			assert.NilError(t, err)
			var actual claims.Custom
			assert.NilError(t, parsed.Claims(&jose.JSONWebKey{Key: public, Algorithm: algorithm, Use: "sig"}, &actual))
			assert.Equal(t, actual.Name, "user@example.com")
			assert.DeepEqual(t, actual.Groups, []string{"edit"})
		})
	}
}
