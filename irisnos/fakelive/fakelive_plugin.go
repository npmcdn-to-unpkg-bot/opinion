//file sstream_plugin.go
package fakelive

import (
	"github.com/boltdb/bolt"
	"github.com/kataras/iris"
	"log"
	"github.com/asdine/storm"
)

type FakelivePlugin struct {
	container          iris.IPluginContainer
	BaseUrl            string
	Authenticator      iris.HandlerFunc
	FakeliveController *FakeliveController
}

func NewFakelivePlugin(baseURL string, authenticator iris.HandlerFunc, dbb *bolt.DB) *FakelivePlugin {
	db = dbb
	stormdb=&storm.DB{Bolt:db}
	stormdb.Init(Video{})
	err := createBoltBuckets()
	if err != nil {
		log.Fatalln(err)
	}
	return &FakelivePlugin{
		BaseUrl:            baseURL,
		Authenticator:      authenticator,
		FakeliveController: &FakeliveController{},
	}
}

// All plugins must at least implements these 3 functions

func (i *FakelivePlugin) Activate(container iris.IPluginContainer) error {
	// use the container if you want to register other plugins to the server, yes it's possible a plugin can registers other plugins too.
	// here we set the container in order to use it's printf later at the PostListen.

	i.container = container
	return nil
}

func (i FakelivePlugin) GetName() string {
	return "Fakelive"
}

func (i FakelivePlugin) GetDescription() string {
	return "Azorestv Fakelive Manager"
}

//
// Implement our plugin, you can view your inject points - listeners on the /kataras/iris/plugin.go too.
//=
// Implement the PostHandle, because this is what we need now, we need to add a listener after a route is registed to our server so we do:
func (i *FakelivePlugin) PostHandle(route iris.IRoute) {

}

// PostListen called after the server is started, here you can do a lot of staff
// you have the right to access the whole iris' Station also, here you can add more routes and do anything you want, for example start a second server too, an admin web interface!
// for example let's print to the server's stdout the routes we collected...
func (i *FakelivePlugin) PreListen(s *iris.Station) {
	fake := s.Party(i.BaseUrl)
	fake.Get("/getplaylist", i.FakeliveController.CurrentPlaylist)
	fake.Get("/getsmilplaylist", i.FakeliveController.CurrentSmilPlaylist)

	fake.Post("/settings", i.Authenticator, i.FakeliveController.SetSettings)
	fake.Get("/settings", i.Authenticator, i.FakeliveController.GetSettings)


	fake.Get("/trim/new", i.Authenticator,i.FakeliveController.GetNewTrim)
	fake.Get("/trim/search/:keyword",i.Authenticator,i.FakeliveController.GetSearchTrim)
	fake.Post("/trim/save",i.Authenticator,i.FakeliveController.PostSaveVideoTrim)

	fake.Post("/reload", i.Authenticator, i.FakeliveController.ReloadNow)

	i.container.Printf("Plugin fakelive registered \n")
}
