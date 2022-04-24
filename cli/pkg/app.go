package pkg

import "github.com/spf13/cobra"

func Start() {
	app := app()
	env := env()
	env.AddCommand(getEnv())
	app.AddCommand(env)
	app.Execute()
}

func app() *cobra.Command {
	return &cobra.Command{
		Use:   "ytool",
		Short: "一些乱七八糟的工具",
	}
}

func env() *cobra.Command {
	return &cobra.Command{
		Use:   "env",
		Short: "env 相关操作",
	}
}

func getEnv() *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "获取环境变量",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			printEnv(args[0])
		},
	}
}
