package auth

import (
	"fmt"
	"github.com/go-rod/rod"
	"log"
	"time"
	"yuanda-go/pkg/config"
	"yuanda-go/pkg/consts"
)

// Login 使用 rod 实现登录流程并保存 cookies
// Login 使用 rod 实现登录流程并判断是否跳转到用户中心
func Login(page *rod.Page, account, password, code string) error {
	// 填写账号密码和验证码
	page.MustElement("#account").MustInput(account)
	page.MustElement("#password").MustInput(password)
	page.MustElement("#veri").MustInput(code)

	fmt.Println("已填写登录信息")

	// 点击登录按钮
	page.MustElement("#loginbtn").MustClick()
	// 检查当前页面 URL 是否符合预期
	time.Sleep(3 * time.Second)
	currentURL := page.MustInfo().URL
	if currentURL != consts.ProfileRedirectURL {
		return fmt.Errorf("登录失败：未跳转到用户中心，当前 URL: %s", currentURL)
	}

	// 获取 Cookies
	cookies, err := page.Cookies([]string{consts.Login})
	if err != nil {
		log.Fatalf("获取 Cookie 失败: %v", err)
		return err
	}

	// 保存到全局配置
	config.Conf.Cookies = cookies
	fmt.Println("登录成功，Cookies 已保存")
	return nil
}

//func Login(ctx context.Context, account, password, code string) error {
//	var cookies []*network.Cookie
//	err := chromedp.Run(ctx,
//		chromedp.WaitVisible(`#account`, chromedp.ByID),
//		chromedp.SendKeys(`#account`, account, chromedp.ByID),
//		chromedp.WaitVisible(`#password`, chromedp.ByID),
//		chromedp.SendKeys(`#password`, password, chromedp.ByID),
//		chromedp.WaitVisible(`#veri`, chromedp.ByID),
//		chromedp.SendKeys(`#veri`, code, chromedp.ByID),
//		chromedp.WaitVisible(`#loginbtn`, chromedp.ByID),
//		chromedp.Click(`#loginbtn`, chromedp.ByID),
//	)
//	if err != nil {
//		log.Fatalf("登录失败: %v", err)
//		return err
//	}
//	time.Sleep(5 * time.Second)
//	err = chromedp.Run(ctx,
//		chromedp.ActionFunc(func(ctx context.Context) error {
//			cookiesRaw, err := network.GetCookies().Do(ctx)
//			if err != nil {
//				return err
//			}
//			cookies = cookiesRaw
//			return nil
//		}),
//	)
//	if err != nil {
//		log.Fatalf("获取cookie失败: %v", err)
//		return err
//	}
//	config.Conf.Cookies = cookies
//	log.Println("登录成功，已保存Cookies")
//	return nil
//}
