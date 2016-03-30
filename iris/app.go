package main

import (
	"github.com/braintree/manners"
	"github.com/kataras/iris"
	"github.com/kataras/iris/plugins/iriscontrol"
	"github.com/thesyncim/opinion/iris/article"
	"github.com/thesyncim/opinion/iris/fakelive"
	"github.com/thesyncim/opinion/iris/publisher"
	"github.com/thesyncim/opinion/iris/securestream"

	"github.com/boltdb/bolt"
	"github.com/kardianos/service"
	"log"
	"time"
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
	opts := iriscontrol.IrisControlOptions{}
	opts.Port = 5555
	opts.Users = map[string]string{}
	opts.Users["thesyncim"] = "Kirk1zodiak"
	err = iris.Plugin(iriscontrol.New(opts))
	if err != nil {
		log.Fatalln(err)
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
