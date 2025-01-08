package nas

import (
	"encoding/xml"
	"errors"
	"home_control_hub/config"
	"home_control_hub/internal/utils"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

var NasConfig = config.GlobalConfig.NasConfig

var c = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")

// ErrAuthFailed is the error returned by Shutdown or
// other functions  to indicate an authentication error.
var ErrAuthFailed = errors.New("authentication failed")

func Wake(c *gin.Context) {
	// 目标 MAC 地址和广播地址
	isWake := utils.WakeOnLanHandler(NasConfig.Mac)
	if !isWake {
		utils.Error(c, "Nas开机失败", "")
		return
	}
	utils.Sueecss(c, "Nas开机成功", map[string]interface{}{})
}

// 模拟登录
func login(baseUrl, user, password string) (sid string, err error) {
	var reqUrl = baseUrl + "/cgi-bin/authLogin.cgi"
	var formData = url.Values{}

	formData.Add("user", user)
	formData.Add("pwd", encodePwd(password))
	resp, err := http.Post(reqUrl, "application/x-www-form-urlencoded", strings.NewReader(formData.Encode()))
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", errors.New(resp.Status)
	}
	var xmlResp = struct {
		AuthPassed bool   `xml:"authPassed"`
		AuthSID    string `xml:"authSid"`
	}{}
	if err = xml.NewDecoder(resp.Body).Decode(&xmlResp); err != nil {
		return "", err
	}
	if !xmlResp.AuthPassed {
		return "", ErrAuthFailed
	}
	return xmlResp.AuthSID, nil
}

// 关机
func Shutdown(c *gin.Context) {
	baseUrl := NasConfig.Url
	user := NasConfig.Username
	password := NasConfig.Password
	sid, err := login(baseUrl, user, password)
	if err != nil {
		utils.Error(c, "关闭失败", err.Error())
		return
	}
	var reqUrl = baseUrl + "/cgi-bin/sys/sysRequest.cgi?sid=" + url.QueryEscape(sid) + "&subfunc=power_mgmt&apply=shutdown"
	resp, err := http.Post(reqUrl, "", nil)
	if err != nil {
		utils.Error(c, "关闭失败", err.Error())
		return
	}

	if resp.StatusCode != http.StatusOK {
		utils.Error(c, "关闭失败", errors.New(resp.Status).Error())
		return
	}

	var xmlResp = struct {
		AuthPassed bool `xml:"authPassed"`
	}{}
	if err = xml.NewDecoder(resp.Body).Decode(&xmlResp); err != nil {
		utils.Error(c, "关闭失败", err.Error())
		return
	}
	if !xmlResp.AuthPassed {
		utils.Error(c, "关闭失败", ErrAuthFailed.Error())
		return
	}
	utils.Json(c, "0", "关闭成功", "", map[string]interface{}{})

}

/*
encodePwd: from QNAPTool.ezEncode:
*/
func encodePwd(in string) string {
	var C = []rune(in)
	var y []rune

	e := len(C)
	A := 0
	B := rune(0)
	for A < e {
		B = C[A] & 255
		A++
		if A == e {
			y = append(y, c[B>>2], c[(B&3)<<4])
			y = append(y, '=', '=')
			break
		}
		z := C[A]
		A++
		if A == e {
			y = append(y, c[B>>2], c[((B&3)<<4)|((z&240)>>4)], c[(z&15)<<2], '=')
			break
		}
		x := C[A]
		A++
		y = append(y, c[B>>2], c[((B&3)<<4)|((z&240)>>4)], c[((z&15)<<2)|((x&192)>>6)], c[x&63])
	}
	return string(y)
}
