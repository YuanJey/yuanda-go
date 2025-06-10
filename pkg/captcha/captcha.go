package captcha

import (
	api2captcha "github.com/2captcha/2captcha-go"
	"log"
	"yuanda-go/pkg/config"
)

func IdentifyCode(filePath string) (string, error) {
	client := api2captcha.NewClient(config.Conf.Captcha.ApiKey)
	captcha := api2captcha.Normal{
		File: filePath,
	}
	code, taskId, err := client.Solve(captcha.ToRequest())
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	log.Printf("Captcha task ID: %s  code : %s", taskId, code)
	return code, nil
}
