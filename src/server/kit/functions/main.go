package functions

import (
	"github.com/bitly/go-simplejson"
	"github.com/levigross/grequests"
	"math/rand"
	"os"
	"strconv"
	"time"
	"regexp"
	"strings"
)

//int数组去重
func Uniq(ls []int) []int {
	intInSlice := func(i int, list []int) bool {
		for _, v := range list {
			if v == i {
				return true
			}
		}
		return false
	}
	var Uniq []int
	for _, v := range ls {
		if !intInSlice(v, Uniq) {
			Uniq = append(Uniq, v)
		}
	}
	return Uniq
}

//在数组中查找,找到返回index 没找到返回-1
func FindIndex(ls []int, value int) int {
	for k, v := range ls {
		if value == v {
			return k
		}
	}
	return -1
}

func StringToInt(value string) int {
	v, e := strconv.Atoi(value)
	if e != nil {
		return -1
	}
	return v
}

func StringToFloat64(value string) float64 {
	v, e := strconv.ParseFloat(value, 10)
	if e != nil {
		return -1
	}
	return v
}

/**
prec: -1 代表输出的精度小数点后的位数，如果是<0的值，则返回最少的位数来表示该数，如果是大于0的则返回对应位数的值
*/
func Float64ToString(num float64, prec int) string {
	return strconv.FormatFloat(num, 'f', prec, 64)
}

func Int64ToString(num int64) string {
	//10为十进制
	return strconv.FormatInt(num, 10)
}

func FormatFloat(num float64, prec int) float64 {
	s := strconv.FormatFloat(num, 'f', prec, 64)
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return float64(-1)
	}
	return f
}

func CreatePathIfNotExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		err = os.Mkdir(path, os.ModePerm)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, err
}

func IntToBool(i int) bool {
	if i == 1 {
		return true
	}
	return false
}

func RandInt64(min, max int64) int64 {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Int63n(max - min) + min
}

func ResponseMap(response *grequests.Response, errMsg string) (map[string]interface{}, error) {
	if response.StatusCode != 200 || !response.Ok {
		e := &DefinedError{}
		e.Msg = errMsg
		err := e
		return nil, err
	}
	_json, err := simplejson.NewJson(response.Bytes())
	if err != nil {
		return nil, err
	}
	respMap, err := _json.Map()
	if err != nil {
		return nil, err
	}
	return respMap, nil
}

func ExtractHandlerName(handlerName string) string {
	reg := regexp.MustCompile(`([\w]+)`)
	nameList := reg.FindAllString(handlerName, -1)
	for idx, name := range nameList {
		if name == "controller" {
			nameList = nameList[idx + 1:len(nameList) - 1]
			break
		}
	}
	return strings.Join(nameList, `_`)
}

//CountRune 统计字符串中的字符数量，可以解决len()无法将一个中文汉字识别为1的问题 TODO 可用utf8.RuneCountInString()来代替
func CountRune(str string) int {
	runes := []rune(str)
	return len(runes)
}
