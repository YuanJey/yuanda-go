package verification

import (
	"log"
	"sync"
	"yuanda-go/pkg/weak_network"
)

func WeakVerification(k, v string) error {
	url := "https://hx.yuanda.biz/Home/Card/writeOffCard"
	req := make(map[string]string)
	req["cardkey"] = k
	req["cardpwd"] = v
	req["cardid"] = "351"
	req["priceid"] = "2846"
	req["typeid"] = "109"
	resp := Resp{}
	err := weak_network.Post(url, req, &resp)
	if err != nil || resp.Status != 1 {
		log.Println("核销失败 卡密信息：", k, v)
		SaveErrOrderToFile(k, v)
		return err
	}
	log.Println("succeed", k, v, resp)
	return nil
}
func WeakVerificationConcurrent(k, v string, times int) {
	var wg sync.WaitGroup
	for i := 0; i < times; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			err := WeakVerification(k, v)
			if err != nil {
				log.Printf("第 %d 次请求失败：%v", i+1, err)
			} else {
				log.Printf("第 %d 次请求成功", i+1)
			}
		}(i)
	}

	wg.Wait()
	log.Println("所有并发请求已完成")
}
