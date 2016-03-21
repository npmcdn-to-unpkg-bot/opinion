package sstream

import (
	"github.com/gin-gonic/gin"
)

type TokenController struct{
	router *gin.RouterGroup
}

func RegisterTokenController(router *gin.RouterGroup){
	var t = new(TokenController)
	router.POST("/create",t.Create)
	router.POST("/update/:id",t.Update)
	router.POST("/delete/:id",t.Delete)
	router.GET("/get/:id",t.Read)
	router.GET("/getall",t.ReadAll)

}

func (tc *TokenController)Create(c *gin.Context)  {
	var Token= &Token{}
	err:=c.BindJSON(Token)
	if err!=nil{
		c.Error(err)
		return
	}

	err=Token.Save()
	if err!=nil{
		c.Error(err)
		return
	}
}

func (tc *TokenController)Read(c *gin.Context)  {
	id:=c.Param("id")
	var Token= &Token{}

	c.JSON(200,Token.Get(id))
}

func (tc *TokenController)ReadAll(c *gin.Context)  {
	var Token = &Token{}

	Tokens,err:=Token.GetAll()
	if err!=nil{
		c.Error(err)
		return
	}
	c.JSON(200,Tokens)
}

func (tc *TokenController)Update(c *gin.Context)  {
	var Token= &Token{}
	err:=c.BindJSON(Token)
	if err!=nil{
		c.Error(err)
		return
	}

	err=Token.Update()
	if err!=nil{
		c.Error(err)
		return
	}
}

func (tc *TokenController)Delete(c *gin.Context)  {
	id:=c.Param("id")
	var Token= &Token{}

	err:=Token.Delete(id)
	if err!=nil{
		c.Error(err)
		return
	}
}