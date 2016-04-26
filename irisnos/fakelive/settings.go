package fakelive

import (
	"encoding/json"
	"github.com/boltdb/bolt"
	"log"
	"time"
)

type LiveStreamSettings struct {
	StartTime string
	EndTime   string
	Activated bool
}

type FakeliveSettings struct {
	LiveStreamSettings LiveStreamSettings
	StartTime          string
	RepeatTimes        []RepeatTimes
	RTimes             []time.Time
}

func SetFakeliveSettings(lss FakeliveSettings) error {

	return db.Update(func(tx *bolt.Tx) error {
		// Retrieve the users bucket.
		// This should be created when the DB is first opened.
		b := tx.Bucket(PlaylistBucket)

		buf, err := json.Marshal(lss)
		if err != nil {
			return err
		}
		// Persist bytes to users bucket.
		return b.Put(FakeLiveSettingsKey, buf)
	})

}

func GetFakeliveSettings() (res FakeliveSettings, err error) {

	err = db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket(PlaylistBucket)

		byteRes := b.Get(FakeLiveSettingsKey)
		log.Println("result ", byteRes)
		if byteRes == nil {
			return nil
		}

		return json.Unmarshal(byteRes, &res)
	})
	return
}
