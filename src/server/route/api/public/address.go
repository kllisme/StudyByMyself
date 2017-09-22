package public

import (
	iris "gopkg.in/kataras/iris.v5"
	publicCtrl "maizuo.com/soda/erp/api/src/server/controller/api/public"
)

func SetupAddress(v iris.MuxAPI) {
	var (
		provinceCtrl = &publicCtrl.ProvinceController{}
		cityCtrl = &publicCtrl.CityController{}
		areaCtrl = &publicCtrl.AreaController{}
		streetCtrl = &publicCtrl.StreetController{}
		schoolCtrl = &publicCtrl.SchoolController{}
	)

	api := v.Party("/addresses")

	api.Get("/provinces", provinceCtrl.Paging)
	api.Post("/provinces", provinceCtrl.Create)
	api.Delete("/provinces/:id", provinceCtrl.Delete)
	api.Put("/provinces/:id", provinceCtrl.Update)
	api.Get("/provinces/:id", provinceCtrl.GetByID)

	api.Get("/cities", cityCtrl.Paging)
	api.Post("/cities", cityCtrl.Create)
	api.Delete("/cities/:id", cityCtrl.Delete)
	api.Put("/cities/:id", cityCtrl.Update)
	api.Get("/cities/:id", cityCtrl.GetByID)

	api.Get("/areas", areaCtrl.Paging)
	api.Post("/areas", areaCtrl.Create)
	api.Delete("/areas/:id", areaCtrl.Delete)
	api.Put("/areas/:id", areaCtrl.Update)
	api.Get("/areas/:id", areaCtrl.GetByID)

	api.Get("/streets", streetCtrl.Paging)
	api.Post("/streets", streetCtrl.Create)
	api.Delete("/streets/:id", streetCtrl.Delete)
	api.Put("/streets/:id", streetCtrl.Update)
	api.Get("/streets/:id", streetCtrl.GetByID)

	api.Get("/schools", schoolCtrl.Paging)
	api.Post("/schools", schoolCtrl.Create)
	api.Delete("/schools/:id", schoolCtrl.Delete)
	api.Put("/schools/:id", schoolCtrl.Update)
	api.Get("/schools/:id", schoolCtrl.GetByID)
}
