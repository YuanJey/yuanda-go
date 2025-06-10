package order

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"yuanda-go/pkg/config"
)

func DownloadOrder(account string) error {
	// 获取昨天日期
	yesterday := time.Now().AddDate(0, 0, -1)
	dateStr := yesterday.Format("2006-01-02")

	// 创建目录
	dirPath := filepath.Join(".", dateStr)
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return fmt.Errorf("创建目录失败: %v", err)
	}
	// 构建 URL 和 文件路径
	url := fmt.Sprintf("https://sc.yuanda.biz/jingdian/index/export.html?start=%s&end=", dateStr)
	filePath := filepath.Join(dirPath, account+".txt")

	// 创建 HTTP 请求
	client := &http.Client{Timeout: 10 * time.Second}
	req, _ := http.NewRequest("GET", url, nil)

	// 设置 Cookie
	for i := range config.Conf.Cookies {
		req.AddCookie(&http.Cookie{Name: config.Conf.Cookies[i].Name, Value: config.Conf.Cookies[i].Value})
	}
	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != 200 {
		return fmt.Errorf("下载失败，状态码: %d", resp.StatusCode)
	}

	// 创建文件并写入内容
	outFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("创建文件失败: %v", err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return fmt.Errorf("写入文件失败: %v", err)
	}

	fmt.Printf("文件已下载到: %s\n", filePath)
	return nil
}
