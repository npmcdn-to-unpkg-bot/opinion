package securestream

import (
	"github.com/kataras/iris"
	"gopkg.in/mgo.v2/bson"

	"time"
)

type Client struct {
	Id      string `storm:"id"`
	Name    string
	Email   string
	Created time.Time
}

type ClientController struct{}

func (cc *ClientController) Create(c *iris.Context) {
	var client Client
	err := c.ReadJSON(&client)
	if err != nil {
		c.RenderJSON(500, err.Error())
		return
	}

	client.Id = bson.NewObjectId().Hex()

	err = stormdb.Save(client)
	if err != nil {
		c.RenderJSON(500, err.Error())
		return
	}

}

func (cc *ClientController) Read(c *iris.Context) {
	id := c.Param("id")
	var client Client

	err := stormdb.One("Id", id, &client)
	if err != nil {
		c.RenderJSON(500, err.Error())
		return
	}

	c.JSON(client)
}

func (cc *ClientController) ReadAll(c *iris.Context) {
	var clients  []Client

	err := stormdb.All(&clients)
	if err != nil {
		c.RenderJSON(500, err.Error())
		return
	}

	c.JSON(clients)
}

func (cc *ClientController) Update(c *iris.Context) {
	var client Client
	err := c.ReadJSON(&client)
	if err != nil {
		c.RenderJSON(500, err.Error())
		return
	}

	err = stormdb.Save(client)
	if err != nil {
		c.RenderJSON(500, err.Error())
		return
	}
}

func (cc *ClientController) Delete(c *iris.Context) {
	id := c.Param("id")

	client := &Client{Id:id}

	err := stormdb.Remove(client)
	if err != nil {
		c.RenderJSON(500, err.Error())
		return
	}
}



