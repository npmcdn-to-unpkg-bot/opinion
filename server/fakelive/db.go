package fakelive

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/boltdb/bolt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/palantir/stacktrace"
)

var (
	db              *gorm.DB
	boltdb          *bolt.DB
	PlaylistBucket  = []byte("Playlist")
	PlaylistKey     = []byte("playlist")
	StartTimeKey    = []byte("startTime")
	LiveSettingsKey = []byte("LiveSettings")
	YTVideosBucket  = []byte("Youtube")
)

func createBoltBuckets() error {

	err := boltdb.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(PlaylistBucket)
		if err != nil {
			return fmt.Errorf("create bucket: %s", err, PlaylistBucket)
		}

		return nil
	})
	if err != nil {
		return stacktrace.Propagate(err, "")
	}

	err = boltdb.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(YTVideosBucket)
		if err != nil {
			return fmt.Errorf("create bucket: %s", err, YTVideosBucket)
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
	boltdb, err = bolt.Open("playlist.db", 0600, nil)
	if err != nil {
		log.Fatalln(fmt.Errorf("error opening bolt DB", err))
	}

	err = createBoltBuckets()
	if err != nil {
		log.Fatalln(fmt.Errorf("error creating bolt bucket(s)", err))
	}

	db, err = gorm.Open(
		"mysql",
		"thesyncim:Kirk1zodiak@tcp(azorestv.com:3306)/azorestv?charset=utf8&parseTime=True",
	)

	if err != nil {
		log.Fatal(stacktrace.Propagate(err, "error connect mysql server"))
	}

}

func SaveCurrentPlaylist(vids *Playlist) error {

	return boltdb.Update(func(tx *bolt.Tx) error {
		// Retrieve the users bucket.
		// This should be created when the DB is first opened.
		b := tx.Bucket(PlaylistBucket)

		// Generate ID for the user.
		// This returns an error only if the Tx is closed or not writeable.
		// That can't happen in an Update() call so I ignore the error check.

		// Marshal vids data into bytes.
		buf, err := json.Marshal(vids)
		if err != nil {
			return stacktrace.Propagate(err, "")
		}

		// Persist bytes to users bucket.
		return b.Put(PlaylistKey, buf)
	})
}

func GetCurrentPlaylist() *Playlist {
	var playlist Playlist
	err := boltdb.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket(PlaylistBucket)

		out := b.Get(PlaylistKey)
		if out == nil {
			return nil
		}

		err := json.Unmarshal(out, &playlist)
		if err != nil {
			return stacktrace.Propagate(err, "")
		}

		return nil
	})
	if err != nil {
		log.Println(stacktrace.Propagate(err, ""))
		return nil
	}

	return &playlist
}

func SetStartTime(st string) error {

	return boltdb.Update(func(tx *bolt.Tx) error {
		// Retrieve the users bucket.
		// This should be created when the DB is first opened.
		b := tx.Bucket(PlaylistBucket)
		// Persist bytes to users bucket.
		return b.Put(StartTimeKey, []byte(st))
	})

}
func getStartTime() (res string, err error) {

	err = boltdb.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket(PlaylistBucket)

		byteRes := b.Get(StartTimeKey)
		res = string(byteRes)

		return nil
	})
	return
}

func SetLiveStreamSettings(lss LiveStreamSettings) error {

	return boltdb.Update(func(tx *bolt.Tx) error {
		// Retrieve the users bucket.
		// This should be created when the DB is first opened.
		b := tx.Bucket(PlaylistBucket)

		buf, err := json.Marshal(lss)
		if err != nil {
			return err
		}
		// Persist bytes to users bucket.
		return b.Put(LiveSettingsKey, buf)
	})

}
func GetLiveStreamSettings() (res LiveStreamSettings, err error) {

	err = boltdb.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket(PlaylistBucket)

		byteRes := b.Get(LiveSettingsKey)

		return json.Unmarshal(byteRes, &res)
	})
	return
}
