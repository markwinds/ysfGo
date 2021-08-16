package main

import (
	"flag"
	"strings"
	"ysf/cmd"
)

var (
	lib1 string // 第一个lib
	lib2 string // 其他几个lib
)

// 检查一个lib和另外几个lib之间是否有依赖关系
func main() {
	flag.StringVar(&lib1, "lib1", "libcrypto.a", "first lib")
	flag.StringVar(&lib2, "lib2", "libssl.a", "second libs, split by comma, eg: libssl.a,libcJSON.a")
	flag.Parse()

	lib2s := strings.Split(lib2, ",")

	nm := cmd.NewNm()
	nm.InitLib1(lib1)
	nm.InitLib2(lib2s)
	nm.Output()
}
