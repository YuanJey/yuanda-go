package order

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func NewReadFiles(dir string) []string {
	entries, err := os.ReadDir("./" + dir)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return nil
	}
	var dirs []string
	for _, entry := range entries {
		if !entry.IsDir() {
			// 只处理文件
			filePath := filepath.Join(dir, entry.Name())
			fmt.Println("找到文件: ", filePath)
			dirs = append(dirs, filePath)
		}
	}
	return dirs
}
func ReadFiles() []string {
	yesterday := time.Now().AddDate(0, 0, -1)
	dir := yesterday.Format("2006-01-02")
	//dir := time.Now().Format("2006-01-02")
	entries, err := os.ReadDir("./" + dir)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return nil
	}
	var dirs []string
	for _, entry := range entries {
		if !entry.IsDir() {
			// 只处理文件
			filePath := filepath.Join(dir, entry.Name())
			fmt.Println("找到文件: ", filePath)
			dirs = append(dirs, filePath)
		}
	}
	return dirs
}
func ReadOrdersFile(filePath string) map[string]string {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("打开文件出错：", err)
		panic(err)
	}
	defer file.Close()
	m := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		//fmt.Println(scanner.Text())
		split := strings.Split(scanner.Text(), "\t")
		if len(split) > 1 {
			m[split[0]] = split[1]
		} else {
			fmt.Println("格式错误", scanner.Text())
			panic("格式错误")
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("读取文件出错：", err)
		panic(err)
	}
	return m
}
