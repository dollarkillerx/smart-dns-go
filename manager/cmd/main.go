package main

import (
	"log"
	"net/http"

	"github.com/dollarkillerx/smart-dns-go/manager/pkg/config"
	"github.com/dollarkillerx/smart-dns-go/manager/router"
	"github.com/gin-gonic/gin"
)

func main() {
	app := gin.New()
	app.Use(gin.Recovery())
	if config.BaseConfig.Debug {
		app.Use(gin.Logger())
	}
	app.Use(Cors())

	router.Router(app)

	if err := app.Run(config.BaseConfig.ListenAddr); err != nil {
		log.Fatalln(err)
	}
}


func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}
