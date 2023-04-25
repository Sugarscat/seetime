package server

import (
	"seetime/server/account"

	"github.com/gin-gonic/gin"
)

func openAPi() {
	r := gin.Default()

	r.Any("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/api/login", account.HandleLogin)

	r.GET("/api/me", account.HandleMe)
	r.PUT("/api/me", account.HandleMeUpdate)
	r.GET("/api/me/info", account.HandleMeInfo)

	r.GET("/api/users", account.HandleUsers)
	r.Run(":6060")
}

func Loading() {
	SendInfo()
	defer openAPi()

}

func init() {}
