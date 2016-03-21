package sstream

import (
	"github.com/palantir/stacktrace"

	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

var (

	boltdb          *bolt.DB
	ClientsBucket  = []byte("Clients")
	SettingsBucket  = []byte("Settings")
	TokenBucket  = []byte("Token")

)

func createBoltBuckets() error {

	err := boltdb.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(ClientsBucket)
		if err != nil {
			return fmt.Errorf("create bucket: %s", err, ClientsBucket)
		}
		return nil
	})
	if err != nil {
		return stacktrace.Propagate(err, "")
	}

	err = boltdb.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(SettingsBucket)
		if err != nil {
			return fmt.Errorf("create bucket: %s", err, SettingsBucket)
		}

		return nil
	})
	if err != nil {
		return (stacktrace.Propagate(err, ""))
	}

	return boltdb.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(TokenBucket)
		if err != nil {
			return fmt.Errorf("create bucket: %s", err, TokenBucket)
		}

		return nil
	})
}

func init() {
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	var err error
	boltdb, err = bolt.Open("sstream.db", 0600, nil)
	if err != nil {
		log.Fatalln(fmt.Errorf("error opening bolt DB", err))
	}

	err = createBoltBuckets()
	if err != nil {
		log.Fatalln(fmt.Errorf("error creating bolt bucket(s)", err))
	}
}