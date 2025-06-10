package config

import (
	"github.com/YuanJey/goconf/pkg/config"
	"github.com/go-rod/rod/lib/proto"
)

var Conf Config

type Config struct {
	Captcha struct {
		ApiKey string `json:"api_key" yaml:"api_key"`
	} `json:"captcha" yaml:"captcha"`
	Cookies []*proto.NetworkCookie
	Path    string `json:"path"`
}

func init() {
	config.UnmarshalConfig(&Conf, "config.yaml")
}
