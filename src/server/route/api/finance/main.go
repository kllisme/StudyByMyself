package finance

import (
	iris "gopkg.in/kataras/iris.v5"
	"maizuo.com/soda/erp/api/src/server/controller/api/finance"
)

func Setup(v iris.MuxAPI) {
	var (
		billCtrl      = &finance.BillController{}
		dailyBillCtrl = &finance.DailyBillController{}
	)

	api := v.Party("/")
	v1.Post("/bills/export",billCtrl.Export)
	api.Get("/bills", billCtrl.ListByAccountType)
	api.Get("/bills/:id", dailyBillCtrl.ListByBillId)

	api.Get("/daily-bills/:id", dailyBillCtrl.DetailsById)

	api.Post("/settlement/actions/pay", billCtrl.BatchPay)
}
