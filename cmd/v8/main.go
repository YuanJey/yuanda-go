package main

import (
	"flag"
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"log"
	"strconv"
	"sync"
	"time"
	"yuanda-go/internal/auth"
	"yuanda-go/internal/buy"
	"yuanda-go/internal/order"
	"yuanda-go/internal/user"
	account2 "yuanda-go/pkg/account"
	"yuanda-go/pkg/captcha"
	"yuanda-go/pkg/consts"
)

func main() {
	var account string
	var password string
	var is3w bool
	var num100 int
	var num200 int
	var num500 int
	var num1000 int
	var num2000 int
	flag.StringVar(&account, "name", "", "name")
	flag.StringVar(&password, "password", "", "password")
	flag.BoolVar(&is3w, "is3w", true, "是否足额购买")
	flag.IntVar(&num100, "num100", 0, "100元")
	flag.IntVar(&num200, "num200", 0, "200元")
	flag.IntVar(&num500, "num500", 0, "500元")
	flag.IntVar(&num1000, "num1000", 0, "1000元")
	flag.IntVar(&num2000, "num2000", 0, "2000元")
	flag.Parse()
	log.Println("账号:", account, "密码:", password)
	log.Println("是否足额购买:", is3w)
	log.Println("购买100的数量:", num100)
	log.Println("购买200的数量:", num200)
	log.Println("购买500的数量:", num500)
	log.Println("购买1000元的数量:", num1000)
	log.Println("购买2000元的数量:", num2000)
	// 启动浏览器（headless 模式可通过参数控制）
	l := launcher.New().Headless(false).MustLaunch()
	browser := rod.New().ControlURL(l).MustConnect()

	defer browser.MustClose()
	page := browser.MustPage(consts.Login)
	fmt.Println("已打开登录页")
	for {
		err := auth.SaveCode(page)
		if err != nil {
			log.Printf("验证码截图失败: %v", err)
			continue
		}
		code, err := captcha.IdentifyCode("veriimg.png")
		if err != nil {
			log.Printf("验证码识别失败: %v", err)
			continue
		}
		log.Printf("验证码识别结果: %s", code)

		loginErr := auth.Login(page, account, password, code)
		if loginErr != nil {
			log.Printf("登录失败: %v", loginErr)
			continue
		}
		// 成功登录
		log.Println(account, "登录成功！")
		break
	}
	orderErr := order.DownloadOrder(account)
	if orderErr != nil {
		log.Fatalf("订单下载失败: %v", orderErr)
		return
	}
	for {
		err := page.Reload()
		if err != nil {
			log.Printf("获取余额 页面刷新失败: %v", err)
			return
		}
		balance, err2 := user.GetBalance(page)
		if err2 != nil {
			log.Fatalf("获取余额失败: %v", err2)
			return
		}
		account2.SaveBalanceToFile(account, strconv.FormatFloat(balance, 'f', 2, 64))
		log.Println("账号:", account, "余额:", balance)
		if balance < 30000 && is3w {
			time.Sleep(30 * time.Second)
			log.Println("余额不足30000元，请充值,等待中。。。。")
			continue
		} else {
			log.Println("开始执行金额:", balance)
			break
		}
	}
	var wg sync.WaitGroup
	maxConcurrency := 5
	semaphore := make(chan struct{}, maxConcurrency)

	startBuy := func(num int, url string, count int) {
		if count <= 0 {
			return
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			buyInstance := buy.Buys{
				Num:            num,
				PageUrl:        url,
				NumberOfCycles: count,
			}

			msg, err := buyInstance.Start(browser)
			if err != nil {
				log.Printf("%d元购买失败: %v | 错误: %s", num, err, msg)
			}
		}()
	}

	// 启动各种面额购买任务
	startBuy(100, consts.JD_100, num100)
	startBuy(200, consts.JD_200, num200)
	startBuy(500, consts.JD_500, num500)
	startBuy(1000, consts.JD_1000, num1000)
	startBuy(2000, consts.JD_2000, num2000)

	wg.Wait()
	newBalance, err := user.GetBalance(page)
	if err != nil {
		log.Fatalf("获取余额失败: %v", err)
	}
	log.Println(account, "购买完成, 当前余额:", newBalance)
}
