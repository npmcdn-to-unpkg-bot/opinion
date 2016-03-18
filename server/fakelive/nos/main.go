package main

import (
	"bytes"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/core/http"
	"github.com/Netherdrake/youtubeId"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/now"
	"github.com/mongodb/mongo-tools/common/json"
	"github.com/otium/ytdl"
	"github.com/palantir/stacktrace"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	db *gorm.DB
)

func init() {
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	var err error

	db, err = gorm.Open(
		"mysql",
		"thesyncim:Kirk1zodiak@tcp(azorestv.com:3306)/azorestv?charset=utf8&parseTime=True",
	)

	if err != nil {
		log.Fatal(stacktrace.Propagate(err, "error connect mysql server"))
	}

}

type youtubeIdMapping struct {
	YtID     string
	VID      int64
	Duration int64
}

var Youtube2Platform = map[string]*youtubeIdMapping{}
var Platform2Youtube = map[int64]*youtubeIdMapping{}

func main() {
	downloadMissingYoutubeVideos()
	getPlaylist()
	gensheculerPlaylist()

}

type Playlist struct {
	Videos    []Video
	StartTime time.Time
}

type Video struct {
	Id              int           `xml:"id,attr"`
	Id_user         int           `xml:"id_user,attr"`
	Title           string        `xml:"title,attr"`
	Thumbnail       string        `xml:"thumbnail,attr"`
	Poster          string        `xml:"poster,attr"`
	Type            int           `xml:"type,attr"`
	Duration        string        `xml:"duration,attr"`
	DurationSeconds time.Duration `xml:"-"`
	YoutubeId       string
	Is_ad           int `xml:"is_ad,attr"`
}

func gensheculerPlaylist() {

	resp, err := http.Get("http://opinion.azorestv.com/api/fakelive/getplaylist")
	if err != nil {
		log.Fatalln(err)
	}

	var buf = bytes.NewBuffer([]byte(""))
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	var p Playlist

	err = json.Unmarshal(buf.Bytes(), &p)
	if err != nil {
		log.Fatalln(err)
	}

	buf.Reset()

	resp, err = http.Get("http://opinion.azorestv.com/api/fakelive/livestreamset")
	if err != nil {
		log.Fatalln(err)
	}

	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	var liveStream LiveStreamSettings

	err = json.Unmarshal(buf.Bytes(), &liveStream)
	if err != nil {
		log.Fatalln(err)
	}

	playlistfile, err := os.Create("playlistScheduler.txt")
	if err != nil {
		log.Fatalln(err)
	}
	var sum time.Time

	var timetostr = func(t time.Time) string {

		return fmt.Sprintf("%02d", t.Hour()) + ":" + fmt.Sprintf("%02d", t.Minute()) + ":" + fmt.Sprintf("%02d", t.Second())

	}

	//Abertura das Festas da Agualva 2014 - parte 2 de 3(1).mp4,11:49:00,V:\AllMyTube Downloaded\Abertura das Festas da Agualva 2014 - parte 2 de 3(1).mp4,12:39:13,0.0,0,1,3013.68

	sum = p.StartTime

	startlivestreamtime := now.MustParse(liveStream.StartTime).Add(-1 * time.Hour)
	endlivestreamtime := now.MustParse(liveStream.EndTime).Add(-1 * time.Hour)
	var once sync.Once

	var onceagain = 0

	sum = sum.Add(-1 * time.Hour)
	for i := range p.Videos {
		startv := timetostr(sum)
		endv := timetostr(sum.Add(p.Videos[i].DurationSeconds))

		if liveStream.Activated {
			//check if at the end of the video we pass the startstream time

			endVideoTime := sum.Add(p.Videos[i].DurationSeconds)

			//the livestream starts before the end of the vod
			if endVideoTime.After(startlivestreamtime) {

				once.Do(func() {

					//TODO resolv this mess <
					cuttime := endVideoTime.Sub(startlivestreamtime)

					playtime := p.Videos[i].DurationSeconds - cuttime

					playlistfile.WriteString(p.Videos[i].YoutubeId + ".mp4" + "," + startv + "," + filepath.Join(basedir, p.Videos[i].YoutubeId+".mp4") + "," + timetostr(sum.Add(playtime)) + ",0.0,0,1," + fmt.Sprint(int(playtime)/int(time.Second)) + "\r\n")

					sum = sum.Add(playtime)

					var livestreamdur time.Duration

					livestreamdur = endlivestreamtime.Sub(startlivestreamtime)

					playlistfile.WriteString("live Azorestv" + "," + timetostr(startlivestreamtime) + "," + "http://azorestv.com:1935/live2/azorestv/playlist.m3u8" + "," + timetostr(endlivestreamtime) + ",0.0,0,1," + fmt.Sprint(int(livestreamdur)/int(time.Second)) + "\r\n")

					sum = endlivestreamtime

				})
				if onceagain == 0 {
					onceagain++
					continue
				}
				//calculate when we need to stop the vod

			}
		}

		if p.Videos[i].Type == 3 {
			_ = strings.Replace(p.Videos[i].Title, ",", "", -1)

			playlistfile.WriteString(p.Videos[i].YoutubeId + ".mp4" + "," + startv + "," + filepath.Join(basedir, p.Videos[i].YoutubeId+".mp4") + "," + endv + ",0.0,0,1," + fmt.Sprint(int(p.Videos[i].DurationSeconds)/int(time.Second)) + "\r\n")

		}

		sum = sum.Add(p.Videos[i].DurationSeconds)
	}
	playlistfile.Close()

}

