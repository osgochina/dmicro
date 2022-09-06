package main

import (
	"github.com/osgochina/dmicro/dserver"
	"github.com/spf13/cobra"
)

func main() {
	dserver.CloseCtl()
	dserver.SetName("DMicro")
	dserver.Cobra(func(rootCmd *cobra.Command) {
		rootCmd.AddCommand(&cobra.Command{
			Use:   "test",
			Short: "test command",
			Run: func(cmd *cobra.Command, args []string) {
				cmd.Println("test")
				cmd.Println(args)
			},
		})
	})
	dserver.Setup(func(svr *dserver.DServer) {

	})
}
