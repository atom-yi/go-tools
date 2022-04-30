package curl

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// 请求类型
type ReqType string

const (
	Unknown = "Unknown"
	Get     = "GET"
	Post    = "POST"
)

type Header struct {
	key   string
	value string
}

// 请求信息结构体
type Request struct {
	reqType ReqType
	url     string
	headers []Header
	body    string
}

func Curl() *cobra.Command {
	curl := &cobra.Command{
		Use:              "curl",
		Short:            "类似 curl 工具",
		TraverseChildren: true,
	}
	curl.AddCommand(curlGet())
	curl.AddCommand(curlPost())
	return curl
}

func curlPost() *cobra.Command {
	var curlFilePath string
	var url string
	var json string
	postCmd := &cobra.Command{
		Use:   "post",
		Short: "发送 post 请求",
		Run: func(cmd *cobra.Command, args []string) {
			request, err := loadCurlFile(curlFilePath)
			if err != nil {
				fmt.Println(err)
				return
			}
			if url != "" {
				request.url = url
			}
			if json != "" {
				request.body = json
			}

			resp, err := doPost(&request)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(resp)
		},
	}
	postCmd.Flags().StringVarP(&url, "url", "u", "", "url")
	postCmd.Flags().StringVarP(&json, "json", "j", "", "json")
	postCmd.Flags().StringVarP(&curlFilePath, "load", "l", "", "加载 curl 文件")
	return postCmd
}

func curlGet() *cobra.Command {
	var curlFilePath string
	var url string
	getCmd := &cobra.Command{
		Use:   "get",
		Short: "发送 get 请求",
		Run: func(cmd *cobra.Command, args []string) {
			// 加载 curl 文件，构造请求
			request, err := loadCurlFile(curlFilePath)
			if err != nil {
				fmt.Println(err)
				return
			}
			// 装配 url
			if url != "" {
				request.url = url
			}
			resp, err := doGet(request)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("url: \n", request.url)
			fmt.Println("headers: \n", request.headers)
			fmt.Println("resp: \n", resp)
		},
	}
	getCmd.Flags().StringVarP(&url, "url", "u", "", "url")
	getCmd.Flags().StringVarP(&curlFilePath, "load", "l", "", "加载 curl 文件")
	return getCmd
}

func doGet(req Request) (string, error) {
	req.reqType = Get
	return request(req)
}

func doPost(req *Request) (string, error) {
	req.reqType = Post
	return request(*req)
}

func request(req Request) (string, error) {
	// todo: 检查请求是否合法
	if err := checkReq(req); err != nil {
		return "", err
	}

	// 调用请求
	resp, err := invokeReq(req)
	if err != nil {
		return "", err
	}
	return resp, nil
}

func checkReq(req Request) error {
	if req.reqType == Unknown {
		return fmt.Errorf("未知请求类型")
	}
	if req.url == "" {
		return fmt.Errorf("url 不能为空")
	}
	return nil
}

func loadCurlFile(curlFilePath string) (Request, error) {
	request := Request{
		reqType: Unknown,
		url:     "",
		headers: nil,
	}
	if curlFilePath == "" {
		return request, nil
	}

	// 解析 curl 文件
	headers := []Header{}
	err := traverseLineInFile(curlFilePath, func(line string) {
		line = strings.Trim(line, " ")
		if strings.HasPrefix(line, "curl") {
			// 解析连接地址，原始链接信息如下格式：
			// curl 'http://xxxx' \
			url := strings.TrimSuffix(line[6:], "' \\")
			request.url = url
		} else if header, match := parseHeader(line, &request); match {
			// 解析 header
			headers = append(headers, header)
		} else if strings.HasPrefix(line, "--data-raw") {
			// 解析请求内容
			request.body = strings.Trim(line[12:], "' \\")
		}
	})
	if err != nil {
		return request, err
	}
	if len(headers) > 0 {
		request.headers = headers
	}
	return request, nil
}

func parseHeader(line string, request *Request) (Header, bool) {
	header := Header{}
	if !strings.HasPrefix(line, "-H '") {
		return header, false
	}
	lastIdx := strings.LastIndex(line, "'")
	if lastIdx == -1 {
		return header, false
	}
	headerStr := line[4:lastIdx]
	words := strings.Split(headerStr, ": ")
	if len(words) != 2 {
		return header, false
	}
	return Header{key: words[0], value: words[1]}, true
}

func traverseLineInFile(filePath string, callback func(string)) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		line = line[:len(line)-1]
		callback(line)
	}
	return nil
}

func invokeReq(request Request) (string, error) {
	client := &http.Client{}

	// 创建请求
	req, err := http.NewRequest(string(request.reqType), request.url, nil)
	if err != nil {
		return "", err
	}

	// 设置请求头
	if request.headers != nil {
		for _, header := range request.headers {
			req.Header.Add(header.key, header.value)
		}
	}

	// 设置请求内容
	if request.body != "" {
		req.Body = ioutil.NopCloser(strings.NewReader(request.body))
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 读取响应内容
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(respBody), nil
}
