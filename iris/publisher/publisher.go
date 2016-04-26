package publisher

import (
	"time"

	"github.com/kataras/iris"

	"gopkg.in/mgo.v2/bson"
	"log"
	"encoding/base64"
)

type Publisher struct {
	Id       string `storm:"id"`
	Email    string  `storm:"unique"`
	Password string
	Salt     string

	Name     string
	Image    *Base64Img `storm:"inline"`
	Admin    bool
	Date     time.Time
	Updated  time.Time
}

type Base64Img struct {
	Filesize int
	Filetype string
	Filename string
	Base64   string
}

func (pub *Publisher) ID() string {

	return pub.Id
}

func (pub *Publisher) PASSWORD() string {

	return pub.Password
}

func (pub *Publisher) FindUser(email string) (User, error) {

	err := stormdb.One("Email", email, pub)
	if err != nil {
		return nil, err
	}

	return pub, nil

}

type PublisherController struct {
}

func (PublisherController) Create(c *iris.Context) {
	var p = &Publisher{}
	err := c.ReadJSON(p)
	if err != nil {
		c.Write(err.Error())
		return
	}

	p.Id = bson.NewObjectId().Hex()

	p.Password = NewSha512Password(p.Password)

	err = stormdb.Save(p)
	if err != nil {
		c.RenderJSON(500, err.Error())
		return
	}
}

func (PublisherController) Edit(c *iris.Context) {

	id := c.Param("id")
	var p Publisher
	err := c.ReadJSON(&p)
	if err != nil {
		if err != nil {
			c.RenderJSON(500, err.Error())
			return
		}
	}

	var old Publisher

	err = stormdb.One("Id", id, &old)
	if err != nil {
		c.RenderJSON(500, err.Error())
		return
	}

	p.Updated = time.Now()
	if p.Password != old.Password {
		p.Password = NewSha512Password(p.Password)
	}

	err = stormdb.Save(p)
	if err != nil {
		c.RenderJSON(500, err)
		return
	}

}

func (PublisherController) GetId(c *iris.Context) {

	id := c.Param("id")
	var pub Publisher

	err := stormdb.One("Id", id, &pub)
	if err != nil {
		c.RenderJSON(500, err.Error())
		return
	}

	c.JSON(pub)
}

func (PublisherController) GetImage(c *iris.Context) {

	id := c.Param("id")

	var pub Publisher
	err := stormdb.One("Id", id, &pub)
	if err != nil {
		c.RenderJSON(500, err.Error())
		return
	}


	buf,err:=base64.StdEncoding.DecodeString(pub.Image.Base64)
	if err != nil {
		c.RenderJSON(500, err.Error())
		return
	}

	c.Response.Header.Set("Content-Type",pub.Image.Filetype)
	c.Write(buf)

}

func (PublisherController) Delete(c *iris.Context) {

	id := c.Param("id")
	p := Publisher{Id:id}

	err := stormdb.Remove(p)
	if err != nil {
		c.RenderJSON(500, err)
		return
	}
}

func (PublisherController) ListAll(c *iris.Context) {
	var publishers []Publisher

	err := stormdb.All(&publishers)
	if err != nil {
		c.RenderJSON(500, err)
		return
	}

	c.JSON(publishers)
}

func AddDefaultPub() error {
	var p = &Publisher{}
	p.Id = bson.NewObjectId().Hex()
	p.Name = "Marcelo Pires"
	p.Password = "Kirk1zodiak"
	p.Email = "thesyncim@gmail.com"
	p.Admin = true

	p.Password = NewSha512Password(p.Password)

	//todo validate existing email

	log.Println(stormdb.Save(p))
	return nil

}

func init() {

	//log.Println(AddDefaultPub())

}
