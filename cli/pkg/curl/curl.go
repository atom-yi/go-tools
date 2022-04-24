package curl

import (
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
)

func Curl() *cobra.Command {
	curl := &cobra.Command{
		Use:   "curl",
		Short: "类似 curl 工具",
	}
	curl.AddCommand(getUrl())
	return curl
}

func getUrl() *cobra.Command {
	return &cobra.Command{
		Use:     "get",
		Short:   "发送 get 请求",
		Args:    cobra.ExactArgs(1),
		Example: "ytool curl get http://www.baidu.com",
		Run: func(cmd *cobra.Command, args []string) {
			content, err := curlGet(args[0])
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(string(content))
		},
	}
}

func curlGet(url string) ([]byte, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
