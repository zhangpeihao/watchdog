// Copyright © 2016 Zhang Peihao <zhangpeihao@gmail.com>

package main

import (
	"fmt"
	"os"

	"github.com/zhangpeihao/watchdog/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
