package order

import (
	"time"
	"fmt"
	"math/rand"
	"strings"
)

//type Payment int
const (
	UNKNOW int = iota
	WECHAT
	ALIPAY
)

func GetPayment(payment string) int {
	var p int
	_payment := strings.ToUpper(payment)
	switch _payment {
	case "WECHAT_OA":
		p = WECHAT
	case "ALIPAY_WEB":
		p = ALIPAY
	default:
		p = UNKNOW
	}
	return p
}

func GetPaymentName(payment string) string {
	name := ""
	switch payment {
	case "WECHAT_OA":
		name = "微信"
	case "ALIPAY_WEB":
		name = "支付宝"
	default:
		name = "未知"
	}
	return name
}

//UNPAID 0等待支付，PAY_TIMEOUT 1支付超时，PAY_FAILURE 2支付失败，PAID 3支付完成，REFUNDED 4已退款，DELIVERY_TIMEOUT 5发货超时，
//DELIVERY_FAILURE 6发货失败，DELIVERED 7已发货，CANCELLED 8已取消, DELIVERY_FAILURE_REFUNDED 9退款失败且已退款， DELIVERY_TIMEOUT_REFUNDED 10退款超时且已退款
//type PayStatus int
const (
	UNPAID int = iota
	PAY_TIMEOUT
	PAY_FAILURE
	PAID
	REFUNDED
	DELIVERY_TIMEOUT
	DELIVERY_FAILURE
	DELIVERED
	CANCELLED
	DELIVERY_FAILURE_REFUNDED
	DELIVERY_TIMEOUT_REFUNDED
)

func GetPayStatus(status int32) string {
	s := ""
	_status := int(status)
	switch _status {
	case UNPAID:
		s = "UNPAID"
	case PAY_TIMEOUT:
		s = "PAY_TIMEOUT"
	case PAY_FAILURE:
		s = "PAY_FAILURE"
	case PAID:
		s = "PAID"
	case REFUNDED:
		s = "REFUNDED"
	case DELIVERY_TIMEOUT:
		s = "DELIVERY_TIMEOUT"
	case DELIVERY_FAILURE:
		s = "DELIVERY_FAILURE"
	case DELIVERED:
		s = "DELIVERED"
	case CANCELLED:
		s = "CANCELLED"
	case DELIVERY_FAILURE_REFUNDED:
		s = "DELIVERY_FAILURE_REFUNDED"
	case DELIVERY_TIMEOUT_REFUNDED:
		s = "DELIVERY_TIMEOUT_REFUNDED"
	}
	return s
}

func GenerateIdByMobile(mobile string) string {
	if len(mobile) != 11 {
		return ""
	}
	prefix := mobile[len(mobile)-4:]
	ymd := time.Now().Format("060102")
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := fmt.Sprintf("%06v", rnd.Int31n(1000000))
	id := ymd + prefix + code
	return id
}
