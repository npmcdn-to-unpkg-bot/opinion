package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/plugins/routesinfo"
	"github.com/kataras/iris/plugins/iriscontrol"
)

func main() {

	info := routesinfo.RoutesInfo()
	iris.Plugin(info)


	iris.Plugin(iriscontrol.New(iriscontrol.IrisControlOptions{
		Port:5555,
		Users:map[string]string{
			"teste":"test",
		},

	}))

	iris.Get("/yourpath", func(c *iris.Context) {
		c.Write("yourpath")
	})

	iris.Post("/otherpostpath", func(c *iris.Context) {
		c.Write("other post path")
	})

	all := info.All()
	// allget := info.ByMethod("GET") -> slice
	// alllocalhost := info.ByDomain("localhost") -> slice
	// bypath:= info.ByPath("/yourpath") -> slice
	// bydomainandmethod:= info.ByDomainAndMethod("localhost","GET") -> slice
	// bymethodandpath:= info.ByMethodAndPath("GET","/yourpath") -> single (it could be slice for all domains too but it's not)

	println("The first registed route was: ", all[0].Path, "registed at: ", all[0].RegistedAt.String())
	println("All routes info:")
	for i:=0; i<len(all); i ++ {
		println(all[i].String())
		//outputs->
		// Domain: localhost Method: GET Path: /yourpath RegistedAt: 2016/03/27 15:27:05:029 ...
		// Domain: localhost Method: POST Path: /otherpostpath RegistedAt: 2016/03/27 15:27:05:030 ...
	}
	iris.Listen(":8080")

}

