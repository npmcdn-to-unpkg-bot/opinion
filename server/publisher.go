package main

import (
	"encoding/json"
	"net/http"

	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"

	"log"
	"time"

	"github.com/syndtr/goleveldb/leveldb/errors"
	"gopkg.in/mgo.v2/bson"
)

/*
type UserRepository interface {
	Signin(email, pass string) (string, error)
	Login(email, pass string) (string, error)
	UserByEmail(email string) (User, error)
}
User
*/

type Publisher struct {
	Id       string
	Email    string
	Password string
	Salt     string

	Name     string
	Image    *Base64Img
	Admin    bool
	Date     time.Time
}

/*
func (pub *Publisher)Signin(email, pass string) (string, error) {
	_, err := pub.UserByEmail(email)
	if err != nil {
		return "", store.ErrEmailDuplication
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(PublishersBucket)

		// email is going to be the Id of the user
		user, err := store.NewUser(email, email, pass)
		if err != nil {
			return err
		}



		pub.User=user
		buf, err := json.Marshal(pub)
		if err != nil {
			return err
		}

		pretty.Println(pub)

		err = b.Put([]byte(bson.NewObjectId().Hex()), buf)
		return err
	})

	if err != nil {
		return "", err
	}

	return email, nil
}
*/

/*
func (pub *Publisher) UserByEmail(email string) (store.User, error) {

	var user store.User
	err :=db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket(PublishersBucket)
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var p Publisher
			err := json.Unmarshal(v, &p)
			if err != nil {
				return err
			}
			pretty.Println(p)
			if p.Email == email {
				*pub = p
				user = pub.User
			}
		}

		return nil

	})
	if err!=nil{
		return store.User{},store.ErrUserNotFound
	}

	if pub == nil {
		return store.User{},store.ErrUserNotFound
	}

	return user,nil
}*/

/*
func (pub *Publisher) Login(email, pass string) (string, error){

	user, err := pub.UserByEmail(email)
	if err != nil {
		return "", store.ErrWrongPassword
	}


	pretty.Println("pass",user.Password)
	passStored := user.Password
	salt := user.Salt

	hpass, err := crypto.HashPassword(pass, []byte(salt))
	if err != nil {
		return "", err
	}
	passOk := crypto.SecureCompare(hpass, []byte(passStored))
	if !passOk {
		log.Println("wrong")
		return "", store.ErrWrongPassword
	}
	return user.Id, nil
}
*/

func (pub *Publisher) Delete(id string) error {

	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(PublishersBucket)
		return b.Delete([]byte(id))
	})
}

func (pub *Publisher) ID() string {

	return pub.Id
}

func (pub *Publisher) PASSWORD() string {

	return pub.Password
}

func (pub *Publisher) Get(id string) *Publisher {

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(PublishersBucket)
		bp := b.Get([]byte(id))

		err := json.Unmarshal(bp, pub)
		if err != nil {
			return err
		}
		pub.Id = id

		return err
	})
	if err != nil {
		return nil
	}

	return pub
}

func (pub *Publisher) Update(p *Publisher) error {

	return db.Update(func(tx *bolt.Tx) error {

		buf, err := json.Marshal(p)
		if err != nil {

			return err
		}

		b := tx.Bucket(PublishersBucket)

		return b.Put([]byte(p.Id), buf)
	})

}

func (pub *Publisher) GetUser(email string) *Publisher {

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(PublishersBucket)
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {

			var p Publisher
			err := json.Unmarshal(v, &p)
			if err != nil {
				return err
			}
			if p.Email == email {
				*pub = p
			}

		}

		return nil
	})
	if err != nil {
		log.Println(err)
		return nil
	}

	return pub

}

func (pub *Publisher) FindUser(email string) (User, error) {

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(PublishersBucket)
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {

			var p Publisher
			err := json.Unmarshal(v, &p)
			if err != nil {
				return err
			}
			p.Id = string(k)
			if p.Email == email {
				*pub = p
			}

		}

		return nil
	})
	if err != nil {
		log.Println(err)
		return nil, errors.New("Error")
	}

	if pub.Id == "" {
		return nil, errors.New("not found")
	}

	return pub, nil

}

type PublisherController struct {
}

func (PublisherController) Create(c *gin.Context) {
	var p = &Publisher{}
	err := c.BindJSON(p)
	if err != nil {
		c.Error(err)
		return
	}

	p.Password = NewSha512Password(p.Password)

	//todo validate existing email

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(PublishersBucket)

		// email is going to be the Id of the user

		buf, err := json.Marshal(p)
		if err != nil {
			return err
		}

		err = b.Put([]byte(bson.NewObjectId().Hex()), buf)
		return err
	})

	if err != nil {
		c.Error(err)
	}
}

func (PublisherController) Edit(c *gin.Context) {

	id := c.Param("id")
	var p = &Publisher{}
	err := c.BindJSON(p)
	if err != nil {
		c.Error(err)
		return
	}

	old := &Publisher{}
	log.Println(id)
	old.Get(id)

	log.Println(old)

	if p.Image == nil {
		p.Image = old.Image
	}

	if p.Password != old.Password {
		p.Password = NewSha512Password(p.Password)
	}

	p.Id = id

	err = old.Update(p)
	if err != nil {
		log.Println("---------->", err)
		c.Error(err)
	}

}

func (PublisherController) GetId(c *gin.Context) {

	id := c.Param("id")
	var p = &Publisher{}

	p.Get(id)

	c.JSON(http.StatusOK, p)

}

func (PublisherController) Delete(c *gin.Context) {

	id := c.Param("id")

	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(PublishersBucket)
		err := b.Delete([]byte(id))
		return err
	})
	if err != nil {
		c.Error(err)
	}
}

func (PublisherController) ListAll(c *gin.Context) {

	var publishers []Publisher
	err := db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket(PublishersBucket)
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {

			var p Publisher
			err := json.Unmarshal(v, &p)
			if err != nil {
				return err
			}
			p.Id = string(k)

			publishers = append(publishers, p)
		}

		return nil
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, publishers)
}

func AddDefaultPub() error {
	var p = &Publisher{}
	p.Name = "Marcelo Pires"
	p.Password = "Kirk1zodiak"
	p.Email = "thesyncim@gmail.com"
	p.Admin = true

	p.Password = NewSha512Password(p.Password)

	//todo validate existing email

	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(PublishersBucket)

		// email is going to be the Id of the user

		buf, err := json.Marshal(p)
		if err != nil {
			return err
		}

		err = b.Put([]byte(bson.NewObjectId().Hex()), buf)
		return err
	})

}

func init() {

	//log.Println(AddDefaultPub())

}
