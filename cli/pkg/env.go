package pkg

import (
	"fmt"
	"os"

	"github.com/atom-yi/cli/pkg/yutil"
)

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
