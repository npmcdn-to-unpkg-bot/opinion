package article

import (


	"fmt"

	"github.com/boltdb/bolt"
)

var db *bolt.DB

var (
	ArticlesBucket   = []byte("Articles")
)

func createBoltBuckets()error {
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(ArticlesBucket)
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}


		return nil
	})

}
