package main

import (
	"log"

	"fmt"

	"github.com/boltdb/bolt"
)

var db *bolt.DB

var (
	ArticlesBucket = []byte("Articles")
	PublishersBucket = []byte("Publisher")
	SessionsBucket = []byte("Sessions")
	Sessions         *bolt.Bucket
)

func init() {
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	var err error
	db, err = bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(ArticlesBucket)
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		_, err = tx.CreateBucketIfNotExists(PublishersBucket)
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		Sessions, err = tx.CreateBucketIfNotExists(SessionsBucket)
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

}
