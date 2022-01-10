package main

import (
	"github.com/osgochina/dmicro/easyservice"
)

func main() {
	easyservice.SetProcessName("test-set-process-title")
	easyservice.Setup(func(service *easyservice.EasyService) {
	})
}
