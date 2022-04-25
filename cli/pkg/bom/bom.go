package bom

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func Bom() *cobra.Command {
	bom := &cobra.Command{
		Use:   "bom",
		Short: "文件BOM头操作",
	}
	bom.AddCommand(addBom(), removeBom(), existBom())
	return bom
}

func addBom() *cobra.Command {
	return &cobra.Command{
		Use:     "add",
		Short:   "给文件添加BOM头",
		Args:    cobra.ExactArgs(1),
		Example: "ytool bom add ./test.txt",
		Run: func(cmd *cobra.Command, args []string) {
			err := addBomToFile(args[0])
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("添加BOM头成功")
		},
	}
}

func removeBom() *cobra.Command {
	return &cobra.Command{
		Use:     "remove",
		Short:   "移除文件BOM头",
		Args:    cobra.ExactArgs(1),
		Example: "ytool bom remove ./test.txt",
		Run: func(cmd *cobra.Command, args []string) {
			err := removeBomFromFile(args[0])
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("移除BOM头成功")
		},
	}
}

func existBom() *cobra.Command {
	return &cobra.Command{
		Use:     "exist",
		Short:   "判断文件是否存在BOM头",
		Args:    cobra.ExactArgs(1),
		Example: "ytool bom exist ./test.txt",
		Run: func(cmd *cobra.Command, args []string) {
			exist, err := existBomInFile(args[0])
			if err != nil {
				fmt.Println(err)
				return
			}
			if exist {
				fmt.Println("存在BOM头")
			} else {
				fmt.Println("不存在BOM头")
			}
		},
	}
}

func addBomToFile(filePath string) error {
	existBom, err := existBomInFile(filePath)
	if err != nil {
		return err
	}
	if existBom {
		return fmt.Errorf("文件已存在BOM头")
	}

	return rwFileInHead(filePath, func(reader *bufio.Reader, writer *bufio.Writer) {
		writer.Write([]byte{0xEF, 0xBB, 0xBF})
	})
}

func removeBomFromFile(filePath string) error {
	existBom, err := existBomInFile(filePath)
	if err != nil {
		return err
	}
	if !existBom {
		return fmt.Errorf("文件不存在BOM头")
	}

	return rwFileInHead(filePath, func(reader *bufio.Reader, writer *bufio.Writer) {
		reader.Discard(3)
	})
}

func rwFileInHead(filePath string, preHandle func(*bufio.Reader, *bufio.Writer)) error {
	return selfCopy(filePath, preHandle, nil)
}

func selfCopy(filePath string, preHandle func(*bufio.Reader, *bufio.Writer),
	postHandle func(*bufio.Reader, *bufio.Writer)) error {
	file, tmpFile, err := getBindTmpFile(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	defer tmpFile.Close()
	reader := bufio.NewReader(file)
	writer := bufio.NewWriter(tmpFile)
	if preHandle != nil {
		preHandle(reader, writer)
	}
	io.Copy(writer, reader)
	if postHandle != nil {
		postHandle(reader, writer)
	}
	writer.Flush()
	defer replaceFile(tmpFile.Name(), filePath)
	return nil
}

func getBindTmpFile(filePath string) (*os.File, *os.File, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, err
	}
	tmpFile, err := createTempFile(filepath.Dir(filePath))
	if err != nil {
		defer file.Close()
		return nil, nil, err
	}
	return file, tmpFile, nil
}

func replaceFile(source string, target string) {
	targetBak := target + ".ybak"
	os.Rename(target, targetBak)
	os.Rename(source, target)
	os.Remove(targetBak)
}

func createTempFile(dir string) (*os.File, error) {
	tmpFile, err := ioutil.TempFile(dir, "*.ytmp")
	if err != nil {
		return nil, err
	}
	return tmpFile, nil
}

func existBomInFile(filePath string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	if err != nil {
		return false, err
	}

	buf := make([]byte, 3)
	_, err = reader.Read(buf)
	if err == io.EOF {
		// 文件头不足3字节
		return false, nil
	} else if err != nil {
		// 文件读取错误
		return false, err
	} else if buf[0] != 0xEF || buf[1] != 0xBB || buf[2] != 0xBF {
		// 文件不存在BOM头
		return false, nil
	} else {
		// 文件存在BOM头
		return true, nil
	}
}
