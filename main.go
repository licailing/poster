package main

import (
	"github.com/daleyshek/poster/common"
)

func main() {
	common.Init()
	common.ServeHTTP()
	select {}
}
