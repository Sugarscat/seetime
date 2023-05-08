package server

import (
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sugarscat/seetime/server/account"
)

func HandleSystem(ctx *gin.Context) {
	token := ctx.Request.Header.Get("Authorization")
	success, id := account.ChecKToken(token)
	// 获取时区
	name, offset := time.Now().Zone()
	fixedZone := time.FixedZone(name, offset)
	local := time.Now().In(fixedZone)

	if success {
		if account.ParsingPermissions(id, "situation") {
			ctx.JSON(200, gin.H{
				"code":    200,
				"success": true,
				"message": "加载成功",
				"data": gin.H{
					"os":   runtime.GOOS,
					"arch": runtime.GOARCH,
					"time": time.Now().Unix(),
					"utc":  local.Format("-07:00"),
				},
			})
		} else {
			ctx.JSON(200, gin.H{
				"code":    400,
				"success": false,
				"message": "无权限",
				"data":    nil,
			})
		}
		return
	}

	ctx.JSON(200, gin.H{
		"code":    403,
		"success": false,
		"message": "身份令牌过期，请重新登录",
		"data":    nil,
	})
}
