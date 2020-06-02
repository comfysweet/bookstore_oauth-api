package http

import (
	"github.com/comfysweet/bookstore_oauth-api/domain/access_token"
	"github.com/comfysweet/bookstore_oauth-api/services"
	"github.com/comfysweet/bookstore_oauth-api/utils/errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AccessTokenHandler interface {
	GetById(*gin.Context)
	Create(*gin.Context)
}

type accessTokenHandler struct {
	service services.Service
}

func NewAccessTokenHandler(service services.Service) AccessTokenHandler {
	return &accessTokenHandler{service: service,}
}

func (handler *accessTokenHandler) GetById(c *gin.Context) {
	accessToken, err := handler.service.GetById(c.Param("access_token_id"))
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, accessToken)
}

func (handler *accessTokenHandler) Create(c *gin.Context) {
	var accessToken access_token.AccessToken
	if err := c.ShouldBindJSON(&accessToken); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}
	if err := handler.service.Create(accessToken); err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, accessToken)
}
