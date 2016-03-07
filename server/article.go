package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
	"labix.org/v2/mgo/bson"
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

func (art *Article) Delete(id string) error {

	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(ArticlesBucket)
		return b.Delete([]byte(id))
	})
}

func (art *Article) Publisher(id string) *Publisher {

	var p = Publisher{}
	db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket(PublishersBucket)

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {

			var pp Publisher
			err := json.Unmarshal(v, &pp)
			if err != nil {
				//todo handle error
				continue
			}
			if string(k) == id {
				p = pp

			}

		}
		return nil

	})

	return &p

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

		b := tx.Bucket(PublishersBucket)

		return b.Put([]byte(a.Id), buf)
	})

}

type ArticlesController struct {
}

func (ArticlesController) Create(c *gin.Context) {
	var a = Article{}
	err := c.BindJSON(&a)
	if err != nil {
		c.Error(err)
		return
	}

	val, ok := c.Get(GinContextField)
	if ok {
		a.Publisherid = val.(Session).UserID

	}

	p := a.Publisher(a.Publisherid)
	a.PublisherName = p.Name

	a.Date = time.Now()
	a.Updated = time.Now()

	buf, err := json.Marshal(a)
	if err != nil {
		c.Error(err)
		return
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(ArticlesBucket)
		err := b.Put([]byte(bson.NewObjectId().Hex()), buf)
		return err
	})
	if err != nil {
		c.Error(err)
	}
}

func (ArticlesController) Edit(c *gin.Context) {

	id := c.Param("id")
	var a Article
	err := c.BindJSON(&a)
	if err != nil {
		c.Error(err)
		return
	}

	a.Updated = time.Now()
	buf, err := json.Marshal(a)
	if err != nil {
		c.Error(err)
		return
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(ArticlesBucket)
		err := b.Put([]byte(id), buf)
		return err
	})
	if err != nil {
		c.Error(err)
	}

}

func (ArticlesController) Delete(c *gin.Context) {

	id := c.Param("id")

	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(ArticlesBucket)
		err := b.Delete([]byte(id))
		return err
	})
	if err != nil {
		c.Error(err)
	}
}

func (ArticlesController) GetId(c *gin.Context) {

	id := c.Param("id")

	var a Article

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(ArticlesBucket)
		bp := b.Get([]byte(id))

		err := json.Unmarshal(bp, &a)
		if err != nil {
			return err
		}

		return err
	})
	if err != nil {
		c.Error(err)
	}
	c.JSON(http.StatusOK, a)

}

func (ArticlesController) GetPublisher(c *gin.Context) {

	id := c.Param("id")

	var a = &Article{}

	p := a.Publisher(id)

	if p == nil {
		c.JSON(500, "invalid publisher")

		return
	}

	c.JSON(http.StatusOK, p)

}
func prepareArticlesforUser(userID string, articles []Article) []Article {

	var p = &Publisher{}

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

func (ArticlesController) ListAll(c *gin.Context) {
	var articles []Article
	var userid string
	if val, ok := c.Get(GinContextField); ok {
		userid = val.(Session).UserID

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
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, prepareArticlesforUser(userid, articles))
}

func (ArticlesController) ListFrontend(c *gin.Context) {
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
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, articles)
}
