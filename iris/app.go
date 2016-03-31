package main

import (
	"github.com/braintree/manners"
	"github.com/kataras/iris"

	//"github.com/kataras/iris/plugins/iriscontrol"
	"github.com/thesyncim/opinion/iris/article"
	"github.com/thesyncim/opinion/iris/fakelive"
	"github.com/thesyncim/opinion/iris/publisher"
	"github.com/thesyncim/opinion/iris/securestream"

	"github.com/boltdb/bolt"
	"github.com/kardianos/service"
	"log"
	"time"

	"github.com/kataras/iris/middleware/cors"
	"github.com/kataras/iris/plugins/routesinfo"
)

type app struct {
	Quit chan bool
}

func (a *app) run() error {

	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		return err
	}
	authenticator := publisher.AngularAuth(db)

	info := routesinfo.RoutesInfo()
	iris.Plugin(info)

	iris.Use(cors.New(cors.CorsOptions{AllowCredentials: true}))

	err = iris.Plugin(publisher.NewPublisherPlugin("/publisher", authenticator, db))
	if err != nil {
		log.Fatalln(err)
	}
	err = iris.Plugin(article.NewArticlesPlugin("/article", authenticator, db))
	if err != nil {
		log.Fatalln(err)
	}
	err = iris.Plugin(securestream.NewSecureStreamPlugin("/tokens", "/clients", authenticator, db))
	if err != nil {
		log.Fatalln(err)
	}
	err = iris.Plugin(fakelive.NewFakelivePlugin("/fakelive", authenticator, db))
	if err != nil {
		log.Fatalln(err)
	}
	iris.Post("/auth/login", publisher.AngularSignIn(db, (&publisher.Publisher{}).FindUser, publisher.NewSha512Password, time.Hour*48))
	all := info.All()

	println("The first registed route was: ", all[0].Path, "registed at: ", all[0].RegistedAt.String())
	println("All routes info:")
	for i := 0; i < len(all); i++ {
		println(all[i].String())
		//outputs->
		// Domain: localhost Method: GET Path: /yourpath RegistedAt: 2016/03/27 15:27:05:029 ...
		// Domain: localhost Method: POST Path: /otherpostpath RegistedAt: 2016/03/27 15:27:05:030 ...
	}

	j := fakelive.RunBackgroundScheduler()
	a.Quit = j.Quit

	return manners.ListenAndServe(":9999", iris.Serve())
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
