package fakelive

import (
	"log"
	"os"
	"os/exec"

	"github.com/carlescere/scheduler"
	"github.com/gin-gonic/gin"

	"github.com/palantir/stacktrace"
	"sort"
	"time"
)

var basedir = "/var/www/vhosts/azorestv.com/httpdocs/uploads/movies/yt"
var j *scheduler.Job

func RunBackgroundScheduler() *scheduler.Job {
	work()
	strttm, _ := getStartTime()
	if strttm == "" {
		strttm = "00:00"
	}
	var err error
	j, err = scheduler.Every().Day().At(strttm).Run(func() {
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

	allvids, err := getVideoCategoryPlaylist()
	if err != nil {
		log.Fatalln(err)
	}

	videos := generatePlalist(allvids)

	sort.Sort(ByID(videos))

	videos = firstNAndShuffle(3, videos)

	smil := genSmilWithLive(videos, calcScheduleDate())

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

func HandlerCurrentPlaylist(c *gin.Context) {

	c.JSON(200, GetCurrentPlaylist())

}

func HandlerCurrentSmilPlaylist(c *gin.Context) {

	c.JSON(200, GetCurrentSmilPlaylist())

}


func HandlerGetStartTime(c *gin.Context) {

	start, err := getStartTime()
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, start)

}

func HandlerSetStartTime(c *gin.Context) {
	type startTime struct {
		StartTime string
	}
	var Ss startTime
	err := c.BindJSON(&Ss)
	if err != nil {
		c.Error(err)
		return
	}

	log.Println(Ss)
	err = SetStartTime(Ss.StartTime)
	if err != nil {
		c.Error(err)
		return
	}

	j.Quit <- true

	j, err = scheduler.Every().Day().At(Ss.StartTime).Run(func() {
		work()
	})

	if err != nil {
		log.Fatalln(stacktrace.Propagate(err, ""))

	}

}

func HandlerGetLiveStreamSettings(c *gin.Context) {

	settings, err := GetLiveStreamSettings()
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, settings)

}

func HandlerSetLiveStreamSettings(c *gin.Context) {
	var lss LiveStreamSettings

	err := c.BindJSON(&lss)
	if err != nil {
		c.Error(err)
		return
	}

	err = SetLiveStreamSettings(lss)
	if err != nil {
		c.Error(err)
		return
	}

}

func HandlerReloadNow(c *gin.Context) {
	work()
	cmd := exec.Command("service", "WowzaStreamingEngine", "restart")
	err := cmd.Start()
	if err != nil {
		c.Error(err)
		log.Fatalln(stacktrace.Propagate(err, ""))
		return
	}
	cmd.Wait()

}
