package access_token

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAccessTokenConstants(t *testing.T) {
	assert.EqualValues(t, 24, expirationTime, "expiration time should be 24 hours")
}

func TestGetNewAccessToken(t *testing.T) {
	token := GetNewAccessToken(0)
	assert.False(t, token.isExpired(), "new access token should'n be expired")
	assert.EqualValues(t, "", token.AccessToken, "new access token should'n have defined access token id")
	assert.True(t, token.UserId == 0, "new access token should'n have associated user id")
	assert.True(t, token.ClientId == 0, "new access token should'n have associated client id")
}

func TestAccessTokenIsExpired(t *testing.T) {
	token := AccessToken{}
	assert.True(t, token.isExpired(), "empty access token should be expired by default")

	token.Expires = time.Now().UTC().Add(3 * time.Hour).Unix()
	assert.False(t, token.isExpired(), "access token expiring three hours from now should not be expired")
}
