package sstream

import "github.com/gin-gonic/gin"

type ClientController struct{
	router *gin.RouterGroup
}

func RegisterClientController(router *gin.RouterGroup){
	var c = new(ClientController)
	router.POST("/create",c.Create)
	router.POST("/update/:id",c.Update)
	router.POST("/delete/:id",c.Delete)
	router.GET("/get/:id",c.Read)
	router.GET("/getall",c.ReadAll)

}

func (cc *ClientController)Create(c *gin.Context)  {
	var client= &Client{}
	err:=c.BindJSON(client)
	if err!=nil{
		c.Error(err)
		return
	}

	err=client.Save()
	if err!=nil{
		c.Error(err)
		return
	}

}

func (cc *ClientController)Read(c *gin.Context)  {
	id:=c.Param("id")
	var client= &Client{}

	c.JSON(200,client.Get(id))
}

func (cc *ClientController)ReadAll(c *gin.Context)  {
	var client = &Client{}

	clients,err:=client.GetAll()
	if err!=nil{
		c.Error(err)
		return
	}
	c.JSON(200,clients)
}

func (cc *ClientController)Update(c *gin.Context)  {
	var client= &Client{}
	err:=c.BindJSON(client)
	if err!=nil{
		c.Error(err)
		return
	}

	err=client.Update()
	if err!=nil{
		c.Error(err)
		return
	}
}

func (cc *ClientController)Delete(c *gin.Context)  {
	id:=c.Param("id")
	var client= &Client{}

	err:=client.Delete(id)
	if err!=nil{
		c.Error(err)
		return
	}
}