package fakelive

import (
	"fmt"
	"log"

	"github.com/boltdb/bolt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/palantir/stacktrace"
	"github.com/asdine/storm"
)

var (
	sqldb               *gorm.DB
	db                  *bolt.DB
	stormdb *storm.DB
	PlaylistBucket = []byte("Playlist")
	PlaylistKey = []byte("playlist")
	PlaylistSmilKey = []byte("playlistSmil")
	StartTimeKey = []byte("startTime")
	LiveSettingsKey = []byte("LiveSettings")
	FakeLiveSettingsKey = []byte("FakeliveSettings")
	YTVideosBucket = []byte("Youtube")
)

func createBoltBuckets() error {

	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(PlaylistBucket)
		if err != nil {
			return fmt.Errorf("create bucket: %s  %s", err, PlaylistBucket)
		}

		return nil
	})
	if err != nil {
		return stacktrace.Propagate(err, "")
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(YTVideosBucket)
		if err != nil {
			return fmt.Errorf("create bucket: %s  %s", err, YTVideosBucket)
		}

		return nil
	})
	if err != nil {
		return (stacktrace.Propagate(err, ""))
	}

	return nil
}

func init() {
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	var err error

	sqldb, err = gorm.Open(
		"mysql",
		"thesyncim:Kirk1zodiak@tcp(azorestv.com:3306)/azorestv?charset=utf8&parseTime=True",
	)

	if err != nil {
		log.Fatal(stacktrace.Propagate(err, "error connect mysql server"))
	}

}
