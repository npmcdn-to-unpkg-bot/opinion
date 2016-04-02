package fakelive

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/palantir/stacktrace"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type Playlist struct {
	Videos             []Video
	StartTime          time.Time
	LiveStreamSettings LiveStreamSettings
}

type playlist struct {
	VideoList []Video `xml:"video"`
}

type Video struct {
	Id int `xml:"id,attr"`
	StartTime int
	EndEndTime int
	Excluded bool
	Id_user         int           `xml:"id_user,attr"`
	Title           string        `xml:"title,attr"`
	Thumbnail       string        `xml:"thumbnail,attr"`
	Poster          string        `xml:"poster,attr"`
	Type            int           `xml:"type,attr"`
	Duration        string        `xml:"duration,attr"`
	DurationSeconds time.Duration `xml:"-"`
	YoutubeId       string        `xml:"-"`
	Is_ad           int           `xml:"is_ad,attr"`
}

func SaveCurrentPlaylist(vids []SmilPlaylist) error {

	return db.Update(func(tx *bolt.Tx) error {
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

func GetCurrentPlaylist()[]SmilPlaylist {
	var playlist []SmilPlaylist
	err := db.View(func(tx *bolt.Tx) error {
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

	return playlist
}
func SaveCurrentSmilPlaylist(smil string) error {

	return db.Update(func(tx *bolt.Tx) error {
		// Retrieve the users bucket.
		// This should be created when the DB is first opened.
		b := tx.Bucket(PlaylistBucket)

		// Generate ID for the user.
		// This returns an error only if the Tx is closed or not writeable.
		// That can't happen in an Update() call so I ignore the error check.
		// Marshal vids data into bytes.
		// Persist bytes to users bucket.
		return b.Put(PlaylistSmilKey, []byte(smil))
	})
}

func GetCurrentSmilPlaylist() string {

	var smil []byte

	err := db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket(PlaylistBucket)

		smil = b.Get(PlaylistSmilKey)

		return nil
	})
	if err != nil {
		log.Println(stacktrace.Propagate(err, ""))
		return ""
	}

	return string(smil)
}

type ByID []Video

func (a ByID) Len() int {
	return len(a)
}
func (a ByID) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a ByID) Less(i, j int) bool {
	return a[i].Id > a[j].Id
}

func firstNAndShuffle(n int, videos []Video) []Video {

	newVideos := videos[:n]
	toShuffle := videos[n:]
	shuffle(toShuffle)
	return append(newVideos, toShuffle...)
}

func shuffle(arr []Video) {
	t := time.Now()
	rand.Seed(int64(t.Nanosecond())) // no shuffling without this line

	for i := len(arr) - 1; i > 0; i-- {
		j := rand.Intn(i)
		arr[i], arr[j] = arr[j], arr[i]
	}
}

func generatePlalist(videos []Video) []Video {

	var selected []Video
	targetduration := 82800 //23h
	sumdur := 0

	//start at the beginning and stop at the end or when we got targetduration
	for i := 0; i < len(videos); i++ {
		if sumdur > targetduration {
			break

		}
		//exclude ads
		if videos[i].Is_ad == 1 || videos[i].Id == 952 {
			continue
		}
		var dur int

		if videoType(videos[i].Type) == embedclip {

			if val, ok := Platform2Youtube[int64(videos[i].Id)]; ok {
				var err error
				dur, err = getYoutubeVideoDuration(val.YtID)
				if err != nil {
					log.Println(stacktrace.Propagate(err, "%s %d", val.YtID, videos[i].Id))
					continue
				}
				videos[i].DurationSeconds = time.Duration(dur) * time.Second
				videos[i].Duration = fmt.Sprint(dur)
				videos[i].YoutubeId = val.YtID
				//do something here
			}

		} else {
			dur = getIntDuration(videos[i].Duration)
			videos[i].Duration = fmt.Sprint(dur)
			videos[i].DurationSeconds = time.Duration(dur) * time.Second
		}
		selected = append(selected, videos[i])
		sumdur += dur

	}

	return selected
}

//returns a playlist based on VIDEOS category
func getPlaylist() ([]Video, error) {

	resp, err := http.Get("http://www.azorestv.com/playlist.php?type=channel&id=1")
	if err != nil {
		return nil, stacktrace.Propagate(err, "")
	}

	buf := bytes.NewBuffer([]byte(""))
	_, err = buf.ReadFrom(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, stacktrace.Propagate(err, "")
	}

	var p playlist
	err = xml.Unmarshal(buf.Bytes(), &p)
	if err != nil {
		return nil, stacktrace.Propagate(err, "")

	}

	if len(p.VideoList) < 10 {
		return nil, stacktrace.Propagate(errors.New("invalid Playlist"), "")
	}

	return p.VideoList, nil
}
