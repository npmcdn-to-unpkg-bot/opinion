//file sstream_plugin.go
package securestream

import (
	"github.com/kataras/iris"
	"github.com/boltdb/bolt"

	"log"
)

type SecureStreamPlugin struct {
	container        iris.IPluginContainer
	TokenBaseUrl     string
	ClientBaseUrl    string
	Authenticator iris.HandlerFunc
	TokenController  *TokenController
	ClientController *ClientController
}

func NewSecureStreamPlugin(tokenURL, clientURL string, authenticator iris.HandlerFunc,dbb *bolt.DB) *SecureStreamPlugin {
	db=dbb

	err := createBoltBuckets()
	if err != nil {
		log.Fatalln(err)

	}
	return &SecureStreamPlugin{
		TokenBaseUrl:tokenURL,
		ClientBaseUrl:clientURL,
		Authenticator:authenticator,
		TokenController:&TokenController{},
		ClientController:&ClientController{},

	}
}

// All plugins must at least implements these 3 functions

func (i *SecureStreamPlugin) Activate(container iris.IPluginContainer) error {
	// use the container if you want to register other plugins to the server, yes it's possible a plugin can registers other plugins too.
	// here we set the container in order to use it's printf later at the PostListen.
	i.container = container

	return nil
}

func (i SecureStreamPlugin) GetName() string {
	return "SecureStream"
}

func (i SecureStreamPlugin) GetDescription() string {
	return "Secure Stream Plugin to manage Wowza streams and Clients"
}

//
// Implement our plugin, you can view your inject points - listeners on the /kataras/iris/plugin.go too.
//=
// Implement the PostHandle, because this is what we need now, we need to add a listener after a route is registed to our server so we do:
func (i *SecureStreamPlugin) PostHandle(route iris.IRoute) {

}

// PostListen called after the server is started, here you can do a lot of staff
// you have the right to access the whole iris' Station also, here you can add more routes and do anything you want, for example start a second server too, an admin web interface!
// for example let's print to the server's stdout the routes we collected...
func (i *SecureStreamPlugin) PreListen(s *iris.Station) {

	tokens := s.Party(i.TokenBaseUrl)
	tokens.Post("/create", i.TokenController.Create)
	tokens.Post("/update/:id", i.TokenController.Update)
	tokens.Post("/delete/:id", i.TokenController.Delete)
	tokens.Get("/get/:id", i.TokenController.Read)
	tokens.Get("/getall", i.TokenController.ReadAll)

	clients := s.Party(i.ClientBaseUrl)
	clients.Post("/create", i.ClientController.Create)
	clients.Post("/update/:id", i.ClientController.Update)
	clients.Post("/delete/:id", i.ClientController.Delete)
	clients.Get("/get/:id", i.ClientController.Read)
	clients.Get("/getall", i.ClientController.ReadAll)

	//do what ever you want, if you have imagination you can do a lot
}

//