func getPlaylist() {

	resp, err := http.Get("http://opinion.azorestv.com/api/fakelive/getplaylist")
	if err != nil {
		log.Fatalln(err)
	}

	var buf = bytes.NewBuffer([]byte(""))
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	var p Playlist

	err = json.Unmarshal(buf.Bytes(), &p)
	if err != nil {
		log.Fatalln(err)
	}

	playlistfile, err := os.Create("playlist.txt")
	if err != nil {
		log.Fatalln(err)
	}
	var sum time.Time

	var timetostr = func(t time.Time) string {

		return fmt.Sprintf("%02d", t.Hour()) + ":" + fmt.Sprintf("%02d", t.Minute())

	}
	sum = p.StartTime
	for i := range p.Videos {
		playlistfile.WriteString(p.Videos[i].Title + " no ar às: " + timetostr(sum) + " - " + timetostr(sum.Add(p.Videos[i].DurationSeconds)) + "\r\n")
		sum = sum.Add(p.Videos[i].DurationSeconds)
	}
	playlistfile.Close()

}

func getYoutubeVideoIds() (list []string, err error) {

	var clip_files []Clip_files

	db.DB()
	err = db.Table("clip_files").Where("embed_flash <> ? and Id_quality = 1", "").Find(&clip_files).Error
	if err != nil {
		return
	}

	for i := range clip_files {

		ytId, errr := youtubeId.New().Parse(clip_files[i].Embed_flash)
		if errr != nil {
			log.Println(stacktrace.Propagate(errr, "error getting video id for: %s id %d", clip_files[i].Embed_flash, clip_files[i].Id_clip))
			continue
		}

		list = append(list, ytId)

		Youtube2Platform[ytId] = &youtubeIdMapping{
			YtID: ytId,
			VID:  clip_files[i].Id_clip,
		}

		Platform2Youtube[clip_files[i].Id_clip] = &youtubeIdMapping{
			YtID: ytId,
			VID:  clip_files[i].Id_clip,
		}
	}
	return
}

var basedir = "v:\\yt"

func downloadMissingYoutubeVideos() error {

	allvid, err := getYoutubeVideoIds()
	if err != nil {
		return err
	}

	var failedctr int

	for i := range allvid {
		_, err := os.Stat(filepath.Join(basedir, allvid[i]+".mp4"))
		//if not exists local download it
		if os.IsNotExist(err) {
			//getVideoInfo

			vid, err := ytdl.GetVideoInfo("https://www.youtube.com/watch?v=" + allvid[i])

			if err != nil {
				failedctr++
				log.Println(stacktrace.Propagate(err, "error getting youtube video info:  %s -> %d", allvid[i]), Youtube2Platform[allvid[i]])
				continue
			}

			file, _ := os.Create(filepath.Join(basedir, (allvid[i])+".mp4"))
			defer file.Close()

			for i := range vid.Formats {
				log.Println(vid.Formats[i])

				break
			}

			if len(vid.Formats) < 1 {
				log.Println(stacktrace.Propagate(err, "error getting youtube video sources to download:  %s", allvid[i]))
				continue

			}

			err = vid.Download(vid.Formats[0], file)
			if err != nil {
				failedctr++
				log.Println(stacktrace.Propagate(
					err,
					"failed to Download video %s to file %s",
					allvid[i],
					filepath.Join(basedir, (allvid[i])+".mp4")),
				)
			}
		}
	}

	return nil
}

type LiveStreamSettings struct {
	StartTime string
	EndTime   string
	Activated bool
}

type Clip_files struct {
	Embed_flash            string `json:"embed_flash"`
	Embed_html5            string `json:"embed_html5"`
	Encoding_source        string `json:"encoding_source"`
	Id                     int64  `json:"id,string"`
	Id_clip                int64  `json:"id_clip,string"`
	Id_quality             int64  `json:"id_quality,string"`
	Live_flash             string `json:"live_flash"`
	Live_html5_dash        string `json:"live_html5_dash"`
	Live_ios               string `json:"live_ios"`
	Live_ms                string `json:"live_ms"`
	Live_rtsp              string `json:"live_rtsp"`
	Vod_flash              string `json:"vod_flash"`
	Vod_flash_trailer      string `json:"vod_flash_trailer"`
	Vod_html5_h264         string `json:"vod_html5_h264"`
	Vod_html5_h264_trailer string `json:"vod_html5_h264_trailer"`
	Vod_html5_webm         string `json:"vod_html5_webm"`
	Vod_html5_webm_trailer string `json:"vod_html5_webm_trailer"`
}
