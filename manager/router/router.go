package router

import (
	"log"

	"github.com/dollarkillerx/smart-dns-go/manager/pkg/model"
	"github.com/dollarkillerx/smart-dns-go/manager/standard"
	"github.com/gin-gonic/gin"
)

func Router(app *gin.Engine) {
	user := app.Group("/user")
	user.POST("/login", func(ctx *gin.Context) {
		var user model.User
		if err := ctx.BindJSON(&user); err != nil {
			ctx.JSON(400, standard.ParamsError)
			return
		}

		log.Println(user)
		ctx.JSON(200, standard.Response{
			Code:    200,
			Success: true,
			Data:    user,
		})
	})
}
