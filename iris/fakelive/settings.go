package fakelive

import (
	"encoding/json"
	"github.com/boltdb/bolt"
	"time"
)

type LiveStreamSettings struct {
	StartTime string
	EndTime   string
	Activated bool
}

type FakeliveSettings struct{
	LiveStreamSettings LiveStreamSettings
	StartTime string
	RepeatTimes []RepeatTimes
	RTimes []time.Time
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

func SetLiveStreamSettings(lss FakeliveSettings) error {

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
		if byteRes==nil{
			return nil
		}


		return json.Unmarshal(byteRes, &res)
	})
	return
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
		if byteRes==nil{
			return nil
		}

		return json.Unmarshal(byteRes, &res)
	})
	return
}
