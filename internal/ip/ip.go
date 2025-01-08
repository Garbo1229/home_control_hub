package ip

import (
	"fmt"
	"home_control_hub/config"
	"home_control_hub/internal/utils"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
)

// 配置
var Ip2regionConfig = config.GlobalConfig.Ip2regionConfig

// 查询IP库
func QueryIp(c *gin.Context) {
	clientIP := c.Param("ip") // 获取路由参数
	if clientIP == "" {
		clientIP = c.ClientIP()
	}

	searcher, err := xdb.NewWithFileOnly(Ip2regionConfig.DbPath)
	if err != nil {
		errStr := fmt.Sprintf("创建搜索器失败: %s", err.Error())
		fmt.Println("[ip2region]创建搜索器失败:", errStr)
		utils.Error(c, "配置出错", errStr)
		return
	}

	defer searcher.Close()

	region, err := searcher.SearchByStr(clientIP)
	if err != nil {
		errStr := fmt.Sprintf("搜索IP失败(%s): %s", clientIP, err)
		fmt.Println("搜索IP失败:", errStr)
		utils.Error(c, "搜索IP失败", errStr)
		return
	}

	parts := strings.Split(region, "|")
	for key, part := range parts {
		if part == "0" {
			parts[key] = ""
		}
	}
	data := map[string]interface{}{
		"ip":       clientIP,
		"country":  parts[0],
		"region":   parts[1],
		"province": parts[2],
		"city":     parts[3],
		"isp":      parts[4],
	}
	utils.Sueecss(c, "查询成功", data)
}

// 更新IP数据库
func Update(c *gin.Context) {
	// 立即返回响应，表示请求已接收并正在处理
	utils.Sueecss(c, "已经正在更新", map[string]interface{}{})
	// 启动 goroutine 异步执行命令
	go func() {
		// 获取当前工作目录
		root := utils.GetPwd()

		// 确保目标文件夹存在
		assetsPath := root + "/assets"
		if _, err := os.Stat(assetsPath); os.IsNotExist(err) {
			if err := os.Mkdir(assetsPath, 0775); err != nil {
				fmt.Println("[ip2region]创建目标文件夹失败:", err)
				return
			}
		}
		// 判断是否存在文件
		ip2regionMakerPath :=
			Ip2regionConfig.MakerPath + "/xdb_maker"

		_, err := os.Stat(ip2regionMakerPath)
		if os.IsNotExist(err) {
			fmt.Println(ip2regionMakerPath + " 不存在")
			return
		}

		// 执行 IP2Region 数据库更新操作
		ip2regionMaker := exec.Command(ip2regionMakerPath, "gen",
			"--src="+Ip2regionConfig.Path+"/data/ip.merge.txt",
			"--dst="+Ip2regionConfig.DbPath+"_new")

		// 捕获命令输出和错误
		output, err := ip2regionMaker.CombinedOutput()
		if err != nil {
			fmt.Println("[ip2region]命令执行时出错:", err)
			fmt.Println("[ip2region]命令输出:", string(output))
			return
		}

		fmt.Println("[ip2region]命令执行成功，输出:", string(output))

		// 等待新文件的生成
		newFilePath := Ip2regionConfig.DbPath + "_new"
		timeout := time.After(30 * time.Second)
		tick := time.Tick(500 * time.Millisecond)

	loop:
		for {
			select {
			case <-timeout:
				fmt.Println("[ip2region]等待新文件生成超时")
				return
			case <-tick:
				if _, err := os.Stat(newFilePath); err == nil {
					break loop
				}
			}
		}

		// 判断目标文件是否存在
		_, err = os.Stat(Ip2regionConfig.DbPath)
		if !os.IsNotExist(err) {
			// 删除目标文件
			if err := os.Remove(Ip2regionConfig.DbPath); err != nil && !os.IsNotExist(err) {
				fmt.Println("[ip2region]删除已有文件失败:", err)
				return
			}
		}

		// 重命名文件
		if err := os.Rename(newFilePath, Ip2regionConfig.DbPath); err != nil {
			fmt.Println("[ip2region]重命名文件失败:", err)
			return
		}

		fmt.Println("[ip2region]文件处理完成")
	}()
}
