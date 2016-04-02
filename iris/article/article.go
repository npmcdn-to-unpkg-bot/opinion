package article

import (
	"encoding/base64"
	"encoding/json"
	"github.com/disintegration/imaging"
	"log"

	"time"

	"bytes"
	"github.com/boltdb/bolt"

	"github.com/kataras/iris"
	"github.com/thesyncim/opinion/iris/publisher"
	"image"
	"image/jpeg"
	"labix.org/v2/mgo/bson"
	"net/http"
	"strings"
)

type Article struct {
	Id            string
	Title         string
	Description   string
	Text          string
	Date          time.Time
	Updated       time.Time
	Approved      bool
	Image         *Base64Img
	Publisherid   string
	PublisherName string
}
type Base64Img struct {
	Filesize int
	Filetype string
	Filename string
	Base64   string
}

func (art *Article) Delete(id string) error {

	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(ArticlesBucket)
		return b.Delete([]byte(id))
	})
}

func (art *Article) Publisher(id string) *publisher.Publisher {
	p := &publisher.Publisher{}
	return p.Get(id)
}

func (art *Article) Get(id string) *Article {

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(ArticlesBucket)
		bp := b.Get([]byte(id))

		err := json.Unmarshal(bp, art)
		if err != nil {
			return err
		}

		return err
	})
	if err != nil {
		return nil
	}

	return art
}

func (art *Article) Update(a *Article) error {

	return db.Update(func(tx *bolt.Tx) error {

		buf, err := json.Marshal(a)
		if err != nil {

			return err
		}
		b := tx.Bucket(ArticlesBucket)
		return b.Put([]byte(a.Id), buf)
	})

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
	var a = Article{}
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

	p := a.Publisher(a.Publisherid)

	log.Println(&a == nil, p == nil)
	a.PublisherName = p.Name

	a.Date = time.Now()
	a.Updated = time.Now()

	if a.Image != nil {

		tmp, err := resizeImage(&a.Image.Base64)
		if err != nil {
			c.Write(err.Error())
			return
		}

		a.Image.Base64 = *tmp

	}

	buf, err := json.Marshal(a)
	if err != nil {
		c.Write(err.Error())
		return
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(ArticlesBucket)
		err := b.Put([]byte(bson.NewObjectId().Hex()), buf)
		return err
	})
	if err != nil {
		c.Write(err.Error())
	}
}

func (ArticlesController) Edit(c *iris.Context) {

	id := c.Param("id")
	var a Article
	err := c.ReadJSON(&a)
	if err != nil {
		c.Write(err.Error())
		return
	}

	a.Updated = time.Now()
	buf, err := json.Marshal(a)
	if err != nil {
		c.Write(err.Error())
		return
	}
	tmp, err := resizeImage(&a.Image.Base64)
	if err != nil {
		c.Write(err.Error())
		return
	}

	a.Image.Base64 = *tmp

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(ArticlesBucket)
		err := b.Put([]byte(id), buf)
		return err
	})
	if err != nil {
		c.Write(err.Error())
	}

}

func (ArticlesController) Delete(c *iris.Context) {

	id := c.Param("id")

	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(ArticlesBucket)
		err := b.Delete([]byte(id))
		return err
	})
	if err != nil {
		c.Write(err.Error())
	}
}

func (ArticlesController) GetId(c *iris.Context) {

	id := c.Param("id")

	var a Article

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(ArticlesBucket)
		bp := b.Get([]byte(id))

		err := json.Unmarshal(bp, &a)
		if err != nil {
			return err
		}

		a.Id = id

		return err
	})
	if err != nil {
		c.Write(err.Error())
	}
	c.JSON(a)

}

func (ArticlesController) GetPublisher(c *iris.Context) {

	id := c.Param("id")

	var a = &Article{}

	p := a.Publisher(id)

	if p == nil {
		c.RenderJSON(http.StatusInternalServerError, "invalid publisher")

		return
	}

	c.JSON(p)

}
func prepareArticlesforUser(userID string, articles []Article) []Article {

	var p = &publisher.Publisher{}

	pub := p.Get(userID)
	if p == nil {
		return nil
	}

	//if is Admin show all
	if pub.Admin {
		return articles
	}

	var filtered []Article
	for i := range articles {
		if articles[i].Publisherid == userID {
			filtered = append(filtered, articles[i])
		}

	}

	return filtered

}

func (ArticlesController) ListAll(c *iris.Context) {
	var articles []Article
	var userid string
	if val := c.Get(publisher.IrisContextField); val != nil {
		userid = val.(publisher.Session).UserID

	}

	err := db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket(ArticlesBucket)

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {

			var a Article
			err := json.Unmarshal(v, &a)
			if err != nil {
				return err
			}
			a.Id = string(k)

			log.Println(a)
			articles = append(articles, a)
		}

		return nil
	})

	if err != nil {
		c.Write(err.Error())
		return
	}

	if articles == nil || userid == "" {
		c.JSON(nil)
		return
	}

	c.JSON(prepareArticlesforUser(userid, articles))
}

func (ArticlesController) ListFrontend(c *iris.Context) {
	var articles []Article

	err := db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket(ArticlesBucket)

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {

			var a Article
			err := json.Unmarshal(v, &a)
			if err != nil {
				return err
			}
			if !a.Approved {
				continue
			}
			a.Id = string(k)

			log.Println(a)
			articles = append(articles, a)
		}

		return nil
	})

	if err != nil {
		c.Write(err.Error())
		return
	}

	c.JSON(articles)
}
