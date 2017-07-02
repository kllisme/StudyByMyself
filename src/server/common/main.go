package common

var (
	common_msg = map[string]string{
		"UNAUTHORIZED":          "会话信息无效",
		"INTERNAL_SERVER_ERROR": "后台系统异常,请重试或稍后再试!",
		"NOT_FOUND":             "你所请求的API不存在,请检查后再试!",
		"TOO_MANY_REQUESTS":     "你所请求的API超过频率限制,请稍后再试!",
	}
)

func SetupCommon() {

}

var ()
