package http_client

import (
	"encoding/json"
	"errors"
	"github.com/YuanJey/goutils2/pkg/utils"
	"io"
	"log"
	"net/http"
	url2 "net/url"
	"strings"
	"time"
	"yuanda-go/pkg/config"
)

func Get(operationID string, url string, req interface{}, resp interface{}) error {
	body := strings.NewReader("")
	if req != nil {
		jsonStr, err := json.Marshal(req)
		if err != nil {
			return err
		}
		body = strings.NewReader(string(jsonStr))
	}
	request, err := http.NewRequest("GET", url, body)
	if err != nil {
		return err
	}
	client := http.Client{Timeout: 5 * time.Second}
	httpResponse, err := client.Do(request)
	if err != nil {
		return err
	}
	result, err := io.ReadAll(httpResponse.Body)
	if httpResponse.StatusCode != 200 {
		return utils.Wrap(errors.New(httpResponse.Status), "status code failed "+url+string(result))
	}
	log.Printf(operationID, "api request success url is "+url, string(result))
	err = utils.JsonStringToStruct(string(result), &resp)
	if err != nil {
		return err
	}
	return nil
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
	//request.Header.Add("Cookie", "PHPSESSID="+ssid.Get()+";think_language=zh-CN")
	client := http.Client{Timeout: 5 * time.Second}
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
