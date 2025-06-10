package auth

import (
	"fmt"
	"github.com/go-rod/rod"
	"log"
	"os"
)

func SaveCode(page *rod.Page) error {
	// 等待验证码图片出现并截图
	img := page.MustElement("#veriimg")
	buf, err := img.Resource()
	if err != nil {
		log.Fatalf("获取验证码图片资源失败: %v", err)
		return err
	}

	// 保存图片
	err = os.WriteFile("veriimg.png", buf, 0644)
	if err != nil {
		log.Fatalf("保存验证码图片失败: %v", err)
		return err
	}

	fmt.Println("验证码图片已保存为 veriimg.png")
	return nil
}

//func SaveCode(ctx context.Context) error {
//	err := chromedp.Run(ctx,
//		chromedp.Navigate(consts.Login),
//	)
//	if err != nil {
//		log.Fatalf("导航到登录页面失败: %v", err)
//		return err
//	}
//	var buf []byte
//	err = chromedp.Run(ctx,
//		chromedp.WaitVisible(`#veriimg`, chromedp.ByID),
//		chromedp.Screenshot(`#veriimg`, &buf, chromedp.NodeVisible, chromedp.ByID),
//	)
//	if err != nil {
//		log.Fatalf("获取验证码失败: %v", err)
//		return err
//	}
//	err = os.WriteFile("veriimg.png", buf, 0644)
//	if err != nil {
//		log.Fatalf("保存验证码图片失败: %v", err)
//		return err
//	}
//	return nil
//}
