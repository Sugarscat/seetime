package server

import (
	"time"

	"github.com/gin-gonic/gin"
)

func HandleTime(ctx *gin.Context) {
	now := time.Now()                  //获取当前时间
	timestamp := now.Unix()            //时间戳
	timeObj := time.Unix(timestamp, 0) //将时间戳转为时间格式
	response := gin.H{
		"time": timeObj,
	}
	ctx.JSON(200, response)
}
