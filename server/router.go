package server

import (
	"github.com/sugarscat/seetime/server/account"
	"github.com/sugarscat/seetime/server/tasks"

	"github.com/gin-gonic/gin"
)

func OpenRouter() {
	router := gin.Default()

	router.Any("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// 在此处添加您的React应用程序路由处理程序
	router.Static("/static", "./build/static")
	router.NoRoute(func(c *gin.Context) {
		c.File("./build/index.html")
	})

	api := router.Group("/api")
	{
		api.GET("/time", HandleTime)

		api.GET("/login", account.HandleLogin)

		api.GET("/me", account.HandleMe)
		api.PUT("/me", account.HandleMeUpdate)

		api.GET("/user", account.HandleUser)

		api.GET("/users", account.HandleUsers)
		api.POST("/users", account.HandleUsersAdd)
		api.PUT("/users", account.HandleUsersUpdate)
		api.DELETE("/users", account.HandleUsersDelete)

		api.GET("/tasks", tasks.HandleTasks)
		api.POST("/tasks", tasks.HandleTasksAdd)
		api.DELETE("/tasks", tasks.HandleTasksDelete)

		api.GET("/task", tasks.HandleTask)
	}

	router.Run(":6060")
}

func init() {}
