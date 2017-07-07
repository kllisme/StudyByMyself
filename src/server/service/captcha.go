package service

import (
	"maizuo.com/soda/erp/api/src/server/common"
	"github.com/spf13/viper"
	"time"
	"github.com/levigross/grequests"
	"github.com/bitly/go-simplejson"
)

type CaptchaService struct {
}

func (self *CaptchaService) Count(sessionKey string, functions string, expires int) (int, error) {
	prefix := viper.GetString("resource.captcha.prefix") + functions + ":" + sessionKey
	count, _ := common.Redis.Get(prefix).Int64()
	count++
	_, err := common.Redis.Set(prefix, count, time.Minute * time.Duration(expires)).Result()
	if err != nil {
		return int(count), err
	}
	return int(count), nil
}

func (self *CaptchaService) Del(sessionKey string, functions string) (bool, error) {
	prefix := viper.GetString("resource.captcha.prefix") + functions + ":" + sessionKey
	_, err := common.Redis.Del(prefix).Result()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (self *CaptchaService) GetKey() (string, *common.Result) {
	result, err := grequests.Post(viper.GetString("resource.captcha.server"), &grequests.RequestOptions{})
	if err != nil {
		return "", common.Error("", err)
	}
	switch result.StatusCode {
	case 429:
		return "", common.Error("040001", result.String())
	case 403:
		return "", common.Error("040002", result.String())
	case 410:
		return "", common.Error("040003", result.String())
	}
	if result.StatusCode != 200 || !result.Ok {
		return "", common.Error("040004", result.String())
	}
	return result.String(), nil
}

func (self *CaptchaService) Verity(code string, key string) (bool, *common.Result) {
	result, err := grequests.Post(viper.GetString("resource.captcha.server") + key, &grequests.RequestOptions{
		Data: map[string]string{
			"code": code,
		},
	})
	if err != nil {
		common.Logger.Debugln("err====", err.Error())
		return false, common.Error("", err)
	}

	switch result.StatusCode {
	case 429:
		return false, common.Error("040001", result.String())
	case 403:
		return false, common.Error("040002", result.String())
	case 410:
		return false, common.Error("040003", result.String())
	}
	if result.StatusCode != 200 || !result.Ok {
		return false, common.Error("040004", result.String())
	}
	if result.String() != "Matched" {
		return false, nil
	}
	return true, nil
}

func (self *CaptchaService) Middleware(captcha *simplejson.Json, sessionKey string, functions string, expires int, maxRequest int) (map[string]interface{}, *common.Result) {
	errStatuts := ""
	count, _ := self.Count(sessionKey, functions, expires)
	if count > maxRequest {
		isVeritied := false
		_err := &common.Result{}
		if captcha != nil {
			key := captcha.Get("key").MustString()
			code := captcha.Get("code").MustString()
			isVeritied, _err = self.Verity(code, key)
			if _err != nil {
				if _err.Status == "040004" || _err.Status == "000001"{
					return nil, common.Error("040004", _err)
				} else if _err.Status == "040003" {
					errStatuts = "040003"
				}
			}
			if isVeritied {
				self.Del(sessionKey, functions)
			}
		}
		if captcha == nil || !isVeritied {
			captchaMap := make(map[string]interface{}, 0)
			captchaData := make(map[string]interface{}, 0)
			key, _err := self.GetKey()
			if _err != nil {
				return nil, _err
			}
			captchaMap["key"] = key
			captchaData["captcha"] = captchaMap
			if errStatuts == "040003" {
				return captchaData, common.Error("040003", nil)
			}
			if count == maxRequest + 1 {
				return captchaData, common.Error("040006", nil)
			} else {
				return captchaData, common.Error("040007", nil)
			}
		}
	}
	return nil, nil
}
