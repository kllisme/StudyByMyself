package crm

import (
	iris "gopkg.in/kataras/iris.v5"
	"maizuo.com/soda/erp/api/src/server/controller/api/crm"
)

func Setup(v iris.MuxAPI)  {
	var (
		consumeCtrl = &crm.ConsumptionController{}
		deviceCtrl = &crm.DeviceController{}

	)

	api := v.Party("/")
	api.Get("/consumptions/excel/export",consumeCtrl.Export)
	api.Get("/consumptions", consumeCtrl.Paging)
	api.Post("/consumptions/:ticketId/refund", consumeCtrl.Refund)

	api.Get("/devices", deviceCtrl.Paging)
	api.Put("/devices/batch/reset", deviceCtrl.Reset)
	api.Post("/devices/:id/status", deviceCtrl.UpdateStatus)
}
