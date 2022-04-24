package env

import (
	"fmt"
	"os"

	"github.com/atom-yi/cli/pkg/yutil"
	"github.com/spf13/cobra"
)

func Env() *cobra.Command {
	env := &cobra.Command{
		Use:   "env",
		Short: "env 相关操作",
	}
	env.AddCommand(getEnv())
	return env
}

func printEnv(envPropName string) {
	if yutil.IsBlankStr(envPropName) {
		fmt.Println("empty env prop name")
		return
	}

	envPropValue := os.Getenv(envPropName)
	if yutil.IsBlankStr(envPropValue) {
		fmt.Printf("env prop %s is empty\n", envPropName)
		return
	}

	fmt.Println(envPropValue)
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
