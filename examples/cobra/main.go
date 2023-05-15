package main

import (
	"fmt"
	"github.com/osgochina/dmicro/dserver"
	"github.com/spf13/cobra"
)

func main() {
	dserver.CloseCtl()
	dserver.SetName("DMicro")
	dserver.Cobra(func(rootCmd *cobra.Command) {
		cmds := rootCmd.Commands()
		for _, cmd := range cmds {
			fmt.Println(cmd.Name())
			if cmd.Name() == "start" {
				cmd.Flags().StringP("mycmd", "m", "test flag", "test flag")
			}
		}
		// 添加自定义命令
		rootCmd.AddCommand(&cobra.Command{
			Use:   "test",
			Short: "test command",
			Run: func(cmd *cobra.Command, args []string) {
				cmd.Println("test")
				cmd.Println(args)
			},
		})
	})
	dserver.Setup()
}
