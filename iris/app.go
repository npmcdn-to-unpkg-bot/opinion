package main

import (
	"github.com/braintree/manners"
	"github.com/thesyncim/opinion/iris/securestream"
	"github.com/thesyncim/opinion/iris/publisher"
	"github.com/thesyncim/opinion/iris/fakelive"
	"github.com/thesyncim/opinion/iris/article"
	"github.com/kataras/iris"
	"github.com/kataras/iris/plugins/iriscontrol"
	"github.com/kardianos/service"
	"time"
	"github.com/boltdb/bolt"
)

type app struct {
	Quit chan bool
}

func (a *app)run() error{

	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		return err
	}
	authenticator :=publisher.AngularAuth(db)

	options := iris.StationOptions{
		Profile:            true,
		ProfilePath:        iris.DefaultProfilePath,
		Cache:              true,
		CacheMaxItems:      0,
		CacheResetDuration: 5 * time.Minute,
		PathCorrection:     true, //explanation at the end of this chapter
	}
	i:=iris.Custom(options)

	i.Plugin(publisher.NewPublisherPlugin("/fakelive",authenticator,db))
	i.Plugin(article.NewArticlesPlugin("/articles",authenticator,db))
	i.Plugin(securestream.NewSecureStreamPlugin("/tokens","/clients",authenticator,db))
	i.Plugin(fakelive.NewFakelivePlugin("/fakelive",authenticator,db))
	i.Post("/auth/login",publisher.AngularSignIn(db, (&publisher.Publisher{}).FindUser, publisher.NewSha512Password, time.Hour*48))
	opts:=iriscontrol.IrisControlOptions{}
	opts.Port=5555
	opts.Users=map[string]string{}
	opts.Users["thesyncim"]="Kirk1zodiak"
	i.Plugin(iriscontrol.New(opts))

	j := fakelive.RunBackgroundScheduler()
	a.Quit = j.Quit

	return manners.ListenAndServe(":9999", i.Serve())
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