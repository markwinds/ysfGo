package main

import (
	"flag"
	"fmt"
	"time"
	"ysf/utils"
)

// 数字键盘的位置
var numX = []uint{551, 237, 530, 902, 158, 554, 895, 188, 550, 890}
var numY = []uint{2086, 1673, 1681, 1676, 1837, 1814, 1816, 1950, 1950, 1950}

// 支付密码
var password = "159632"

// 按键整个过程坐标
var ysfX = []uint{663, 891, 310, 160, 955, 536, 206, 576, 902, 920, 879, 547, 1014}
var ysfY = []uint{332, 1878, 448, 419, 206, 2074, 1659, 1831, 1961, 1834, 1688, 1689, 118}

// 设备名
var device = "beff28c7"
var app = "unionpay"

var aim = 1

func init() {
	for i, item := range password {
		ysfX[i+6] = numX[item-'0']
		ysfY[i+6] = numY[item-'0']
	}
}

func main() {
	flag.IntVar(&aim, "n", 1, "num of ysf")
	flag.Parse()

	adb := utils.NewAdb("adb")
	err := adb.ConnectDev(device, false)
	if err != nil {
		panic(err)
	}

	// 启动程序
	_, err = adb.RunCmdToDevice(device, []string{"shell", "am start -n com.unionpay/.activity.UPActivityMain"})
	if err != nil {
		panic(err)
	}
	// 关闭振动
	_, err = adb.RunCmdToDevice(device, []string{"shell", "cmd appops set com.unionpay VIBRATE ignore"})
	if err != nil {
		panic(err)
	}

	errNum := 0
	num := 0
	for num < aim {
		fmt.Println("---------------------")
		fmt.Println(time.Now())
		fmt.Printf("num[%d] err[%d]", num+1, errNum)
		num += 1
		// 开始点击
		for i, item := range ysfX {
			err := adb.Clink(device, item, ysfY[i])
			if err != nil {
				panic(err)
			}
			// 点击完等待
			switch i {
			case 4:
				time.Sleep(2000 * time.Millisecond)
			case 6, 7, 8, 9, 10:
				time.Sleep(100 * time.Millisecond)
			case 11:
				time.Sleep(1500 * time.Millisecond)
			default:
				time.Sleep(500 * time.Millisecond)
			}
		}
		// 错误检测
		name, err := adb.GetTopActivity(device, app)
		if err != nil {
			panic(err)
		}
		delay := 2000
		for name != ".activity.UPActivityMain" {
			// 出现错误就重启app
			_, err := adb.RunCmdToDevice(device, []string{"shell", "am force-stop  com.unionpay"})
			if err != nil {
				panic(err)
			}
			time.Sleep(time.Duration(delay) * time.Millisecond)
			_, err = adb.RunCmdToDevice(device, []string{"shell", "am start -n com.unionpay/.activity.UPActivityMain"})
			if err != nil {
				panic(err)
			}
			_, err = adb.RunCmdToDevice(device, []string{"shell", "cmd appops set com.unionpay VIBRATE ignore"})
			if err != nil {
				panic(err)
			}
			time.Sleep(time.Duration(delay) * time.Millisecond)
			name, err = adb.GetTopActivity(device, app)
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
