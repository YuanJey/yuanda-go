package verification

import (
	"fmt"
	"log"
	"os"
	"time"
	"yuanda-go/pkg/http_client"
)

type Resp struct {
	Status int    `json:"status"`
	Info   string `json:"info"`
	Data   string `json:"data"`
}

func Verification(k, v string) error {
	url := "https://hx.yuanda.biz/Home/Card/writeOffCard"
	req := make(map[string]string)
	req["cardkey"] = k
	req["cardpwd"] = v
	req["cardid"] = "351"
	req["priceid"] = "2846"
	req["typeid"] = "109"
	resp := Resp{}
	err := http_client.Post(url, req, &resp)
	if err != nil || resp.Status != 1 {
		log.Println("核销失败 卡密信息：", k, v)
		SaveErrOrderToFile(k, v)
		return err
	}
	log.Println("succeed", k, v, resp)
	return nil
}
func init() {
	currentDate := time.Now().Format("2006-01-02")
	filename := currentDate + "_err.txt"
	_, err := os.ReadFile(filename)
	if err != nil {
		file, err1 := os.Create(filename)
		if err1 != nil {
			log.Println("创建文件失败：", err1)
			return
		}
		defer file.Close()
	}
}
func SaveErrOrderToFile(k, v string) {
	currentDate := time.Now().Format("2006-01-02")
	filename := currentDate + "_err.txt"
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("打开文件失败：", err)
		return
	}
	defer file.Close()
	_, err = file.WriteString(fmt.Sprintf("%s%s%s\n", k, "\t", v))
}
