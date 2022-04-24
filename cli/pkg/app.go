package pkg

import (
	"github.com/atom-yi/cli/pkg/curl"
	"github.com/atom-yi/cli/pkg/env"
	"github.com/spf13/cobra"
)

func Start() {
	app := app()
	app.AddCommand(env.Env())
	app.AddCommand(curl.Curl())
	app.Execute()
}

func app() *cobra.Command {
	return &cobra.Command{
		Use:   "ytool",
		Short: "一些乱七八糟的工具",
	}
}
