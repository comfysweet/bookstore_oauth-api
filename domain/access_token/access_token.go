package access_token

import (
	"github.com/comfysweet/bookstore_oauth-api/utils/errors"
	"strings"
	"time"
)

const expirationTime = 24

type AccessToken struct {
	AccessToken string `json:"access_token"`
	UserId      int64  `json:"user_id"`
	ClientId    int64  `json:"client_id"`
	Expires     int64  `json:"expires"`
}

func (accessToken *AccessToken) Validate() *errors.RestErr {
	accessToken.AccessToken = strings.TrimSpace(accessToken.AccessToken)
	if len(accessToken.AccessToken) == 0 {
		return errors.NewBadRequestError("invalid access token id")
	}
	if accessToken.UserId <= 0 {
		return errors.NewBadRequestError("invalid user id")
	}
	if accessToken.ClientId <= 0 {
		return errors.NewBadRequestError("invalid client id")
	}
	if accessToken.Expires <= 0 {
		return errors.NewBadRequestError("invalid expiration time")
	}
	return nil
}

func GetNewAccessToken() AccessToken {
	return AccessToken{Expires: time.Now().UTC().Add(expirationTime * time.Hour).Unix()}
}
func (accessToken *AccessToken) isExpired() bool {
	return time.Unix(accessToken.Expires, 0).Before(time.Now().UTC())
}
