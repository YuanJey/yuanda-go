package buy

import (
	"github.com/go-rod/rod"
	"log"
	"time"
)

type Buys struct {
	Num            int
	NumberOfCycles int
	PageUrl        string
}

func (b *Buys) Start(browser *rod.Browser) (string, error) {
	page := browser.MustPage().Timeout(15 * time.Second)
	defer page.MustClose()
	for i := 0; i < b.NumberOfCycles; i++ {

		page.MustNavigate(b.PageUrl).MustWaitLoad()

		page.MustElement("div.cart-buy > a.buy-btn").MustWaitVisible().MustScrollIntoView().MustClick()
		log.Println("已点击购买按钮，等待页面跳转...")

		page.MustElement("#alipay").MustWaitVisible().MustScrollIntoView().MustClick()
		log.Println("已点击【找人代付】按钮")

		page.MustElement("#jiesuan").MustWaitVisible().MustScrollIntoView().MustClick()
		log.Println("已点击【结算】按钮")
		messageText := page.MustElement("#zhengwen").MustWaitVisible().MustText()
		log.Printf("%v 元 第 %v 次购买成功", b.Num, i+1)
		log.Printf("成功信息: %s", messageText)

		page.MustClose()
	}
	return "success", nil
}
