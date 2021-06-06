package main

import (
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

func init() {
	for i, item := range password {
		ysfX[i+6] = numX[item-'0']
		ysfY[i+6] = numY[item-'0']
		fmt.Println(i, item-'0')
	}
}

func main() {
	adb := utils.NewAdb("adb")
	err := adb.ConnectDev(device, false)
	if err != nil {
		panic(err)
	}

	errNum := 0
	num := 0
	aim := 3000
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
				time.Sleep(3000 * time.Millisecond)
			case 6, 7, 8, 9, 10:
				time.Sleep(400 * time.Millisecond)
			case 11:
				time.Sleep(1800 * time.Millisecond)
			default:
				time.Sleep(1000 * time.Millisecond)
			}
		}
		// 错误检测
		name, err := adb.GetTopActivity(device, app)
		if err != nil {
			panic(err)
		}
		if name == ".activity.UPActivityMain" {
			continue
		}
		if name == ".activity.react.UPActivityReactNative" { // 卡在支付界面
			errNum++
			err := adb.Clink(device, 70, 136)
			if err != nil {
				panic(err)
			}
			time.Sleep(1100 * time.Millisecond)
			name, err = adb.GetTopActivity(device, app)
			if err != nil {
				panic(err)
			}
			if name == ".activity.react.UPActivityReactNative" {
				for name == ".activity.react.UPActivityReactNative" {
					err := adb.Clink(device, 997, 139)
					if err != nil {
						panic(err)
					}
					time.Sleep(1100 * time.Millisecond)
					name, err = adb.GetTopActivity(device, app)
					if err != nil {
						panic(err)
					}
				}
			}
			continue
		}
		fmt.Println(name)
		panic("activity error")
	}

}
