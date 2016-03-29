package securestream

import (
	"github.com/kataras/iris"
	"github.com/palantir/stacktrace"
	"github.com/boltdb/bolt"
	"time"
	"encoding/json"
	"log"
	"gopkg.in/mgo.v2/bson"
)

type ClientController struct{}

func (cc *ClientController)Create(c *iris.Context) {
	var client = &Client{}
	err := c.ReadJSON(client)
	if err != nil {
		c.Write(err.Error())
		return
	}

	err = client.Save()
	if err != nil {
		c.Write(err.Error())
		return
	}

}

func (cc *ClientController)Read(c *iris.Context) {
	id := c.Param("id")
	var client = &Client{}

	c.JSON(client.Get(id))
}

func (cc *ClientController)ReadAll(c *iris.Context) {
	var client = &Client{}

	clients, err := client.GetAll()
	if err != nil {
		c.Write(err.Error())
		return
	}
	c.JSON(clients)
}

func (cc *ClientController)Update(c *iris.Context) {
	var client = &Client{}
	err := c.ReadJSON(client)
	if err != nil {
		c.Write(err.Error())
		return
	}

	err = client.Update()
	if err != nil {
		c.Write(err.Error())
		return
	}
}

func (cc *ClientController)Delete(c *iris.Context) {
	id := c.Param("id")
	var client = &Client{}

	err := client.Delete(id)
	if err != nil {
		c.Write(err.Error())
		return
	}
}

type Client struct {
	Id      string
	Name    string
	Email   string
	Created time.Time
}

func (c *Client)Save() error {
	c.Id = bson.NewObjectId().Hex()
	return db.Update(func(tx *bolt.Tx) error {
		// Retrieve the users bucket.
		// This should be created when the DB is first opened.
		b := tx.Bucket(ClientsBucket)
		// Persist bytes to users bucket.
		out, err := json.Marshal(c)
		if err != nil {
			return err
		}
		return b.Put([]byte(c.Id), out)
	})

	return nil
}

func (c *Client)Update() error {
	return db.Update(func(tx *bolt.Tx) error {
		// Retrieve the users bucket.
		// This should be created when the DB is first opened.
		b := tx.Bucket(ClientsBucket)
		// Persist bytes to users bucket.
		out, err := json.Marshal(c)
		if err != nil {
			return err
		}
		return b.Put([]byte(c.Id), out)
	})
}

func (c *Client)Delete(id string) error {
	return db.Update(func(tx *bolt.Tx) error {
		// Retrieve the users bucket.
		// This should be created when the DB is first opened.
		b := tx.Bucket(ClientsBucket)
		// Persist bytes to users bucket.
		return b.Delete([]byte(id))
	})
}

func (c *Client)Get(id string) *Client {
	err := db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket(ClientsBucket)

		out := b.Get([]byte(id))
		if out == nil {
			return nil
		}

		if out != nil {
			err := json.Unmarshal(out, c)
			if err != nil {
				return stacktrace.Propagate(err, "")
			}
		}

		return nil
	})
	if err != nil {
		log.Println(stacktrace.Propagate(err, ""))
		return nil
	}

	return c
}

func (c *Client)GetAll() ([]Client, error) {
	var clients []Client
	err := db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket(ClientsBucket)
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var client Client
			err := json.Unmarshal(v, &client)
			if err != nil {
				return err
			}
			clients = append(clients, client)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return clients, nil
}


