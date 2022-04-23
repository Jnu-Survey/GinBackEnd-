package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"wechatGin/common"
)

// VerifyToken 验证Token是否在缓存中
func VerifyToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		// todo 拿到Token
		urlToken := c.Query("token")
		bodyToken := c.DefaultPostForm("token", "")
		if urlToken == "" && bodyToken == "" {
			ResponseError(c, 10001, errors.New("Token不存在"))
			c.Abort()
			return
		}
		temp := ""
		if urlToken != "" {
			temp = urlToken
		} else {
			temp = bodyToken
		}
		// todo 找缓存
		res, err := common.RedisConfDo("GET", temp)
		if err != nil {
			ResponseError(c, 10002, errors.New("缓存错误"))
			c.Abort()
			return
		}
		if res == nil {
			ResponseError(c, 10004, errors.New("Token已经过期"))
			c.Abort()
			return
		}
		// todo 处理并加入上下文
		uid, ok := res.([]byte)
		if !ok {
			ResponseError(c, 10003, errors.New("服务器错误"))
			c.Abort()
			return
		}
		c.Set("uid", uid) // 能够匹配上则设置上下问信息
		c.Next()
	}
}
