package main

import (
	"github.com/desertbit/grumble"
	"github.com/osgochina/dmicro/supervisor/ctrl"
)

func main() {
	grumble.Main(ctrl.App)
}
