package buy

import (
	"fmt"
	"github.com/go-rod/rod"
	"log"
	"time"
)

// Process 模拟用户购买流程（使用 rod）
func Process(page *rod.Page, number int, url string) (string, error) {
	page = page.Timeout(5 * time.Second).MustNavigate(url).MustWaitLoad()
	// 步骤1: 导航到商品页面
	page = page.MustNavigate(url).MustWaitLoad()
	log.Printf("已打开商品页: %s", url)

	// 步骤2: 等待并点击购买按钮
	//log.Println("等待购买按钮出现...")
	page.MustElement("div.cart-buy > a.buy-btn").MustWaitVisible().MustScrollIntoView().MustClick()
	log.Println("已点击购买按钮，等待页面跳转...")

	// 步骤4: 等待“找人代付”按钮并点击
	//log.Println("等待【找人代付】按钮出现...")
	alipayBtn := page.MustElement("#alipay").MustWaitVisible().MustScrollIntoView().MustClick()
	log.Printf("已点击【%s】按钮", alipayBtn.MustText())

	// 步骤5: 等待结算按钮并点击
	//log.Println("等待【结算】按钮出现...")
	jiesuanBtn := page.MustElement("#jiesuan").MustWaitVisible().MustScrollIntoView().MustClick()
	log.Printf("已点击【%s】按钮", jiesuanBtn.MustText())

	// 步骤6: 获取成功信息https://sc.yuanda.biz/jingdian/Getmail/index/mpid
	//log.Println("等待成功信息出现...")
	messageElm := page.MustElement("#zhengwen").MustWaitVisible()
	messageText, err := messageElm.Text()
	if err != nil {
		return "", fmt.Errorf("获取成功信息失败: %v", err)
	}
	log.Printf("成功信息: %v %s", number, messageText)
	return messageText, nil
}
