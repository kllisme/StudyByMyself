package alipay

import (
	"sort"
	"fmt"
	"crypto/x509"
	"crypto/rsa"
	"crypto/sha1"
	"crypto"
	"encoding/base64"
	"encoding/hex"
	"strings"
	"github.com/spf13/viper"
	"maizuo.com/soda/erp/api/src/server/common"
)

type AlipayKit struct {
}

func StringToSign(m interface{}, isVerify bool) string{
	//STEP 1, 对key进行升序排序.
	sorted_keys := make([]string, 0)
	if m == nil {
		return ""
	}

	switch mValue := m.(type) {
	case map[string]interface{}:
		for k, _ := range mValue {
			sorted_keys = append(sorted_keys, k)
		}
	case map[string]string:
		for k, _ := range mValue {
			sorted_keys = append(sorted_keys, k)
		}
	default:
		return ""
	}

	sort.Strings(sorted_keys)
	//STEP2, 对key=value的键值对用&连接起来，略过空值
	var signStrings string
	for _, k := range sorted_keys {
		value := ""
		switch mValue := m.(type) {
		case map[string]interface{}:
			//fmt.Printf("k=%v, v=%v\n", k, mValue[k])
			value = fmt.Sprintf("%v", mValue[k])
		case map[string]string:
			//fmt.Printf("k=%v, v=%v\n", k, mValue[k])
			value = fmt.Sprintf("%v", mValue[k])
		}


		if value != "" && k != "sign" {
			if isVerify && k == "sign_type" {
				continue
			}
			signStrings = signStrings + k + "=" + value + "&"
		}
	}

	if strings.HasSuffix(signStrings, "&") {
		signStrings = signStrings[: len(signStrings) - 1 ]
	}

	//fmt.Println("signStrings=============", signStrings)
	return signStrings
}

func (self *AlipayKit) CreateRsaSign(mReq map[string]interface{}) (string) {
	signStrings := StringToSign(mReq, false)
	common.Logger.Debugln("生成的待签名---->", signStrings)
	//============================= 开始签名 ==================================
	key := viper.GetString("resource.pay.alipay.privateKey")

	common.Logger.Debugln("key=================", `-----BEGIN RSA PRIVATE KEY-----` + key + `----END RSA PRIVATE KEY-----`)
	encodedKey, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		common.Logger.Debugln("rsaSign private_key error", err.Error())
		return ""
	}
	privateKey, err := x509.ParsePKCS8PrivateKey(encodedKey)
	if err != nil {
		common.Logger.Debugln("x509.ParsePKCS1PrivateKey-------privateKey----- error : %v\n", err)
		return ""
	} else {
		common.Logger.Debugln("x509.ParsePKCS1PrivateKey-------privateKey-----", privateKey)
	}

	if privateKey == nil {
		return ""
	}

	result, err := RsaSign(signStrings, privateKey.(*rsa.PrivateKey))
	common.Logger.Debugln("alipay.RsaSign=========", result, err)
	return result
}

/**
 * RSA签名
 */
func RsaSign(origData string, privateKey *rsa.PrivateKey) (string, error) {
	h := sha1.New()
	h.Write([]byte(origData))
	digest := h.Sum(nil)
	s, err := rsa.SignPKCS1v15(nil, privateKey, crypto.SHA1, digest)
	if err != nil {
		fmt.Errorf("rsaSign SignPKCS1v15 error")
		return "", err
	}
	data := base64.StdEncoding.EncodeToString(s)
	return string(data), nil
}

/**
 * RSA签名验证
 */
func (self *AlipayKit) VerifyRsaSign(m interface{}, sign string) (bool, error) {
	src := StringToSign(m, true)
	fmt.Println("src=====", src)
	//步骤1，加载RSA的公钥
	key := viper.GetString("resource.pay.alipay.publicKey")
	//block, _ := pem.Decode([]byte(`-----BEGIN RSA PRIVATE KEY-----` + key + `----END RSA PRIVATE KEY-----`))
	common.Logger.Warningln("==================begin===============")
	// 解base64
	encodedKey, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		common.Logger.Warningln("==============Failed to decode RSA public key: %s\n", err.Error())
		return false, err
	}
	pub, err := x509.ParsePKIXPublicKey(encodedKey)
	//pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		common.Logger.Warningln("================Failed to parse RSA public key: %s\n", err.Error())
		return false, err
	}
	common.Logger.Warningln("==================end===============")
	rsaPub, _ := pub.(*rsa.PublicKey)
	RsaVerify(src, rsaPub, sign)


	return true, nil
}

func RsaVerify(origData string, publicKey *rsa.PublicKey, sign string) (bool, error) {
	//步骤2，计算代签名字串的SHA1哈希
	h := sha1.New()
	//io.WriteString(h, src)
	h.Write([]byte(origData))
	digest := h.Sum(nil)

	//步骤3，base64 decode,必须步骤，支付宝对返回的签名做过base64 encode必须要反过来decode才能通过验证
	data, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return false, err
	}

	hexSig := hex.EncodeToString(data)
	common.Logger.Debugln("base decoder: %v, %v\n", sign, hexSig)

	//步骤4，调用rsa包的VerifyPKCS1v15验证签名有效性
	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA1, digest, data)
	if err != nil {
		common.Logger.Warningln("Verify sig error, reason: ", err)
		return false, err
	}
	return true, nil
}

