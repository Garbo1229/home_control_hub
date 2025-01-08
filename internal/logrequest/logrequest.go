package logrequest

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func LogRequest(c *gin.Context) {
	method := c.Request.Method
	path := c.Request.URL.Path
	rawQuery := c.Request.URL.RawQuery
	host := c.Request.Host
	headers := c.Request.Header

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to read request body"})
		return
	}

	// 获取请求的原始地址
	remoteAddr := c.Request.Header.Get("X-Real-IP")
	if remoteAddr == "" {
		remoteAddr = c.ClientIP()
	}

	// 获取通过Cloudflare代理后的实际IP地址
	realIP := c.Request.Header.Get("CF-Connecting-IP")
	if realIP == "" {
		// 如果没有通过Cloudflare，使用X-Forwarded-For头部
		xForwardedFor := c.Request.Header.Get("X-Forwarded-For")
		if xForwardedFor != "" {
			// X-Forwarded-For 可能包含多个IP地址，取第一个
			realIP = strings.Split(xForwardedFor, ",")[0]
		} else {
			// 如果没有X-Forwarded-For头部，则使用RemoteAddr
			realIP = c.ClientIP()
		}
	}

	// 构建完整的URL
	url := fmt.Sprintf("%s://%s%s", c.Request.Header.Get("X-Forwarded-Proto"), host, path)
	if rawQuery != "" {
		url += "?" + rawQuery
	}

	// Log the request details
	fmt.Printf("Method: %s\n", method)
	fmt.Printf("URL: %s\n", url)
	fmt.Printf("Host: %s\n", host)
	fmt.Printf("RemoteAddr: %s\n", remoteAddr)
	fmt.Printf("RealIP: %s\n", realIP)
	fmt.Printf("Headers: %v\n", headers)
	fmt.Printf("Body: %s\n", string(body))

	// Optionally, you can write the logs to a file instead of printing them to the console
	// f, err := os.OpenFile("requests.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// if err != nil {
	//     c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to open log file"})
	//     return
	// }
	// defer f.Close()
	// fmt.SetOutput(f)
	// fmt.Printf("Method: %s\nURL: %s\nHost: %s\nRemoteAddr: %s\nRealIP: %s\nHeaders: %v\nBody: %s\n", method, url, host, remoteAddr, realIP, headers, string(body))

	c.JSON(http.StatusOK, gin.H{"status": "Request logged"})
}
