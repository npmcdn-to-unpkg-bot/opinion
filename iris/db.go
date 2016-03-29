package main

import (
	"log"

	"github.com/boltdb/bolt"
)

var db *bolt.DB



func init() {
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	var err error
	db, err = bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
}
