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

	const longForm = "2006-01-02 15:04:05"
	t, _ := time.Parse(longForm, calcScheduleDate())
	err = SaveCurrentPlaylist(&Playlist{Videos: videos, StartTime: t})
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
	seconds := fmt.Sprintf("%02d", r.At.Second())
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
	cmd := exec.Command("service", "WowzaStreamingEngine", "restart")
	err := cmd.Start()
	if err != nil {
		c.Write(err.Error())
		log.Fatalln(stacktrace.Propagate(err, ""))
		return
	}
	cmd.Wait()

}
