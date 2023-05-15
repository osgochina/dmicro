# Cobra自定义扩展命令

通过Cobra自定义扩展命令，可以在Cobra的基础上，实现更多的功能。
```go
dserver.Cobra(func(rootCmd *cobra.Command) {
	// 添加自定义命令
	rootCmd.AddCommand(&cobra.Command{
        Use:   "test",
        Short: "test command",
        Run: func(cmd *cobra.Command, args []string) {
        cmd.Println("test")
        cmd.Println(args)
        },
    }))
})
```