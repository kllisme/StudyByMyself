package public

import (
iris "gopkg.in/kataras/iris.v5"
"maizuo.com/soda/erp/api/src/server/controller/api/public"
)

func Setup(v iris.MuxAPI) {
	var (
		regionCtrl      = &public.RegionController{}
		//schoolCtrl = &public.SchoolController{}
	)

	api := v.Party("/")

	api.Get("/regions/provinces", regionCtrl.GetProvinces)
	api.Get("/regions/provinces/:id/cities", regionCtrl.GetCities)


}
