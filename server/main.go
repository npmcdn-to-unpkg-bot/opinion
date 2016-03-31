package main

import (
	"time"

	"github.com/braintree/manners"
	"github.com/gin-gonic/gin"
	"github.com/itsjamie/gin-cors"
	fakelive "github.com/thesyncim/opinion/server/fakelivev2"
	"github.com/kardianos/service"
	"log"
	"os"

	"github.com/thesyncim/opinion/server/sstream"
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

func (a *app)run() error{
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
	fake.GET("getsmilplaylist", fakelive.HandlerCurrentSmilPlaylist)

	fake.GET("starttime", fakelive.HandlerGetStartTime)
	fake.GET("livestreamset", fakelive.HandlerGetLiveStreamSettings)
	fake.POST("starttime", fakelive.HandlerSetStartTime).Use(authenticator)
	fake.POST("livestreamset", fakelive.HandlerSetLiveStreamSettings).Use(authenticator)
	fake.POST("reload", fakelive.HandlerReloadNow).Use(authenticator)


	tokens:=router.Group("/tokens")
	sstream.RegisterTokenController(tokens)
	clients:=router.Group("/clients")
	sstream.RegisterClientController(clients)

	j := fakelive.RunBackgroundScheduler()

	a.Quit = j.Quit

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.



	return manners.ListenAndServe(":9999", router)

}

func (a *app) Start(s service.Service) error {

	go a.run()
	return nil
}

func (a *app) Stop(s service.Service) error {
	a.Quit <- true


	manners.Close()

	return nil
}

type Config struct{}

func main() {


	svcConfig := &service.Config{
		Name:        "fakelive",
		DisplayName: "fakelive and opinion server",
		Description: "",
	}

	prg:=&app{Quit: make(chan bool)}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	if len(os.Args) > 1 {
		err = service.Control(s, os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	logger, err := s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}
	err = s.Run()
	if err != nil {
		logger.Error(err)
	}

}
