package article

import (
	"encoding/base64"

	"github.com/disintegration/imaging"
	"log"

	"time"

	"bytes"

	"github.com/kataras/iris"
	"github.com/thesyncim/opinion/iris/publisher"
	"image"
	"image/jpeg"
	"labix.org/v2/mgo/bson"

	"strings"
)

type Article struct {
	Id            string `storm:"id"`
	Title         string  `storm:"unique"`
	Description   string
	Text          string
	Date          time.Time
	Updated       time.Time
	Approved      bool   `storm:"index"`
	Image         *Base64Img  `storm:"inline"`
	Publisherid   string
	PublisherName string
}
type Base64Img struct {
	Filesize int
	Filetype string
	Filename string
	Base64   string
}

type ArticlesController struct {
}

func resizeImage(str *string) (*string, error) {

	//decode from base64
	read := base64.NewDecoder(base64.StdEncoding, strings.NewReader(*str))
	m, _, err := image.Decode(read)
	if err != nil {
		return nil, err
	}

	//resize Image
	result := imaging.Resize(m, 600, 0, imaging.Lanczos)

	out := bytes.NewBuffer([]byte(""))

	err = jpeg.Encode(out, result, nil)
	if err != nil {
		return nil, err
	}

	//base64 again
	base64resized := base64.StdEncoding.EncodeToString(out.Bytes())

	return &base64resized, nil
}

func (ArticlesController) Create(c *iris.Context) {
	var a Article
	a.Id = bson.NewObjectId().Hex()
	err := c.ReadJSON(&a)
	if err != nil {
		c.Write(err.Error())
		log.Println(err)
		return
	}

	val := c.Get(publisher.IrisContextField)
	if val != nil {
		a.Publisherid = val.(publisher.Session).UserID
	}

	var pub publisher.Publisher

	err = stormdb.One("Id", a.Publisherid, &pub)
	if err != nil {
		log.Println(err)
		c.RenderJSON(500, err.Error())
		return
	}

	a.PublisherName = pub.Name

	a.Date = time.Now()
	a.Updated = time.Now()

	if a.Image != nil {
		tmp, err := resizeImage(&a.Image.Base64)
		if err != nil {
			log.Println(err)
			c.Write(err.Error())
			return
		}
		a.Image.Base64 = *tmp
	}

	err = stormdb.Save(a)

	if err != nil {
		c.RenderJSON(500, err.Error())
		return
	}
}

func (ArticlesController) Edit(c *iris.Context) {
	log.Println("edit")

	var a Article
	err := c.ReadJSON(&a)
	if err != nil {
		log.Println(err)
		c.RenderJSON(500, err.Error())
		return
	}
	var old Article

	err = stormdb.One("Id", a.Id, &old)

	if a.Image == nil {
		a.Image = old.Image
	} else {
		tmp, err := resizeImage(&a.Image.Base64)
		if err != nil {
			log.Println(err)
			c.RenderJSON(500, err.Error())
			return
		}

		a.Image.Base64 = *tmp

	}

	a.Updated = time.Now()

	log.Println("ai")
	err = stormdb.Save(a)
	if err != nil {
		log.Println(err)
		c.RenderJSON(500, err.Error())
		return
	}

}

func (ArticlesController) Delete(c *iris.Context) {

	id := c.Param("id")

	a := Article{Id:id}

	err := stormdb.Remove(a)
	if err != nil {
		c.RenderJSON(500, err)
		return
	}
}

func (ArticlesController) GetId(c *iris.Context) {

	id := c.Param("id")

	var a Article

	err := stormdb.One("Id", id, &a)
	if err != nil {
		c.RenderJSON(500, err)
		return
	}

	c.JSON(a)

}

func (ArticlesController) GetPublisher(c *iris.Context) {

	id := c.Param("id")

	var a = &Article{Id:id}

	err := stormdb.One("Id", id, &a)
	if err != nil {
		c.RenderJSON(500, err)
		return
	}
	var p publisher.Publisher
	err = stormdb.One("Id", id, &p)
	if err != nil {
		c.RenderJSON(500, err)
		return
	}

	c.JSON(p)

}
func prepareArticlesforUser(userID string, articles []Article) ([]Article, error) {
	log.Println("enter")
	var pub publisher.Publisher

	err := stormdb.One("Id", userID, &pub)
	if err != nil {
		log.Println("eerror")
		return nil, err
	}

	log.Println("ai")

	log.Println(pub)




	//if is Admin show all
	if pub.Admin {
		return articles, nil
	}

	var filtered []Article
	for i := range articles {
		if articles[i].Publisherid == userID {
			filtered = append(filtered, articles[i])
		}

	}

	return filtered, nil

}

func (ArticlesController) ListAll(c *iris.Context) {
	var articles []Article
	var userid string
	if val := c.Get(publisher.IrisContextField); val != nil {
		userid = val.(publisher.Session).UserID

	}
	log.Println("fodasse")

	err := stormdb.All(&articles)
	if err != nil {
		log.Println("fodasse all")
		c.RenderJSON(400, err)
		log.Println(err)
		return
	}

	log.Println("ai o caralho", articles)
	result, err := prepareArticlesforUser(userid, articles)
	log.Println("adeus")
	if err != nil {
		log.Println("fodasse prepare")
		c.RenderJSON(400, err)
		log.Println(err)
		return
	}

	log.Println("what??")

	c.JSON(result)
}

func (ArticlesController) ListFrontend(c *iris.Context) {
	var articles []Article

	err := stormdb.Find("Approved", true, &articles)
	if err != nil {
		c.RenderJSON(500, err.Error())
		return
	}

	c.JSON(articles)
}
