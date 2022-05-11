package rest

import (
	"github.com/mercadolibre/golang-restclient/rest"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	rest.StartMockupServer()
	os.Exit(m.Run())
}

func TestLoginUserTimeoutFromApi(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		URL:          "http://localhost:8080/users/login",
		HTTPMethod:   http.MethodPost,
		ReqBody:      `{"email":"test@mail.com", "password":"qwe"}`,
		RespHTTPCode: 0,
		RespBody:     `{}`,
	})

	repository := NewRepository()
	user, err := repository.LoginUser("test@mail.com", "qwe")
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "invalid rest client response when trying to login user", err.Message)
}

func TestLoginUserInvalidErrorInterface(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		URL:          "http://localhost:8080/users/login",
		HTTPMethod:   http.MethodPost,
		ReqBody:      `{"email":"test@mail.com", "password":"qwe"}`,
		RespHTTPCode: http.StatusNotFound,
		RespBody:     `{"message":"invalid login credentials","status":"404"","error":"not_found"}`,
	})

	repository := NewRepository()
	user, err := repository.LoginUser("test@mail.com", "qwe")
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "invalid error interface", err.Message)
}

func TestLoginUserInvalidLoginCredentials(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		URL:          "http://localhost:8080/users/login",
		HTTPMethod:   http.MethodPost,
		ReqBody:      `{"email":"test@mail.com", "password":"qwe"}`,
		RespHTTPCode: http.StatusNotFound,
		RespBody:     `{"message":"invalid login credentials","status":404,"error":"not_found"}`,
	})

	repository := NewRepository()
	user, err := repository.LoginUser("test@mail.com", "qwe")
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusNotFound, err.Status)
	assert.EqualValues(t, "invalid login credentials", err.Message)
}

func TestLoginUserInvalidJsonResponse(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		URL:          "http://localhost:8080/users/login",
		HTTPMethod:   http.MethodPost,
		ReqBody:      `{"email":"test@mail.com", "password":"qwe"}`,
		RespHTTPCode: http.StatusOK,
		RespBody:     `{"id":"1","first_name":"ks","last_name":"zemsk","email":"test@mail.com"}`,
	})

	repository := NewRepository()
	user, err := repository.LoginUser("test@mail.com", "qwe")
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "error when unmarshal users response", err.Message)
}

func TestLoginUserValidJsonResponse(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		URL:          "http://localhost:8080/users/login",
		HTTPMethod:   http.MethodPost,
		ReqBody:      `{"email":"test@mail.com", "password":"qwe"}`,
		RespHTTPCode: http.StatusOK,
		RespBody:     `{"id":1,"first_name":"ks","last_name":"zemsk","email":"test@mail.com"}`,
	})

	repository := NewRepository()
	user, err := repository.LoginUser("test@mail.com", "qwe")
	assert.NotNil(t, user)
	assert.Nil(t, err)
	assert.EqualValues(t, 1, user.Id)
	assert.EqualValues(t, "ks", user.FirstName)
	assert.EqualValues(t, "zemsk", user.LastName)
	assert.EqualValues(t, "test@mail.com", user.Email)
}