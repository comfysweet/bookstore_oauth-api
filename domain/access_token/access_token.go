package access_token

import (
	"fmt"
	"github.com/comfysweet/bookstore_oauth-api/utils/crypto_utils"
	"github.com/comfysweet/bookstore_oauth-api/utils/errors"
	"strings"
	"time"
)

const (
	expirationTime             = 24
	grantTypePassword          = "password"
	grantTypeClientCredentials = "client_credentials"
)

type AccessToken struct {
	AccessToken string `json:"access_token"`
	UserId      int64  `json:"user_id"`
	ClientId    int64  `json:"client_id"`
	Expires     int64  `json:"expires"`
}

type AccessTokenRequest struct {
	GrantType string `json:"grant_type"`
	Scope     string `json:"scope"`

	// Used for password grant type
	Username string `json:"username"`
	Password string `json:"password"`

	// Used for client credentials grant type
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func (accessToken *AccessTokenRequest) Validate() *errors.RestErr {
	switch accessToken.GrantType {
	case grantTypePassword:
		break
	case grantTypeClientCredentials:
		break
	default:
		return errors.NewBadRequestError("invalid grant_type parameter")
	}
	return nil
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

func GetNewAccessToken(userId int64) AccessToken {
	return AccessToken{
		UserId:  userId,
		Expires: time.Now().UTC().Add(expirationTime * time.Hour).Unix()}
}
func (accessToken *AccessToken) isExpired() bool {
	return time.Unix(accessToken.Expires, 0).Before(time.Now().UTC())
}

func (accessToken *AccessToken) Generate() {
	accessToken.AccessToken = crypto_utils.GetMd5(fmt.Sprintf("at-%d-%d-ran", accessToken.UserId, accessToken.Expires))
}
