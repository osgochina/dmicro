package main

import (
	"fmt"
	"github.com/osgochina/dmicro/dserver"
)

func main() {
	dserver.Authors = "osgochina@gmail.com"
	dserver.SetName("DMicro_simple")
	dserver.Setup(func(svr *dserver.DServer) {
		fmt.Println("start success!")
	})
}
