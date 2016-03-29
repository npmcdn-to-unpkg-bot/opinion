package publisher

import (
	"log"

	"fmt"

	"github.com/boltdb/bolt"
)

var db *bolt.DB

var (
	PublishersBucket = []byte("Publisher")
	SessionsBucket   = []byte("Sessions")
	Sessions         *bolt.Bucket
)

func createBoltBuckets() error {
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.

	return db.Update(func(tx *bolt.Tx) error {

		_, err := tx.CreateBucketIfNotExists(PublishersBucket)
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
