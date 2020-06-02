package db

import (
	"github.com/comfysweet/bookstore_oauth-api/clients/cassandra"
	"github.com/comfysweet/bookstore_oauth-api/domain/access_token"
	"github.com/comfysweet/bookstore_oauth-api/utils/errors"
	"github.com/gocql/gocql"
)

const (
	queryGetAccessToken    = "select access_token, user_id, client_id, expires from access_tokens where access_token=?;"
	queryCreateAccessToken = "insert into access_tokens(access_token, user_id, client_id, expires) values (?, ?, ?, ?);"
	queryUpdateExpires     = "update access_tokens set expires=? where access_token=?;"
)

func NewRepository() DbRepository {
	return &repository{}
}

type repository struct {
}

type DbRepository interface {
	GetById(string) (*access_token.AccessToken, *errors.RestErr)
	Create(access_token.AccessToken) *errors.RestErr
	UpdateExpirationTime(access_token.AccessToken) *errors.RestErr
}

func (repo *repository) GetById(id string) (*access_token.AccessToken, *errors.RestErr) {
	var result = access_token.AccessToken{}
	if err := cassandra.GetSession().Query(queryGetAccessToken, id).Scan(&result.AccessToken, &result.UserId, &result.ClientId, &result.Expires); err != nil {
		if err == gocql.ErrNotFound {
			return nil, errors.NewNotFoundError("no access token found by given id")
		}
		return nil, errors.NewInternalServiceError(err.Error())
	}
	return &result, nil
}

func (repo *repository) Create(token access_token.AccessToken) *errors.RestErr {
	if err := cassandra.GetSession().Query(queryCreateAccessToken, token.AccessToken, token.UserId, token.ClientId, token.Expires).Exec(); err != nil {
		return errors.NewInternalServiceError(err.Error())
	}
	return nil
}

func (repo *repository) UpdateExpirationTime(token access_token.AccessToken) *errors.RestErr {
	if err := cassandra.GetSession().Query(queryUpdateExpires, token.Expires, token.AccessToken).Exec(); err != nil {
		return errors.NewInternalServiceError(err.Error())
	}
	return nil
}
