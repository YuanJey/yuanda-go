package main

import (
	"flag"
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"log"
	"strconv"
	"time"
	"yuanda-go/internal/auth"
	"yuanda-go/internal/buy"
	"yuanda-go/internal/order"
	"yuanda-go/internal/user"
	account2 "yuanda-go/pkg/account"
	"yuanda-go/pkg/captcha"
	"yuanda-go/pkg/check"
	"yuanda-go/pkg/consts"
)

func main() {
	if check.Check() == 0 {
		log.Fatalf("未通过检测，请联系管理员")
		return
	}
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
	rod.New()
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
		log.Println("账号:", account, "余额:", balance)
		if balance < 30000 && is3w {
			time.Sleep(30 * time.Second)
			log.Fatalf("余额不足30000元，请充值,等待中。。。。")
			continue
		} else {
			log.Println("开始执行金额:", balance)
			break
		}
	}
	m := make(map[string]int)
	m["100"] = 0
	m["200"] = 0
	m["500"] = 0
	m["1000"] = 0
	m["2000"] = 0
	for i := range num100 {
		process, err := buy.Process(page, 100, consts.JD_100)
		if err != nil {
			m["100"] = m["100"] + 1
			log.Fatalf("第 %v 次100 购买失败: %v", i, process)
		}
	}
	for i := range num200 {
		process, err := buy.Process(page, 200, consts.JD_200)
		if err != nil {
			m["200"] = m["200"] + 1
			log.Fatalf("第 %v 次200 购买失败: %v", i, process)
		}
	}
	for i := range num500 {
		process, err := buy.Process(page, 500, consts.JD_500)
		if err != nil {
			m["500"] = m["500"] + 1
			log.Fatalf("第 %v 次500 购买失败: %v", i, process)
		}
	}
	for i := range num1000 {
		process, err := buy.Process(page, 1000, consts.JD_1000)
		if err != nil {
			m["1000"] = m["1000"] + 1
			log.Fatalf("第 %v 次1000 购买失败: %v", i, process)
		}
	}
	for i := range num2000 {
		process, err := buy.Process(page, 2000, consts.JD_2000)
		if err != nil {
			m["2000"] = m["2000"] + 1
			log.Fatalf("第 %v 次2000 购买失败: %v", i, process)
		}
	}
	balance, err3 := user.GetBalance(page)
	if err3 != nil {
		log.Fatalf("获取余额失败: %v", err3)
		return
	}
	log.Println(account, "购买完成,当前余额:", balance)
	if sum, ok := m["100"]; ok {
		log.Println(account, "100购买失败次数:", sum)
	}
	if sum, ok := m["200"]; ok {
		log.Println(account, "200购买失败次数:", sum)
	}
	if sum, ok := m["500"]; ok {
		log.Println(account, "500购买失败次数:", sum)
	}
	if sum, ok := m["1000"]; ok {
		log.Println(account, "1000购买失败次数:", sum)
	}
	if sum, ok := m["2000"]; ok {
		log.Println(account, "2000购买失败次数:", sum)
	}
}
