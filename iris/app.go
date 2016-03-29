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
)

type app struct {
	Quit chan bool
}

func (a *app)run() error{
	authenticator :=publisher.AngularAuth(db)

	iris.Plugin(publisher.NewPublisherPlugin("/fakelive",authenticator,db))
	iris.Plugin(article.NewArticlesPlugin("/articles",authenticator,db))
	iris.Plugin(securestream.NewSecureStreamPlugin("/tokens","/clients",authenticator,db))
	iris.Plugin(fakelive.NewFakelivePlugin("/fakelive",authenticator,db))
	iris.Post("/auth/login",publisher.AngularSignIn(db, (&publisher.Publisher{}).FindUser, publisher.NewSha512Password, time.Hour*48))
	opts:=iriscontrol.IrisControlOptions{}
	opts.Port=5555
	opts.Users=make(map[string]string)
	opts.Users["thesyncim"]="Kirk1zodiak"
	iris.Plugin(iriscontrol.New(opts))

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