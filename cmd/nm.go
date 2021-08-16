package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// nm命令相关操作

// SymbolInfo lib中的符号信息
type SymbolInfo struct {
	LibName    string // 库名
	Filename   string // 文件名
	Type       string // 符号类型
	SymbolName string // 符号名
}

type Nm struct {
	lib1Symbol map[string]*SymbolInfo // 存放第一个lib的符号信息
	lib2Symbol map[string]*SymbolInfo // 存放其他几个lib的符号信息
}

func NewNm() *Nm {
	return &Nm{
		lib1Symbol: make(map[string]*SymbolInfo),
		lib2Symbol: make(map[string]*SymbolInfo),
	}
}

// InitLib1 加载第一个lib的符号信息到图
func (nm *Nm) InitLib1(libPath string) {
	nm.readLibSymbol(libPath, nm.lib1Symbol)
}

// InitLib2 加载其他lib的符号信息到图
func (nm *Nm) InitLib2(libsPath []string) {
	for _, libPath := range libsPath {
		nm.readLibSymbol(libPath, nm.lib2Symbol)
	}
}

// 读取lib信息到图
func (nm *Nm) readLibSymbol(libPath string, symbolMap map[string]*SymbolInfo) {
	re := regexp.MustCompile(`(.+):(.+):.* (.) (.+)`) // libcJSON.a:cJSON.c.o:0000000000002548 T cJSON_CreateFalse
	// 检查路径是否正确
	_, err := os.Stat(libPath)
	if err != nil {
		fmt.Printf("error:lib path is not exit [%s]\n", libPath)
		panic(err)
	}
	// 执行nm命令 收集symbol
	outStr, err := RunCmd("nm", []string{"-A", libPath})
	if err != nil {
		fmt.Printf("error:nm cmd exec failed! err[%s]\n", err.Error())
		panic(err)
	}
	outStr = strings.Replace(outStr, "\r", "", -1)
	lines := strings.Split(outStr, "\n")
	for _, line := range lines {
		line = strings.Trim(line, " ")
		if len(line) == 0 {
			continue
		}
		res := re.FindStringSubmatch(line)
		if res == nil || len(res) != 5 {
			fmt.Printf("error:get symbol info from line failed! line[%s]\n", line)
			os.Exit(1)
		}
		if res[4][0] == '$' {
			continue
		}
		var symbolInfo SymbolInfo
		symbolInfo.LibName = res[1]
		symbolInfo.Filename = res[2]
		symbolInfo.Type = res[3]
		symbolInfo.SymbolName = res[4]
		symbolMap[res[4]] = &symbolInfo
	}
}

// Output 输出比对结果
func (nm *Nm) Output() {
	resMap := make(map[string][][2]*SymbolInfo) // 存放比对结果
	for key1, value1 := range nm.lib1Symbol {
		value2, ok := nm.lib2Symbol[key1]
		if !ok {
			continue
		}
		// 两个图中存在相同的symbol则将该symbol的信息存放到结果
		libRes, ok := resMap[value2.LibName]
		if !ok {
			libRes = make([][2]*SymbolInfo, 0, 0)
		}
		libRes = append(libRes, [2]*SymbolInfo{value1, value2})
		resMap[value2.LibName] = libRes // 如果切片的容量不足 append将重新开辟一块内存 所以这里要重新赋值
	}
	// 输出结果
	for key, val := range resMap {
		fmt.Printf("------------------------------------%s vs %s------------------------------------\n", val[0][0].LibName, key)
		fmt.Printf("%-50s%-35s%-5s  |  %-5s%-35s%-50s\n", "Symbol", "Filename", "type", "type", "Filename", "Symbol")
		for _, line := range val {
			fmt.Printf("%-50s%-35s%-5s  |  %-5s%-35s%-50s\n", line[0].SymbolName, line[0].Filename, line[0].Type,
				line[1].Type, line[1].Filename, line[1].SymbolName)
		}
		fmt.Print("\n\n\n")
	}
}

// Clean 清除已经加载的lib
func (nm *Nm) Clean() {
	// 重新make即可 内存回收交给gc
	nm.lib1Symbol = make(map[string]*SymbolInfo)
	nm.lib2Symbol = make(map[string]*SymbolInfo)
}
