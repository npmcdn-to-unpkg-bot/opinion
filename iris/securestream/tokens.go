package securestream

import (
	"encoding/json"
	"github.com/boltdb/bolt"
	"github.com/kataras/iris"
	"github.com/palantir/stacktrace"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

type TokenController struct{}

func (tc *TokenController) Create(c *iris.Context) {
	var Token = &Token{}
	err := c.ReadJSON(Token)
	if err != nil {
		c.Write(err.Error())
		return
	}

	err = Token.Save()
	if err != nil {
		c.Write(err.Error())
		return
	}
}

func (tc *TokenController) Read(c *iris.Context) {
	id := c.Param("id")
	var Token = &Token{}

	c.JSON(Token.Get(id))
}

func (tc *TokenController) ReadAll(c *iris.Context) {
	var Token = &Token{}

	Tokens, err := Token.GetAll()
	if err != nil {
		c.Write(err.Error())
		return
	}
	c.JSON(Tokens)
}

func (tc *TokenController) Update(c *iris.Context) {
	var Token = &Token{}

	err := c.ReadJSON(Token)
	if err != nil {
		c.Write(err.Error())
		return
	}

	err = Token.Update()
	if err != nil {
		c.Write(err.Error())
		return
	}
}

func (tc *TokenController) Delete(c *iris.Context) {
	id := c.Param("id")
	var Token = &Token{}

	err := Token.Delete(id)
	if err != nil {
		c.Write(err.Error())
		return
	}
}

type Token struct {
	Id         string
	ClientId   string
	ClientName string
	StreamName string
	Expire     time.Time
}

func (t *Token) Save() error {
	t.Id = bson.NewObjectId().Hex()
	t.ClientName = new(Client).Get(t.ClientId).Name
	return db.Update(func(tx *bolt.Tx) error {
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

func (c *Token) GetAll() ([]Token, error) {
	var tokens []Token
	err := db.View(func(tx *bolt.Tx) error {
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
		return nil, err
	}

	return tokens, nil
}

func (t *Token) Update() error {
	t.ClientName = new(Client).Get(t.ClientId).Name
	return db.Update(func(tx *bolt.Tx) error {
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

func (t *Token) Delete(id string) error {
	return db.Update(func(tx *bolt.Tx) error {
		// Retrieve the users bucket.
		// This should be created when the DB is first opened.
		b := tx.Bucket(TokenBucket)
		// Persist bytes to users bucket.
		return b.Delete([]byte(id))
	})
}

func (t *Token) Get(id string) *Token {
	err := db.View(func(tx *bolt.Tx) error {
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

func (t *Token) GetClient() *Client {
	if t == nil {
		return nil
	}
	return new(Client).Get(t.ClientId)
}
