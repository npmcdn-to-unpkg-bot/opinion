//file sstream_plugin.go
package article

import (
	"github.com/kataras/iris"
	"github.com/boltdb/bolt"
)

type ArticlesPlugin struct {
	container          iris.IPluginContainer
	BaseUrl            string
	Authenticator      iris.HandlerFunc
	ArticlesController *ArticlesController
}

func NewArticlesPlugin(baseURL string, authenticator iris.HandlerFunc,dbb *bolt.DB) *ArticlesPlugin {
	db=dbb
	return &ArticlesPlugin{
		BaseUrl:            baseURL,
		Authenticator:      authenticator,
		ArticlesController: &ArticlesController{},
	}
}

// All plugins must at least implements these 3 functions

func (i *ArticlesPlugin) Activate(container iris.IPluginContainer) error {
	// use the container if you want to register other plugins to the server, yes it's possible a plugin can registers other plugins too.
	// here we set the container in order to use it's printf later at the PostListen.
	err:=createBoltBuckets()
	if err!=nil{
		return err
	}
	i.container = container
	return nil
}

func (i ArticlesPlugin) GetName() string {
	return "Fakelive"
}

func (i ArticlesPlugin) GetDescription() string {
	return "Azorestv Fakelive Manager"
}

//
// Implement our plugin, you can view your inject points - listeners on the /kataras/iris/plugin.go too.
//=
// Implement the PostHandle, because this is what we need now, we need to add a listener after a route is registed to our server so we do:
func (i *ArticlesPlugin) PostHandle(route iris.IRoute) {

}

// PostListen called after the server is started, here you can do a lot of staff
// you have the right to access the whole iris' Station also, here you can add more routes and do anything you want, for example start a second server too, an admin web interface!
// for example let's print to the server's stdout the routes we collected...
func (i *ArticlesPlugin) PostListen(s *iris.Station) {
	article := s.Party(i.BaseUrl)
	article.Use(i.Authenticator)
	article.Post("/create", i.ArticlesController.Create)
	article.Post("/edit/:id", i.ArticlesController.Edit)
	article.Get("/getid/:id", i.ArticlesController.GetId)
	article.Post("/delete/:id", i.ArticlesController.Delete)
	article.Post("/publisher/:id", i.ArticlesController.GetPublisher)
	article.Get("/listall", i.ArticlesController.ListAll)
	article.Get("/listfrontend", i.ArticlesController.ListFrontend)

	articlesFrontend:= s.Party(i.BaseUrl+"f")
	articlesFrontend.Get("/getid/:id", i.ArticlesController.GetId)
	articlesFrontend.Get("/listfrontend", i.ArticlesController.ListFrontend)
}