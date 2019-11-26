package routers

import (
	"github.com/gin-gonic/gin"
	"html_parse_api/actions"
)

var r *gin.Engine

func init() {
	r = gin.Default()
}

func Load() *gin.Engine {
	r.POST("html-parse/baidu-pc", actions.ParseBaiduPc)
	r.POST("html-parse/so", actions.ParseSo)
	r.POST("html-parse/sogou", actions.ParseSogou)

	return r
}
