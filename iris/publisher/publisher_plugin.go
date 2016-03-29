//file sstream_plugin.go
package publisher

import (
	"github.com/kataras/iris"
	"github.com/boltdb/bolt"
	"log"
)

type PublisherPlugin struct {
	container           iris.IPluginContainer
	BaseUrl             string
	Authenticator       iris.HandlerFunc
	PublisherController *PublisherController
}

func NewPublisherPlugin(baseURL string, authenticator iris.HandlerFunc,dbb *bolt.DB) *PublisherPlugin {
	db=dbb
	err:=createBoltBuckets()
	if err!=nil{
		log.Fatalln(err)
	}
	return &PublisherPlugin{
		BaseUrl:             baseURL,
		Authenticator:       authenticator,
		PublisherController: &PublisherController{},
	}
}

// All plugins must at least implements these 3 functions

func (i *PublisherPlugin) Activate(container iris.IPluginContainer) error {
	// use the container if you want to register other plugins to the server, yes it's possible a plugin can registers other plugins too.
	// here we set the container in order to use it's printf later at the PostListen.

	i.container = container
	return nil
}

func (i PublisherPlugin) GetName() string {
	return "Opinion Publisher"
}

func (i PublisherPlugin) GetDescription() string {
	return "Azorestv Opinion App Publisher"
}

//
// Implement our plugin, you can view your inject points - listeners on the /kataras/iris/plugin.go too.
//=
// Implement the PostHandle, because this is what we need now, we need to add a listener after a route is registed to our server so we do:
func (i *PublisherPlugin) PostHandle(route iris.IRoute) {

}

// PostListen called after the server is started, here you can do a lot of staff
// you have the right to access the whole iris' Station also, here you can add more routes and do anything you want, for example start a second server too, an admin web interface!
// for example let's print to the server's stdout the routes we collected...
func (i *PublisherPlugin) PreListen(s *iris.Station) {
	pub := s.Party(i.BaseUrl)
	pub.Post("/create", i.Authenticator, i.PublisherController.Create)
	pub.Post("/edit/:id", i.Authenticator, i.PublisherController.Edit)
	pub.Get("/getid/:id", i.PublisherController.GetId)
	pub.Post("/delete/:id", i.Authenticator, i.PublisherController.Delete)
	pub.Get("/listall", i.PublisherController.ListAll)
	pub.Get("/publisher/image/:id", i.PublisherController.GetImage)

	i.container.Printf("Plugin routes registereds  %+v",pub )

}
