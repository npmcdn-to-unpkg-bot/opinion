package main

import (
	"time"

	"github.com/braintree/manners"
	"github.com/gin-gonic/gin"
	"github.com/itsjamie/gin-cors"
	"github.com/thesyncim/opinion/server/fakelive"
	"gopkg.in/hlandau/easyconfig.v1"
	"gopkg.in/hlandau/service.v2"
)

/*
type Runnable interface {
	Start() error
	Stop() error
}
*/

type app struct {
	Quit chan bool
}

func (a *app) Start() error {

	router := gin.Default()
	router.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))

	var publisher = PublisherController{}
	var articles = ArticlesController{}

	authenticator := AngularAuth(db)
	publisherrouter := router.Group("/publisher")
	//publisher.Use(authenticator)
	{
		publisherrouter.Any("/create", publisher.Create)
		publisherrouter.POST("/edit/:id", publisher.Edit)
		publisherrouter.GET("/getid/:id", publisher.GetId)

		publisherrouter.POST("/delete/:id", publisher.Delete)
		publisherrouter.GET("/listall", publisher.ListAll)

	}

	router.GET("/publisher/image/:id", publisher.GetImage)

	article := router.Group("/article")
	article.Use(authenticator)
	{
		article.POST("/create", articles.Create)
		article.POST("/edit/:id", articles.Edit)
		article.GET("/getid/:id", articles.GetId)
		article.POST("/delete/:id", articles.Delete)
		article.POST("/publisher/:id", articles.GetPublisher)
		article.GET("/listall", articles.ListAll)
		article.GET("/listfrontend", articles.ListFrontend)

	}

	articleFrontEnd := router.Group("/articlef")

	{

		articleFrontEnd.GET("/getid/:id", articles.GetId)

		articleFrontEnd.GET("/listfrontend", articles.ListFrontend)

	}

	auth := router.Group("/auth")
	{
		auth.POST("/login", AngularSignIn(db, (&Publisher{}).FindUser, NewSha512Password, time.Hour*48))
	}

	fake := router.Group("/fakelive")

	fake.GET("getplaylist", fakelive.HandlerCurrentPlaylist)

	fake.GET("starttime", fakelive.HandlerGetStartTime)
	fake.GET("livestreamset", fakelive.HandlerGetLiveStreamSettings)
	fake.POST("starttime", fakelive.HandlerSetStartTime).Use(authenticator)
	fake.POST("livestreamset", fakelive.HandlerSetLiveStreamSettings).Use(authenticator)
	fake.POST("reload", fakelive.HandlerReloadNow).Use(authenticator)

	j := fakelive.RunBackgroundScheduler()

	a.Quit = j.Quit

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.



	return manners.ListenAndServe(":9999", router)
}

func (a *app) Stop() error {
	a.Quit <- true

	manners.Close()

	return nil
}

type Config struct{}

func main() {

	cfg := Config{}

	(&easyconfig.Configurator{
		ProgramName: "Azorestv Software",
	}).ParseFatal(&cfg)

	service.Main(&service.Info{
		Name: "Azorestv Software",

		NewFunc: func() (service.Runnable, error) {
			return &app{Quit: make(chan bool)},nil
		},
		AllowRoot:true,
	})

}
