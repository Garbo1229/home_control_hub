package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code string                 `json:"code"`
	Msg  string                 `json:"msg"`
	Err  string                 `json:"err"`
	Data map[string]interface{} `json:"data"`
}

// 失败Json
func Error(c *gin.Context, msg string, err string) {
	Json(c, "1", msg, err, map[string]interface{}{})
}

// 成功Json
func Sueecss(c *gin.Context, msg string, data map[string]interface{}) {
	Json(c, "0", msg, "", data)
}

func Json(c *gin.Context, code string, msg string, err string, data map[string]interface{}) {
	// 默认响应
	defaultResponse := Response{
		Code: "1",
		Msg:  "",
		Err:  "",
		Data: map[string]interface{}{},
	}

	// 更新响应结构体中的值
	if code != "" {
		defaultResponse.Code = code
	}
	if msg != "" {
		defaultResponse.Msg = msg
	}
	if err != "" {
		defaultResponse.Err = err
	}
	if data != nil {
		defaultResponse.Data = data
	}

	// 返回 JSON 响应
	c.JSON(http.StatusOK, defaultResponse)
}
