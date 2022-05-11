package services

import (
	"github.com/comfysweet/bookstore_oauth-api/domain/access_token"
	"github.com/comfysweet/bookstore_oauth-api/repository/db"
	"github.com/comfysweet/bookstore_oauth-api/repository/rest"
	"github.com/comfysweet/bookstore_oauth-api/utils/errors"
	"strings"
)

type Service interface {
	GetById(string) (*access_token.AccessToken, *errors.RestErr)
	Create(access_token.AccessTokenRequest) (*access_token.AccessToken, *errors.RestErr)
	UpdateExpirationTime(access_token.AccessToken) *errors.RestErr
}

type service struct {
	restUserRepo rest.RestUsersRepository
	dbRepo       db.DbRepository
}

func NewService(userRepo rest.RestUsersRepository, dbRepo db.DbRepository) Service {
	return &service{
		restUserRepo: userRepo,
		dbRepo:       dbRepo,
	}
}

func (s *service) GetById(accessTokenId string) (*access_token.AccessToken, *errors.RestErr) {
	accessTokenId = strings.TrimSpace(accessTokenId)
	if len(accessTokenId) == 0 {
		return nil, errors.NewBadRequestError("invalid access token id")
	}
	accessToken, err := s.dbRepo.GetById(accessTokenId)
	if err != nil {
		return nil, err
	}
	return accessToken, nil
}

func (s *service) Create(request access_token.AccessTokenRequest) (*access_token.AccessToken, *errors.RestErr) {
	if err := request.Validate(); err != nil {
		return nil, err
	}
	user, err := s.restUserRepo.LoginUser(request.Username, request.Password)
	if err != nil {
		return nil, err
	}

	// Generate a new access token:
	accessToken := access_token.GetNewAccessToken(user.Id)
	accessToken.Generate()

	// Save the new access token in Cassandra:
	if err := s.dbRepo.Create(accessToken); err != nil {
		return nil, err
	}
	return &accessToken, nil
}

func (s *service) UpdateExpirationTime(accessToken access_token.AccessToken) *errors.RestErr {
	if err := accessToken.Validate(); err != nil {
		return err
	}
	return s.dbRepo.UpdateExpirationTime(accessToken)
}
