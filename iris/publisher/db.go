package publisher

import (
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/asdine/storm"
)

var db *bolt.DB
var stormdb *storm.DB

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
