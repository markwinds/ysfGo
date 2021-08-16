package cmd

import "os/exec"

// RunCmd 执行命令 返回命令执行结构字符串 传入命令和参数
func RunCmd(cmdStr string, args []string) (string, error) {
	cmd := exec.Command(cmdStr, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(out), nil
}
