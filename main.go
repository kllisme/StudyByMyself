package main

import (
	"maizuo.com/soda/erp/api/src/server"
	"maizuo.com/soda/erp/api/src/server/common"
)

func main() {

	common.SetupConfig()

	common.SetupCORS()

	common.SetupJWT()

	common.SetupCommon()

	common.SetupLogger()

	common.SetupSession()

	common.SetupRedis()

	server.SetUpServer()
}
