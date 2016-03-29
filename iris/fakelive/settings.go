package fakelive

import (
	"encoding/json"
	"github.com/boltdb/bolt"
)

type LiveStreamSettings struct {
	StartTime string
	EndTime   string
	Activated bool
}

func SetStartTime(st string) error {

	return db.Update(func(tx *bolt.Tx) error {
		// Retrieve the users bucket.
		// This should be created when the DB is first opened.
		b := tx.Bucket(PlaylistBucket)
		// Persist bytes to users bucket.
		return b.Put(StartTimeKey, []byte(st))
	})

}
func getStartTime() (res string, err error) {

	err = db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket(PlaylistBucket)

		byteRes := b.Get(StartTimeKey)
		res = string(byteRes)

		return nil
	})
	return
}

func SetLiveStreamSettings(lss LiveStreamSettings) error {

	return db.Update(func(tx *bolt.Tx) error {
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

	err = db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket(PlaylistBucket)

		byteRes := b.Get(LiveSettingsKey)

		return json.Unmarshal(byteRes, &res)
	})
	return
}
