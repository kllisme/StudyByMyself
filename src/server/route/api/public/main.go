package public

import (
	iris "gopkg.in/kataras/iris.v5"
	publicCtrl "maizuo.com/soda/erp/api/src/server/controller/api/public"
)

func Setup(v iris.MuxAPI) {
	var (
		applicationCtrl = &publicCtrl.ApplicationController{}
		advertisementCtrl = &publicCtrl.AdvertisementController{}
		adPositionCtrl = &publicCtrl.ADPositionController{}
	)

	api := v.Party("/")

	SetupAddress(api)

	api.Get("/applications", applicationCtrl.Paging)
	api.Post("/applications", applicationCtrl.Create)
	api.Delete("/applications/:id", applicationCtrl.Delete)
	api.Put("/applications/:id", applicationCtrl.Update)
	api.Get("/applications/:id", applicationCtrl.GetByID)

	api.Get("/advertisements", advertisementCtrl.Paging)
	api.Post("/advertisements", advertisementCtrl.Create)
	api.Delete("/advertisements/:id", advertisementCtrl.Delete)
	api.Put("/advertisements/:id", advertisementCtrl.Update)
	api.Get("/advertisements/:id", advertisementCtrl.GetByID)
	api.Post("/advertisements/images", advertisementCtrl.SaveImage)
	api.Post("/advertisements/batch/orders", advertisementCtrl.BatchUpdateOrder)

	api.Get("/ad-positions", adPositionCtrl.Paging)
	api.Post("/ad-positions", adPositionCtrl.Create)
	api.Delete("/ad-positions/:id", adPositionCtrl.Delete)
	api.Put("/ad-positions/:id", adPositionCtrl.Update)
	api.Get("/ad-positions/:id", adPositionCtrl.GetByID)
}
