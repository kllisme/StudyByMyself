package public

import (
	iris "gopkg.in/kataras/iris.v5"
	"maizuo.com/soda/erp/api/src/server/controller/api/public"
)

func Setup(v iris.MuxAPI) {
	var (
		regionCtrl = &public.RegionController{}
		//schoolCtrl = &public.SchoolController{}
		applicationCtrl = &public.ApplicationController{}
		advertisementCtrl = &public.AdvertisementController{}
		adSpaceCtrl = &public.ADSpaceController{}
	)

	api := v.Party("/")

	api.Get("/regions/provinces", regionCtrl.GetProvinces)
	api.Get("/regions/provinces/:id/cities", regionCtrl.GetCities)

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

	api.Get("/ad-spaces", adSpaceCtrl.Paging)
	api.Post("/ad-spaces", adSpaceCtrl.Create)
	api.Delete("/ad-spaces/:id", adSpaceCtrl.Delete)
	api.Put("/ad-spaces/:id", adSpaceCtrl.Update)
	api.Get("/ad-spaces/:id", adSpaceCtrl.GetByID)
}
