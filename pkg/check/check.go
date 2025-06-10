package check

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func getKey() (string, error) {
	// 打开文件
	file, err := os.Open("key.txt")
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			return line, nil
		}
	}

	// 如果没有找到非空行
	return "", nil
}
func Check() int {
	src, err2 := getKey()
	if err2 != nil {
		fmt.Println("获取key失败:", err2)
		return 0
	}
	url := "https://test-1312265679.cos.ap-chengdu.myqcloud.com/config_v7.json"

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("请求错误:", err)
		return 0
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("网络请求失败或未找到匹配的机器信息")
		return 0
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应失败:", err)
		return 0
	}

	var jsonData map[string][]string
	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		fmt.Println("解析JSON失败:", err)
		return 0
	}

	checksList, ok := jsonData["checks"]
	if !ok {
		fmt.Println("未找到checks字段")
		return 0
	}

	for _, item := range checksList {
		if item == src {
			return 1
		}
	}
	fmt.Println("网络请求失败或未找到匹配的机器信息")
	return 0
}
