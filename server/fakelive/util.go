package fakelive

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Netherdrake/youtubeId"
	"github.com/jinzhu/now"
	"github.com/palantir/stacktrace"
)

//returns a playlist based on VIDEOS category
func getVideoCategoryPlaylist() ([]Video, error) {

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

func calcNextSchedule() string {
	starthour := 14
	//startminute := 00
	//output format like 2015-04-25 16:00:00
	now := time.Now()

	buf := bytes.NewBuffer([]byte(""))
	checkWriteErr(buf.WriteString(fmt.Sprint(now.Year())))
	checkWriteErr(buf.WriteString("-"))
	checkWriteErr(buf.WriteString(fmt.Sprintf("%02d", int(now.Month()))))
	checkWriteErr(buf.WriteString("-"))
	checkWriteErr(buf.WriteString(fmt.Sprint(now.Day())))
	checkWriteErr(buf.WriteString(" "))
	checkWriteErr(buf.WriteString(fmt.Sprint(starthour)))
	checkWriteErr(buf.WriteString(fmt.Sprint(":00:00")))
	return buf.String()
}

func calcScheduleDate() string {
	strttm, err := getStartTime()
	if err!=nil{
		strttm="00:00"

	}

	if strttm==""{
		strttm="00:00"
	}




	//startminute := 00
	//output format like 2015-04-25 16:00:00
	now := time.Now()

	buf := bytes.NewBuffer([]byte(""))
	checkWriteErr(buf.WriteString(fmt.Sprint(now.Year())))
	checkWriteErr(buf.WriteString("-"))
	checkWriteErr(buf.WriteString(fmt.Sprintf("%02d", int(now.Month()))))
	checkWriteErr(buf.WriteString("-"))
	checkWriteErr(buf.WriteString(fmt.Sprint(now.Day())))
	checkWriteErr(buf.WriteString(" "))
	checkWriteErr(buf.WriteString(fmt.Sprint(strttm)))
	checkWriteErr(buf.WriteString(fmt.Sprint(":00")))
	return buf.String()
}

func timeFormated(t time.Time) string {
	const dateformat = "%d-%02d-%02d %02d:%02d:%02d"
	return fmt.Sprintf(dateformat,
		t.Year(),
		t.Month(),
		t.Day(),
		t.Hour(),
		t.Minute(),
		t.Second(),
	)
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

		if videos[i].Type == embedclip {

			if val, ok := Platform2Youtube[int64(videos[i].Id)]; ok {
				var err error
				dur, err = getYoutubeVideoDuration(val.YtID)
				if err != nil {
					log.Println(stacktrace.Propagate(err, ""))
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

func checkWriteErr(n int, err error) {
	if err != nil {
		log.Println(stacktrace.Propagate(err, ""))
	}
}

func checkErr(err error) {
	if err != nil {
		log.Println(stacktrace.Propagate(err, ""))
	}
}

func getIntDuration(dur string) int {
	strdur := strings.Split(dur, ":")
	var duration int

	switch len(strdur) {
	case 1:
		secdur, _ := strconv.Atoi(strdur[0])
		duration = secdur
	case 2:
		secdur, _ := strconv.Atoi(strdur[1])
		minutesdur, _ := strconv.Atoi(strdur[0])
		duration = secdur + (minutesdur * 60)
	case 3:
		secdur, _ := strconv.Atoi(strdur[2])
		minutesdur, _ := strconv.Atoi(strdur[1])
		hourssdur, _ := strconv.Atoi(strdur[0])
		duration = secdur + (minutesdur * 60) + (hourssdur * 60 * 60)
	}

	log.Println(dur, duration)
	return duration
}

func genSmilWithLive(ids []Video, startTime string) string {
	smil := bytes.NewBuffer([]byte(""))
	_, err := smil.WriteString(
		`<smil>
		    <head>
            </head>
            <body>
                <stream name="fakelive"></stream>

`)

	const tpl = `<playlist name="pl%d" playOnStream="fakelive" repeat="true" scheduled="%s">
			<video src="%s" start="0" length="-1"/>
		</playlist>`

	const tpllive = `<playlist name="pl%d" playOnStream="fakelive" repeat="true" scheduled="%s">
			<video src="%s" start="%d" length="%d"/>
		</playlist>`

	var StartTime time.Time




	StartTime = now.MustParse(startTime)




	if err != nil {
		log.Println(stacktrace.Propagate(err, ""))

	}

	ctr := 0

	liveStream, err := GetLiveStreamSettings()
	if err != nil {
		log.Println(err)
	}

	startlivestreamtime,err := now.Parse(liveStream.StartTime)
	if err != nil {
		startlivestreamtime = time.Time{}
		log.Println(err)
	}
	endlivestreamtime,err := now.Parse(liveStream.EndTime)
	if err != nil {
		endlivestreamtime= time.Time{}
		log.Println(err)
	}
	var once sync.Once

	var onceagain = 0

	for i := range ids {

		location, err := GetVideoLocation(ids[i].Id)
		if err != nil {
			log.Println(stacktrace.Propagate(err, "video id : %d", ids[i].Id))
			continue
		}

		if liveStream.Activated {

			//check if at the end of the video we pass the startstream time

			endVideoTime := StartTime.Add(ids[i].DurationSeconds)

			//the livestream starts before the end of the vod
			if endVideoTime.After(startlivestreamtime) {

				once.Do(func() {

					//TODO resolv this mess <
					cuttime := endVideoTime.Sub(startlivestreamtime)

					playtime := ids[i].DurationSeconds - cuttime

					checkWriteErr(smil.WriteString(fmt.Sprintf(tpllive, ctr, timeFormated(StartTime), location, 0, int(playtime/time.Second)) + "\n"))
					ctr++
					StartTime = StartTime.Add(playtime)

					var livestreamdur time.Duration

					livestreamdur = endlivestreamtime.Sub(startlivestreamtime)

					checkWriteErr(smil.WriteString(fmt.Sprintf(tpllive, ctr, timeFormated(startlivestreamtime), "azorestv", -2, int(livestreamdur/time.Second)) + "\n"))
					StartTime = endlivestreamtime
					ctr++

				})
				if onceagain == 0 {
					onceagain++
					continue
				}
				//calculate when we need to stop the vod

			}

		}

		checkWriteErr(smil.WriteString(fmt.Sprintf(tpl, ctr, timeFormated(StartTime), location) + "\n"))
		StartTime = StartTime.Add(ids[i].DurationSeconds)
		ctr++
	}

	checkWriteErr(smil.WriteString(`</body></smil>`))

	return smil.String()
}

func GetVideoLocation(id int) (string, error) {

	var vid_clip Clip_files
	var vid Clips

	err := db.Table("clip_files").Where("id_clip = ? and id_quality = 1", id).First(&vid_clip).Error
	if err != nil {
		return "", err
	}

	err = db.Where("id = ? ", id).First(&vid).Error
	if err != nil {
		return "", err
	}

	var location string

	switch vid.Type {
	case local:
		location = vid_clip.Vod_flash

	case vod:
		if strings.Contains(vid_clip.Vod_flash, "rtmp") {
			res := strings.Split(vid_clip.Vod_flash, "/")

			location = res[len(res)-1]
		} else {
			location = vid_clip.Vod_flash
		}

	case live:

	case embedclip:
		id, err := youtubeId.New().Parse(vid_clip.Embed_flash)
		if err != nil {
			log.Println(stacktrace.Propagate(err, ""))
		}
		location = "yt/" + id + ".mp4"

	case autoencode:

	}
	return location, nil

}
