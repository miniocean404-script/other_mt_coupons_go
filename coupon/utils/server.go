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
			cmd.Stdout = &out
			cmd.Stderr = os.Stderr
			cmd.Run()
		}
	}
}
