package utils

import (
	"bytes"
	"os"
	"os/exec"
	"runtime"
)


func StartServer() {
	switch runtime.GOOS {
	case "windows", "linux", "darwin":
		{
			var out bytes.Buffer
			cmd := exec.Command("node", "../sign/index.js")
			cmd.Stdout = &out  // 把执行命令的标准输出定向到out
			cmd.Stderr = os.Stderr // 把命令的错误输出定向到 os
		     cmd.Start()
		}
	}
}
