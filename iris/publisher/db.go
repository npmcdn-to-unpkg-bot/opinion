package publisher

import (
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/asdine/storm"
)


var stormdb *storm.DB

var (
	PublishersBucket = []byte("Publisher")
	SessionsBucket = []byte("Sessions")
	Sessions         *bolt.Bucket
)


