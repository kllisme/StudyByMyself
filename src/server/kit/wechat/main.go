package wechat

import (
	"github.com/levigross/grequests"
	"maizuo.com/soda/erp-api/src/server/kit/functions"
	"maizuo.com/soda/erp-api/src/server/common"
	"strings"
	"github.com/spf13/viper"
	"strconv"
	"unicode/utf8"
	"regexp"
)

func GetAccess(appId string, code string) (map[string]interface{}, error) {
	href:=viper.GetString("resource.wechat.href.api")
	url := href+"/sns/oauth2/access_token?appid=APPID&secret=SECRET&code=CODE&grant_type=authorization_code"
	secret := viper.GetString("resource.pay.wechat.secret")
	common.Logger.Debugln("secret====", secret)
	url = strings.Replace(url, "APPID", appId, -1)
	url = strings.Replace(url, "SECRET", secret, -1)
	url = strings.Replace(url, "CODE", code, -1)
	common.Logger.Debugln("url===========", url)
	response, err := grequests.Get(url, &grequests.RequestOptions{})
	if err != nil {
		return nil, err
	}
	return functions.ResponseMap(response, "wechat response error")

}

func GetUserInfo(accessToken string, openId string) (map[string]interface{}, error) {
	href:=viper.GetString("resource.wechat.href.api")
	url :=href+ "/sns/userinfo?access_token=ACCESS_TOKEN&openid=OPENID&lang=zh_CN"
	url = strings.Replace(url, "ACCESS_TOKEN", accessToken, -1)
	url = strings.Replace(url, "OPENID", openId, -1)
	response, err := grequests.Get(url, &grequests.RequestOptions{})
	if err != nil {
		return nil, err
	}
	return functions.ResponseMap(response, "wechat response error")
}

//表情解码
func UnicodeEmojiDecode(s string) string {
	//emoji表情的数据表达式
	re := regexp.MustCompile("\\[[\\\\u0-9a-zA-Z]+\\]")
	//提取emoji数据表达式
	reg := regexp.MustCompile("\\[\\\\u|]")
	src := re.FindAllString(s, -1)
	for i := 0; i < len(src); i++ {
		e := reg.ReplaceAllString(src[i], "")
		p, err := strconv.ParseInt(e, 16, 32)
		if err == nil {
			s = strings.Replace(s, src[i], string(rune(p)), -1)
		}
	}
	return s
}

//表情转换
func UnicodeEmojiCode(s string) string {
	ret := ""
	rs := []rune(s)
	for i := 0; i < len(rs); i++ {
		if len(string(rs[i])) == 4 {
			u := `[\u` + strconv.FormatInt(int64(rs[i]), 16) + `]`
			ret += u

		} else {
			ret += string(rs[i])
		}
	}
	return ret
}


func FilterEmoji(content string) string {
	new_content := ""
	for _, value := range content {
		_, size := utf8.DecodeRuneInString(string(value))
		if size <= 3 {
			new_content += string(value)
		}
	}
	return new_content
}
