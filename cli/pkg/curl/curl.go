package curl

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var headerMatcher *regexp.Regexp

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

func init() {
	// todo: 修改正则表达式
	headerMatcher, _ = regexp.Compile("(?<=').*(?=')")
}

func Curl() *cobra.Command {
	curl := &cobra.Command{
		Use:              "curl",
		Short:            "类似 curl 工具",
		TraverseChildren: true,
	}
	curl.AddCommand(getUrl())
	// todo: curl.AddCommand(postUrl())
	return curl
}

func getUrl() *cobra.Command {
	var curlFilePath string
	var url string
	var reqBody string
	getCmd := &cobra.Command{
		Use:     "get",
		Short:   "发送 get 请求",
		Example: "ytool curl get -u http://www.baidu.com",
		Run: func(cmd *cobra.Command, args []string) {
			resp, err := request(curlFilePath, url, reqBody)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("url: \n", url)
			fmt.Println("resp: \n", len(resp))
		},
	}
	getCmd.Flags().StringVarP(&url, "url", "u", "", "url")
	getCmd.Flags().StringVarP(&url, "body", "b", "", "请求内容")
	getCmd.Flags().StringVarP(&curlFilePath, "load", "l", "", "加载 curl 文件")
	return getCmd
}

func request(curlFilePath string, url string, reqBody string) (string, error) {
	// 加载 curl 文件
	request, err := loadCurlFile(curlFilePath)
	if err != nil {
		return "", err
	}

	// 设置 url
	if url != "" {
		request.url = url
	}

	// 设置请求内容
	if reqBody != "" {
		request.body = reqBody
	}

	// todo: 检查请求是否合法

	// 调用请求
	resp, err := invokeReq(request)
	if err != nil {
		return "", err
	}
	return resp, nil
}

func loadCurlFile(curlFilePath string) (*Request, error) {
	request := Request{
		reqType: Unknown,
		url:     "",
		headers: nil,
	}
	if curlFilePath == "" {
		return &request, nil
	}
	headers := []string{}
	err := traverseLineInFile(curlFilePath, func(line string) {
		line = strings.Trim(line, " ")
		parseReqType(line, &request)
		if parseHeader(line, &request) != nil {
			headers = append(headers, line)
		}
		fmt.Println(line)
	})
	if err != nil {
		return nil, err
	}
	return &request, nil
}

func parseHeader(line string, request *Request) *Header {
	if strings.HasPrefix(line, "-H") {
		return nil
	}
	headerStr := headerMatcher.FindString(line)
	if headerStr == "" {
		return nil
	}
	words := strings.Split(headerStr, ": ")
	if len(words) != 2 {
		return nil
	}
	return &Header{key: words[0], value: words[1]}
}

func parseReqType(line string, request *Request) {
	if strings.HasPrefix(line, "-X") {
		return
	} else if strings.Contains(line, "GET") {
		request.reqType = Get
	} else if strings.Contains(line, "POST") {
		request.reqType = Post
	}
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

func invokeReq(request *Request) (string, error) {
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

	// todo: 设置 cookie

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
