package crm

import (
	iris "gopkg.in/kataras/iris.v5"
	"maizuo.com/soda/erp/api/src/server/controller/api/crm"
)

func Setup(v iris.MuxAPI)  {
	var (
		consumeCtrl = &crm.ConsumptionController{}
		deviceCtrl = &crm.DeviceController{}
		customerCtrl = &crm.CustomerController{}
		operatorCtrl = &crm.OperatorController{}
		billCtrl = &crm.BillController{}
	)

	api := v.Party("/crm")
	api.Get("/consumptions/excel/export",consumeCtrl.Export)
	api.Get("/consumptions", consumeCtrl.Paging)
	api.Post("/consumptions/:ticketId/refund", consumeCtrl.Refund)

	api.Get("/devices", deviceCtrl.Paging)
	api.Put("/devices/batch/reset", deviceCtrl.Reset)
	api.Post("/devices/:id/status", deviceCtrl.UpdateStatus)

	api.Get("/customers/:mobile", customerCtrl.GetByMobile)
	api.Put("/customers/:mobile/password", customerCtrl.ChangePassword)

	//api.Get("/wallets/:mobile", customerCtrl.GetWalletByMobile)
	//api.Get("/chipcards/:mobile", customerCtrl.GetchipcardByMobile)
	//api.Get("/tickets/:mobile", customerCtrl.GetTicketByMobile)
	//
	api.Get("/bills", billCtrl.Paging)
	api.Get("/chipcard-bills", billCtrl.PagingChipcardBill)
	//
	api.Get("/operators", operatorCtrl.Paging)
	api.Get("/operators/:id", operatorCtrl.GetByID)
	api.Put("/operators/:id/password", operatorCtrl.ChangePassword)

}
