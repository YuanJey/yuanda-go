package weak_network

import (
	"errors"
	"github.com/YuanJey/goutils2/pkg/utils"
	"io"
	"log"
	"math/rand"
	"net/http"
	url2 "net/url"
	"strings"
	"time"
	"yuanda-go/pkg/config"
)

type unstableTransport struct{}

func (t *unstableTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// 模拟延迟
	time.Sleep(time.Duration(rand.Intn(2000)) * time.Millisecond)
	// 模拟丢包或请求失败
	if rand.Float32() < 0.7 { // 30%的请求会失败
		return nil, errors.New("模拟网络不稳定")
	}
	return http.DefaultTransport.RoundTrip(req)
}
func Post(url string, req map[string]string, resp interface{}) error {
	data := url2.Values{}
	for k, v := range req {
		data.Add(k, v)
	}
	dataString := data.Encode()
	body := strings.NewReader(dataString)
	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}
	for i := range config.Conf.Cookies {
		request.AddCookie(&http.Cookie{Name: config.Conf.Cookies[i].Name, Value: config.Conf.Cookies[i].Value})
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{
		Transport: &unstableTransport{},
		Timeout:   1000 * time.Millisecond,
	}
	httpResponse, err := client.Do(request)
	if err != nil {
		return err
	}
	result, err := io.ReadAll(httpResponse.Body)
	if httpResponse.StatusCode != 200 {
		log.Printf("api request err url is "+url, httpResponse.Status, string(result))
		return utils.Wrap(errors.New(httpResponse.Status), "status code failed "+url+string(result))
	}
	err = utils.JsonStringToStruct(string(result), &resp)
	if err != nil {
		return err
	}
	return nil
}
