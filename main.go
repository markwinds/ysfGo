package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
	"ysf/utils"
)

var Config struct {
	Coordinates []string `yaml:"Coordinates,flow"`
	Delay       int      `yaml:"Delay"`
	Num         int      `yaml:"Num"`
	AdbPath     string   `yaml:"AdbPath"`
	Device      string   `yaml:"Device"`
}

var configPath = "./config.yml"
var app = "com.unionpay"

func main() {
	// 解析参数
	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(err) // 程序要运行一些参数是必要的 没有配置文件直接退出
	}
	err = yaml.Unmarshal(content, &Config)
	if len(Config.Coordinates) == 0 {
		panic("coordinate num is 0")
	}

	// 解析坐标
	var (
		ysfX []uint
		ysfY []uint
	)
	for _, item := range Config.Coordinates {
		c := strings.Split(item, ",")
		if len(c) != 2 {
			panic("coordinate format error")
		}
		x, _ := strconv.Atoi(c[0])
		y, _ := strconv.Atoi(c[1])
		ysfX = append(ysfX, uint(x))
		ysfY = append(ysfY, uint(y))
	}
	pwdIndex := len(ysfX) - 8

	// 初始化adb
	adb := utils.NewAdb(Config.AdbPath)
	err = adb.ConnectDev(Config.Device, false)
	if err != nil {
		panic(err)
	}

	// 启动程序
	_, err = adb.RunCmdToDevice(Config.Device, []string{"shell", "am start -n com.unionpay/.activity.UPActivityMain"})
	if err != nil {
		panic(err)
	}
	// 关闭振动
	_, err = adb.RunCmdToDevice(Config.Device, []string{"shell", "cmd appops set com.unionpay VIBRATE ignore"})
	if err != nil {
		panic(err)
	}
	// 屏幕常亮
	err = adb.AlwaysLight(Config.Device)
	if err != nil {
		panic(err)
	}

	errNum := 0
	num := 0
	for num < Config.Num {
		fmt.Println("---------------------")
		fmt.Println(time.Now())
		fmt.Printf("num[%d] err[%d]", num+1, errNum)
		num += 1
		// 开始点击
		for i, item := range ysfX {
			err := adb.Clink(Config.Device, item, ysfY[i])
			if err != nil {
				panic(err)
			}
			// 点击完等待
			if i >= pwdIndex && i < pwdIndex+6 {
				time.Sleep(500 * time.Millisecond)
			} else {
				time.Sleep(time.Duration(Config.Delay) * time.Millisecond)
			}
		}
		// 错误检测
		activity, err := adb.GetTopActivity(Config.Device, app)
		if err != nil {
			panic(err)
		}
		fragment, err := adb.GetTopFragment(Config.Device, app)
		if err != nil {
			panic(err)
		}
		delay := 2000
		for activity != ".activity.UPActivityMain" || fragment != "0" {
			// 出现错误就重启app
			_, err := adb.RunCmdToDevice(Config.Device, []string{"shell", "am force-stop  com.unionpay"})
			if err != nil {
				panic(err)
			}
			time.Sleep(time.Duration(delay) * time.Millisecond)
			_, err = adb.RunCmdToDevice(Config.Device, []string{"shell", "am start -n com.unionpay/.activity.UPActivityMain"})
			if err != nil {
				panic(err)
			}
			_, err = adb.RunCmdToDevice(Config.Device, []string{"shell", "cmd appops set com.unionpay VIBRATE ignore"})
			if err != nil {
				panic(err)
			}
			time.Sleep(time.Duration(delay) * time.Millisecond)
			activity, err = adb.GetTopActivity(Config.Device, app)
			if err != nil {
				panic(err)
			}
			fragment, err = adb.GetTopFragment(Config.Device, app)
			if err != nil {
				panic(err)
			}
			delay += 1000
		}
		if delay != 2000 {
			errNum++
			num--
		}
	}

}
