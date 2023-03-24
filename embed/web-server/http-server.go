package web_server

import (
	"blog-go/pkg/config"
	"blog-go/pkg/handler"
	"blog-go/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)
func StartWebServer() {
	r:=gin.Default()
	r.Use(Cors())
	Register(r)
	address := config.LoadConfig()
	go func() {
		err := r.Run(address)
		if err != nil {
			logger.ErrorLogger.WithField("Http Start err: ", err).WithError(err)
			panic(err)
		}
	}()
}
func Register(engine *gin.Engine)  {
	chat:=engine.Group("/chat")
	chat.GET("/process",handler.VerifyData)
	chat.POST("/process",handler.MessageHandler)

}
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Set("content-type", "application/json")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}