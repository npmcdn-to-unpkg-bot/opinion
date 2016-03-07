package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/itsjamie/gin-cors"
)

func main() {

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

	var p = PublisherController{}
	var a = ArticlesController{}

	authenticator := AngularAuth(db)
	publisher := router.Group("/publisher")
	//publisher.Use(authenticator)
	{
		publisher.Any("/create", p.Create)
		publisher.POST("/edit/:id", p.Edit)
		publisher.GET("/getid/:id", p.GetId)
		publisher.POST("/delete/:id", p.Delete)
		publisher.GET("/listall", p.ListAll)

	}

	article := router.Group("/article")
	article.Use(authenticator)
	{
		article.POST("/create", a.Create)
		article.POST("/edit/:id", a.Edit)
		article.GET("/getid/:id", a.GetId)
		article.POST("/delete/:id", a.Delete)
		article.POST("/publisher/:id", a.GetPublisher)
		article.GET("/listall", a.ListAll)
		article.GET("/listfrontend", a.ListFrontend)

	}

	articlefront := router.Group("/articlef")

	{

		articlefront.GET("/getid/:id", a.GetId)

		articlefront.GET("/listfrontend", a.ListFrontend)

	}

	auth := router.Group("/auth")
	{
		auth.POST("/login", AngularSignIn(db, (&Publisher{}).FindUser, NewSha512Password, time.Hour*48))
	}

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	router.Run(":9999")

}
