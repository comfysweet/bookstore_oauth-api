package rest

import (
	"encoding/json"
	"github.com/comfysweet/bookstore_oauth-api/domain/users"
	"github.com/comfysweet/bookstore_oauth-api/utils/errors"
	"github.com/mercadolibre/golang-restclient/rest"
	"time"
)

var userRestClient = rest.RequestBuilder{
	Timeout: 100 * time.Millisecond,
	BaseURL: "http://localhost:8081",
}

type RestUsersRepository interface {
	LoginUser(string, string) (*users.User, *errors.RestErr)
}

func NewRepository() RestUsersRepository {
	return &repository{}
}

type repository struct {
}

func (repo *repository) LoginUser(email string, password string) (*users.User, *errors.RestErr) {
	request := users.UserLoginRequest{
		Email:    email,
		Password: password,
	}
	response := userRestClient.Post("/users/login", request)
	if response == nil || response.Response == nil {
		return nil, errors.NewInternalServiceError("invalid rest client response when trying to login user")
	}
	if response.StatusCode > 299 {
		var restErr errors.RestErr
		err := json.Unmarshal(response.Bytes(), &restErr)
		if err != nil {
			return nil, errors.NewInternalServiceError("invalid error interface")
		}
		return nil, &restErr
	}

	var user users.User
	if err := json.Unmarshal(response.Bytes(), &user); err != nil {
		return nil, errors.NewInternalServiceError("error when unmarshal users response")
	}
	return &user, nil
}
