package logic

import (
	"fmt"
	"github.com/avast/retry-go"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"log"
	"sync"
	verification "yuanda-go/internal/hx"
	"yuanda-go/internal/order"
	"yuanda-go/pkg/config"
	"yuanda-go/pkg/consts"
)

func StartHX() {
	l := launcher.New().Headless(false).MustLaunch()
	browser := rod.New().ControlURL(l).MustConnect()

	defer browser.MustClose()
	page := browser.MustPage(consts.HX)
	log.Println("已打开核销页")
	//等待用户回车确定继续执行
	log.Println("订单全部下载完成后请输入1继续执行，或输入其他内容取消...")
	// 等待用户输入
	var input string
	_, err := fmt.Scanln(&input)
	if err != nil || input != "1" {
		log.Println("取消执行")
		return
	}
	log.Println("开始执行核销流程...")
	cookies, err := page.Cookies([]string{consts.HX})
	if err != nil {
		log.Fatalf("获取 Cookie 失败: %v", err)
		return
	}

	// 保存到全局配置
	config.Conf.Cookies = cookies
	// 并发控制
	maxConcurrency := 5 // 控制最大并发数量
	semaphore := make(chan struct{}, maxConcurrency)

	// 获取所有订单文件
	files := order.ReadFiles()

	var wg sync.WaitGroup
start:
	// 并发处理每个文件
	for _, file := range files {
		wg.Add(1)
		go func(file string) {
			defer wg.Done()

			ordersMap := order.ReadOrdersFile(file)
			var innerWg sync.WaitGroup

			for k, v := range ordersMap {
				semaphore <- struct{}{} // 获取信号量
				innerWg.Add(1)
				go func(k, v string) {
					defer func() {
						innerWg.Done()
						<-semaphore // 释放信号量
					}()
					retry.Do(func() error {
						err := verification.Verification(k, v)
						if err != nil {
							log.Printf("订单 %s 核销失败: %v\n", k, err)
							return err
						}
						return nil
					},
						retry.Attempts(10),
					)
				}(k, v)
			}
			innerWg.Wait()
		}(file)
	}

	wg.Wait()
	log.Println("所有订单核销完成")
	log.Println("请输入1再次执行，或输入其他内容结束...")
	// 等待用户输入
	var end string
	_, err1 := fmt.Scanln(&end)
	if err1 != nil || end != "1" {
		log.Println("执行结束。。。")
		page.Close()
		return
	} else {
		goto start
	}
}
func NewStartHX() {
	l := launcher.New().Headless(false).MustLaunch()
	browser := rod.New().ControlURL(l).MustConnect()

	defer browser.MustClose()
	page := browser.MustPage(consts.HX)
	log.Println("已打开核销页")
	//等待用户回车确定继续执行
	log.Println("订单全部下载完成后请输入1继续执行，或输入其他内容取消...")
	// 等待用户输入
	var input string
	_, err := fmt.Scanln(&input)
	if err != nil || input != "1" {
		log.Println("取消执行")
		return
	}
	log.Println("开始执行核销流程...")
	cookies, err := page.Cookies([]string{consts.HX})
	if err != nil {
		log.Fatalf("获取 Cookie 失败: %v", err)
		return
	}

	// 保存到全局配置
	config.Conf.Cookies = cookies
	// 并发控制
	maxConcurrency := 5 // 控制最大并发数量
	semaphore := make(chan struct{}, maxConcurrency)

	// 获取所有订单文件
	files := order.ReadFiles()

	var wg sync.WaitGroup
start:
	// 并发处理每个文件
	for _, file := range files {
		wg.Add(1)
		go func(file string) {
			defer wg.Done()

			ordersMap := order.ReadOrdersFile(file)
			var innerWg sync.WaitGroup

			for k, v := range ordersMap {
				semaphore <- struct{}{} // 获取信号量
				innerWg.Add(1)
				go func(k, v string) {
					defer func() {
						innerWg.Done()
						<-semaphore // 释放信号量
					}()
					//retry.Do(func() error {
					//	err := verification.WeakVerification(k, v)
					//	if err != nil {
					//		log.Printf("订单 %s 核销失败: %v\n", k, err)
					//		return err
					//	}
					//	return nil
					//},
					//	retry.Attempts(10),
					//)
					verification.WeakVerificationConcurrent(k, v, 5)
				}(k, v)
			}
			innerWg.Wait()
		}(file)
	}

	wg.Wait()
	log.Println("所有订单核销完成")
	log.Println("请输入1再次执行，或输入其他内容结束...")
	// 等待用户输入
	var end string
	_, err1 := fmt.Scanln(&end)
	if err1 != nil || end != "1" {
		log.Println("执行结束。。。")
		page.Close()
		return
	} else {
		goto start
	}
}
