//file sstream_plugin.go
package publisher

import (
	"github.com/boltdb/bolt"
	"github.com/kataras/iris"


	"github.com/asdine/storm"
)

type PublisherPlugin struct {
	container           iris.IPluginContainer
	BaseUrl             string
	Authenticator       iris.HandlerFunc
	PublisherController *PublisherController
}

func NewPublisherPlugin(baseURL string, authenticator iris.HandlerFunc, dbb *bolt.DB) *PublisherPlugin {
	db = dbb
	stormdb=&storm.DB{Bolt:db}
	stormdb.Init(Publisher{})

AddDefaultPub()

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

func (i *PublisherPlugin) PostHandle(route iris.IRoute) {

}

func (i *PublisherPlugin) PreListen(s *iris.Station) {
	pub := s.Party(i.BaseUrl)
	pub.Post("/create", i.PublisherController.Create)
	pub.Post("/edit/:id", i.PublisherController.Edit)
	pub.Get("/getid/:id", i.PublisherController.GetId)
	pub.Post("/delete/:id", i.PublisherController.Delete)
	pub.Get("/listall", i.PublisherController.ListAll)
	pub.Get("/publisher/image/:id", i.PublisherController.GetImage)

	i.container.Printf("Plugin publisher registered \n")

}
