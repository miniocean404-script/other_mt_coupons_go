package utils

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
)


func StartServer() {
	switch runtime.GOOS {
	case "windows", "linux", "darwin":
		{
			cmd := exec.Command("node", "../sign/index.js")

			var out bytes.Buffer  // 也可以输出到 bytes.Buffer 的 out 中
			cmd.Stdout = &out  // 把执行命令的标准输出定向到out
			cmd.Stderr = os.Stderr // 把命令的错误输出定向到 os

		    cmd.Start() // start 异步执行 run 同步阻塞


			for {
				server_result := out.String()
				if server_result!="" {
					fmt.Println(server_result)
					break
				}
			}
		}
	}
}
