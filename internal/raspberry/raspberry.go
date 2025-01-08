package raspberry

import (
	"fmt"
	"home_control_hub/internal/utils"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
)

func Restart(c *gin.Context) {
	if IsRaspberryPi() {
		utils.Sueecss(c, "重启成功", map[string]interface{}{})
		go func() {
			cmd := exec.Command("reboot")
			cmd.Run()
		}()
		return
	} else {
		utils.Error(c, "重启树莓派失败", "不是树莓派")
	}

}

func Shutdown(c *gin.Context) {
	if IsRaspberryPi() {
		utils.Sueecss(c, "关闭成功", map[string]interface{}{})
		go func() {
			cmd := exec.Command("shutdown", "-h", "now")
			cmd.Run()
		}()
	} else {
		utils.Json(c, "1", "关闭树莓派失败", "不是树莓派", map[string]interface{}{})
	}
}

// 判断是否树莓派
func IsRaspberryPi() bool {
	cmd := exec.Command("uname", "-m")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error running uname:", err)
		return false
	}

	arch := strings.TrimSpace(string(output))
	fmt.Println("CPU Architecture:", arch)

	if arch == "armv7l" || arch == "aarch64" {
		fmt.Println("Running on a Raspberry Pi")
		return true
	} else {
		fmt.Println("Not running on a Raspberry Pi")
		return false
	}
}
