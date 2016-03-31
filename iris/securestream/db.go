package securestream

import (
	"github.com/palantir/stacktrace"

	"fmt"
	"github.com/boltdb/bolt"
)

var (
	db             *bolt.DB
	ClientsBucket  = []byte("Clients")
	SettingsBucket = []byte("Settings")
	TokenBucket    = []byte("Token")
)

func createBoltBuckets() error {

	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(ClientsBucket)
		if err != nil {
			return fmt.Errorf("create bucket: %s", err, ClientsBucket)
		}
		return nil
	})
	if err != nil {
		return stacktrace.Propagate(err, "")
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(SettingsBucket)
		if err != nil {
			return fmt.Errorf("create bucket: %s", err, SettingsBucket)
		}

		return nil
	})
	if err != nil {
		return (stacktrace.Propagate(err, ""))
	}

	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(TokenBucket)
		if err != nil {
			return fmt.Errorf("create bucket: %s", err, TokenBucket)
		}

		return nil
	})
}
