package main

import (
	"fmt"
	"github.com/bronze-ding/bindview-cli/src"
	"os"
)

func main() {
	args := os.Args[1:] // 去除第一个参数，即程序本身的路径
	if len(args) == 2 && args[0] == "create" {
		// 列表
		src.List(func(url string) {
			src.Download(url, args[1])
		}, args[1])
	} else {
		fmt.Println("参数不正确")
	}
}
