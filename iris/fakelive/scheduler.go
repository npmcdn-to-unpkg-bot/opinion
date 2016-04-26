package fakelive

import (
	"log"
	"os"
	"os/exec"

	"fmt"
	"github.com/carlescere/scheduler"
	"github.com/jinzhu/now"
	"github.com/kataras/iris"
	"github.com/palantir/stacktrace"
	"sort"
	"sync"
	"time"
	"github.com/boltdb/bolt"

	"strings"
	"github.com/Azure/azure-sdk-for-go/core/http"
	"bytes"
)

var basedir = "/var/www/vhosts/azorestv.com/httpdocs/uploads/movies/yt"
var j *scheduler.Job

func RunBackgroundScheduler() *scheduler.Job {
	work()
	fls, _ := GetFakeliveSettings()
	if fls.StartTime == "" {
		fls.StartTime = "00:00"
	}
	var err error
	j, err = scheduler.Every().Day().At(fls.StartTime).Run(func() {
		work()
	})

	if err != nil {
		log.Fatalln(stacktrace.Propagate(err, ""))

	}
	return j
}

func work() {

	err := downloadMissingYoutubeVideos()
	if err != nil {
		log.Println(err)
	}

	allvids, err := getPlaylist()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(syncPlaylist(allvids))

	syncPlaylist(allvids)

	videos := generatePlalist(allvids)

	sort.Sort(ByID(videos))

	videos = firstNAndShuffle(3, videos)

	playlissmill, videos := genSmilPlaylistSlice(videos, calcScheduleDate())
	smil := genSmil(playlissmill)

	SaveCurrentSmilPlaylist(smil)

	s, err := os.Create("/var/www/vhosts/azorestv.com/httpdocs/uploads/movies/streamschedule.smil")
	if err != nil {
		log.Fatalln(stacktrace.Propagate(err, ""))
	}

	checkWriteErr(s.WriteString(smil))
	defer checkErr(s.Close())

	//const longForm = "2006-01-02 15:04:05"
	//t, _ := time.Parse(longForm, calcScheduleDate())

	err = SaveCurrentPlaylist(playlissmill)
	if err != nil {
		log.Fatalln(stacktrace.Propagate(err, ""))
	}

	cmd := exec.Command("service", "WowzaStreamingEngine", "restart")
	err = cmd.Start()
	if err != nil {
		log.Fatalln(stacktrace.Propagate(err, ""))
	}
	cmd.Wait()

}

type FakeliveController struct{}

func (*FakeliveController) CurrentPlaylist(c *iris.Context) {

	c.JSON(GetCurrentPlaylist())

}

func (*FakeliveController) CurrentSmilPlaylist(c *iris.Context) {

	c.WriteText(200, GetCurrentSmilPlaylist())
}

type RepeatTimes struct {
	At   time.Time
	Once sync.Once
}

func (r RepeatTimes) getHourSeconds() time.Time {
	hours := fmt.Sprintf("%02d", r.At.Hour())
	seconds := fmt.Sprintf("%02d", r.At.Minute())
	log.Println(now.MustParse(hours + ":" + seconds))
	return now.MustParse(hours + ":" + seconds)
}

func ToTypeRepeatTimes(times []time.Time) (rt []RepeatTimes) {
	rt = make([]RepeatTimes, len(times))
	for i := range times {
		rt[i].At = times[i]
		rt[i].Once = sync.Once{}
	}
	return
}

func (*FakeliveController) GetSettings(c *iris.Context) {
	settings, err := GetFakeliveSettings()
	if err != nil {
		log.Println("we got error", err)
		c.JSON(FakeliveSettings{})
		return
	}
	log.Println("no error")
	c.JSON(settings)
}

func (*FakeliveController) GetNewTrim(c *iris.Context) {
	var videos []Video
	err := stormdb.AllByIndex("Id", &videos)
	if err != nil {
		c.JSON(FakeliveSettings{})
		return
	}

	max := 30
	sort.Sort(ByID(videos))

	if len(videos) < max {
		max = len(videos)
	}

	c.JSON(videos[:max])
}

func (*FakeliveController) GetSearchTrim(c *iris.Context) {
	keyword := c.Param("keyword")

	var videos []Video

	max := 30

	stormdb.Bolt.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("Video"))
		c := bucket.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			if strings.Contains(strings.ToLower(string(v)), strings.ToLower(keyword)) {
				var video Video

				err := stormdb.Get("Video", k, &video)
				if err != nil {
					log.Println(err)
					continue
				}
				videos = append(videos, video)
			}

		}

		return nil
	})

	if len(videos) < max {
		max = len(videos)
	}

	c.JSON(videos[:max])

}

func (*FakeliveController) PostSaveVideoTrim(c *iris.Context) {
	var video Video

	err := c.ReadJSON(&video)
	if err != nil {
		c.JSON(FakeliveSettings{})
		return
	}

	err = stormdb.Save(video)
	if err != nil {
		c.JSON(FakeliveSettings{})
		return
	}
}

func (*FakeliveController) SetSettings(c *iris.Context) {
	var lss FakeliveSettings

	err := c.ReadJSON(&lss)
	if err != nil {
		c.Write(err.Error())
		return
	}

	lss.RepeatTimes = ToTypeRepeatTimes(lss.RTimes)

	err = SetFakeliveSettings(lss)
	if err != nil {
		c.Write(err.Error())
		return
	}

}

func (*FakeliveController) ReloadNow(c *iris.Context) {
	work()

	resp,err:=http.Get("http://opinion.azorestv.com:1935/reload")
	if err != nil {
		c.Write(err.Error())
		return
	}

	var r = &bytes.Buffer{}

	_,err=r.ReadFrom(resp.Body)
	if err != nil {
		c.Write(err.Error())
		return
	}

	if !strings.Contains(r.String(),"DONE"){
		c.Error("failed",500)

	}


}
