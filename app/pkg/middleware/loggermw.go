package mw

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LoggerMiddleware helps to tidy log format & create zap.logfile
func LoggerMiddleware(g *gin.Engine) {
	g.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		message := fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
		zap.L().Info(message)
		return message
	}))
}
