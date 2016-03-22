package sstream

import (
	"time"

	"github.com/boltdb/bolt"
	"github.com/palantir/stacktrace"
	"encoding/json"
	"log"
	"labix.org/v2/mgo/bson"
)

type Settings struct {
	ProtectedStreamName string
}

func newID() string {
	return bson.NewObjectId().Hex()
}

type Client struct {
	Id      string
	Name    string
	Email   string
	Created time.Time
}

func (c *Client)Save() error {
	c.Id = newID()
	return boltdb.Update(func(tx *bolt.Tx) error {
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
	return boltdb.Update(func(tx *bolt.Tx) error {
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
	return boltdb.Update(func(tx *bolt.Tx) error {
		// Retrieve the users bucket.
		// This should be created when the DB is first opened.
		b := tx.Bucket(ClientsBucket)
		// Persist bytes to users bucket.
		return b.Delete([]byte(id))
	})
}

func (c *Client)Get(id string) *Client {
	err := boltdb.View(func(tx *bolt.Tx) error {
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

func (c *Client)GetAll() ([]Client,error ){
	var clients []Client
	err := boltdb.View(func(tx *bolt.Tx) error {
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
		return nil,err
	}

	return clients,nil
}

type Token struct {
	Id         string
	ClientId   string
	ClientName string
	StreamName string
	Expire     time.Time
}

func (t *Token)Save() error {
	t.Id = newID()
	t.ClientName=new(Client).Get(t.ClientId).Name
	return boltdb.Update(func(tx *bolt.Tx) error {
		// Retrieve the users bucket.
		// This should be created when the DB is first opened.
		b := tx.Bucket(TokenBucket)
		// Persist bytes to users bucket.
		out, err := json.Marshal(t)
		if err != nil {
			return err
		}
		return b.Put([]byte(t.Id), out)
	})

	return nil
}


func (c *Token)GetAll() ([]Token,error ){
	var tokens []Token
	err := boltdb.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket(TokenBucket)
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var token Token
			err := json.Unmarshal(v, &token)
			if err != nil {
				return err
			}
			tokens = append(tokens, token)
		}

		return nil
	})
	if err != nil {
		return nil,err
	}

	return tokens,nil
}

func (t *Token)Update() error {
	t.ClientName=new(Client).Get(t.ClientId).Name
	return boltdb.Update(func(tx *bolt.Tx) error {
		// Retrieve the users bucket.
		// This should be created when the DB is first opened.
		b := tx.Bucket(TokenBucket)
		// Persist bytes to users bucket.
		out, err := json.Marshal(t)
		if err != nil {
			return err
		}
		return b.Put([]byte(t.Id), out)
	})
}

func (t *Token)Delete(id string) error {
	return boltdb.Update(func(tx *bolt.Tx) error {
		// Retrieve the users bucket.
		// This should be created when the DB is first opened.
		b := tx.Bucket(TokenBucket)
		// Persist bytes to users bucket.
		return b.Delete([]byte(id))
	})
}

func (t *Token)Get(id string) *Token {
	err := boltdb.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket(TokenBucket)

		out := b.Get([]byte(id))
		if out == nil {
			return nil
		}

		if out != nil {
			err := json.Unmarshal(out, t)
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

	return t
}

func (t *Token)GetClient() *Client {
	if t == nil {
		return nil
	}
	return new(Client).Get(t.ClientId)
}




