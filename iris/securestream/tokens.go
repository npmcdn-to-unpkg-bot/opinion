package securestream

import (

	"github.com/kataras/iris"

	"gopkg.in/mgo.v2/bson"

	"time"
)
type Token struct {
	Id         string `storm:"id"`
	ClientId   string
	ClientName string
	StreamName string
	Expire     time.Time
}


type TokenController struct{}

func (tc *TokenController) Create(c *iris.Context) {
	var token Token
	err := c.ReadJSON(&token)
	if err != nil {
		c.RenderJSON(500,err.Error())
		return
	}
	token.Id=bson.NewObjectId().Hex()

	var client Client

	err=stormdb.One("Id",token.ClientId,&client)
	if err != nil {
		c.RenderJSON(500,err.Error())
		return
	}

	token.ClientName=client.Name

	err=stormdb.Save(token)
	if err != nil {
		c.RenderJSON(500,err.Error())
		return
	}
}

func (tc *TokenController) Read(c *iris.Context) {
	id := c.Param("id")
	token := Token{Id:id}

	err:=stormdb.One("Id",id,&token)
	if err != nil {
		c.RenderJSON(500,err.Error())
		return
	}

	c.JSON(token)
}

func (tc *TokenController) ReadAll(c *iris.Context) {
	var tokens []Token

	err:=stormdb.All(&tokens)
	if err != nil {
		c.RenderJSON(500,err.Error())
		return
	}
	c.JSON(tokens)
}

func (tc *TokenController) Update(c *iris.Context) {
	var token Token

	err := c.ReadJSON(&token)
	if err != nil {
		c.Write(err.Error())
		return
	}

	var client Client

	err=stormdb.One("Id",token.ClientId,&client)
	if err != nil {
		c.RenderJSON(500,err.Error())
		return
	}

	token.ClientName=client.Name

	err=stormdb.Save(token)
	if err != nil {
		c.RenderJSON(500,err.Error())
		return
	}
}

func (tc *TokenController) Delete(c *iris.Context) {
	id := c.Param("id")
	token := &Token{Id:id}

	err:=stormdb.Remove(token)
	if err != nil {
		c.RenderJSON(500,err.Error())
		return
	}


}

