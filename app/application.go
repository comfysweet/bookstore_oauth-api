package app

import (
	"github.com/comfysweet/bookstore_oauth-api/http"
	"github.com/comfysweet/bookstore_oauth-api/repository/db"
	"github.com/comfysweet/bookstore_oauth-api/services"
	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
)

func StartApplication() {
	atHandler := http.NewAccessTokenHandler(services.NewService(db.NewRepository()))

	router.GET("/oauth/access_token/:access_token_id", atHandler.GetById)
	router.POST("/oauth/access_token", atHandler.Create)
	router.Run(":8080")
}
