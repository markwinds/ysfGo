package utils

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type Adb struct {
	AdbPath string // 系统中adb命令的执行路径
}

func NewAdb(adbPath string) *Adb {
	return &Adb{AdbPath: adbPath}
}

// RunCmd 传入adb参数 执行adb命令 返回命令执行结果
func (a *Adb) RunCmd(args []string) (string, error) {
	cmd := exec.Command(a.AdbPath, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// ConnectDev 连接到设备
// 输入device为设备名比如198.2.4.6
// root指示是否以root身份在设备运行命令
func (a *Adb) ConnectDev(device string, root bool) error {
	_, err := a.RunCmd([]string{"connect", device})
	if err != nil {
		return err
	}
	if !a.CheckDevice(device) {
		return fmt.Errorf("connect to device[%s] failed", device)
	}
	// 获取root权限
	if root {
		res, err := a.RunCmdToDevice(device, []string{"root"})
		if err != nil {
			return err
		}
		// root的uid为0 根据这个判断是否获取到了root权限
		res, err = a.RunCmdToDevice(device, []string{"id"})
		if err != nil {
			return err
		}
		if !strings.Contains(res, "uid=0") {
			fmt.Printf("get root failed res[%s]", res)
			return fmt.Errorf("get root failed res[%s]", res)
		}
	}
	return nil
}

// RunCmdToDevice 在指定的设备上运行命令
func (a *Adb) RunCmdToDevice(device string, args []string) (string, error) {
	args = append([]string{"-s", device}, args...)
	return a.RunCmd(args)
}

// CheckDevice 检测设备是否已经连接
func (a *Adb) CheckDevice(device string) bool {
	// 通过adb devices命令判断连接是否成功
	res, err := a.RunCmd([]string{"devices"})
	if err != nil {
		fmt.Println("cmd devices failed")
		return false
	}
	if !strings.Contains(res, device) {
		//fmt.Printf("connect device failed res[%s]", res)
		fmt.Printf("device[%s] not connected res[%s]", device, res)
		return false
	}
	return true
}

// Clink 模拟点击 输入设备和坐标
func (a *Adb) Clink(device string, x uint, y uint) error {
	// 运行命令前先确认device已经连接
	if !a.CheckDevice(device) {
		return fmt.Errorf("device[%s] not connected", device)
	}
	_, err := a.RunCmdToDevice(device, []string{"shell", "input", "tap", strconv.Itoa(int(x)), strconv.Itoa(int(y))})
	if err != nil {
		return err
	}
	return nil
}

// GetTopActivity 根据应用名获取当前activity
func (a *Adb) GetTopActivity(device string, app string) (string, error) {
	if !a.CheckDevice(device) {
		return "", fmt.Errorf("device[%s] not connected", device)
	}
	// 组命令
	arg := `dumpsys activity top |grep ACTIVITY |grep ` + app + ` |grep -v '^$' |awk -F '/' '{print $2}'  |awk '{print $1}'`
	res, err := a.RunCmdToDevice(device, []string{"shell", arg})
	if err != nil {
		return "", err
	}
	res = strings.Trim(res, "\r\n")
	res = strings.Trim(res, "\n")
	return res, nil
}
