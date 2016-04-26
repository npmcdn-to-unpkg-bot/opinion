package main

import (
	"github.com/braintree/manners"
	"github.com/kataras/iris"

	"github.com/thesyncim/opinion/iris/article"
	"github.com/thesyncim/opinion/iris/fakelive"
	"github.com/thesyncim/opinion/iris/publisher"
	"github.com/thesyncim/opinion/iris/securestream"
	"github.com/boltdb/bolt"
	"github.com/kardianos/service"
	"time"
	"github.com/kataras/iris/middleware/cors"
	"github.com/kataras/iris/middleware/recovery"
	"os"
	"io"
	"github.com/asdine/storm"
	"log"
)

type app struct {
	Quit    chan bool
	Logfile io.ReadWriteCloser
}

func (a *app) run() error {

	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Println(err)
		return err
	}
	authenticator := publisher.AngularAuth(&storm.DB{Bolt:db})

	a.Logfile, err = os.OpenFile("recovery.log", os.O_APPEND | os.O_WRONLY | os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	iris.Use(recovery.Recovery(a.Logfile))
	iris.Use(cors.New(cors.CorsOptions{AllowCredentials: true}))
	iris.Plugin(publisher.NewPublisherPlugin("/publisher", authenticator, db))
	iris.Plugin(article.NewArticlesPlugin("/article", authenticator, db))
	iris.Plugin(securestream.NewSecureStreamPlugin("/tokens", "/clients", authenticator, db))
	iris.Plugin(fakelive.NewFakelivePlugin("/fakelive", authenticator, db))
	iris.Post("/auth/login", publisher.AngularSignIn(&storm.DB{Bolt:db}, (&publisher.Publisher{}).FindUser, publisher.NewSha512Password, time.Hour * 48))
	iris.Options("/auth/login", func(c *iris.Context) {})

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
	a.Logfile.Close()
	return nil
}
