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
	api.Put("/consumptions/:ticketId", consumeCtrl.Refund)

	api.Get("/devices", deviceCtrl.Paging)
	api.Put("/devices/batch/remove", deviceCtrl.Remove)
	api.Put("/devices/:id/lock", deviceCtrl.Lock)
	api.Put("/devices/:id/unlock", deviceCtrl.Unlock)
	api.Put("/devices/:id/free", deviceCtrl.Free)
}
