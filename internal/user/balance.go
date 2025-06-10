package user

import (
	"fmt"
	"github.com/go-rod/rod"
	"regexp"
	"strconv"
)

func GetBalance(page *rod.Page) (float64, error) {
	// 等待元素出现
	page.MustWait(`() => document.querySelector('span.corg') !== null`).MustWaitLoad()

	// 获取文本
	balanceStr := page.MustElement("span.corg").MustText()

	// 提取金额（支持逗号和小数）
	re := regexp.MustCompile(`([0-9\,\.]+)`)
	match := re.FindStringSubmatch(balanceStr)
	if len(match) < 2 {
		return 0, fmt.Errorf("余额格式不正确: %s", balanceStr)
	}

	// 去除逗号后转换为 float64
	cleaned := regexp.MustCompile(`\,`).ReplaceAllString(match[1], "")
	balance, err := strconv.ParseFloat(cleaned, 64)
	if err != nil {
		return 0, fmt.Errorf("转换余额失败: %v", err)
	}
	return balance, nil
}
