package fakelive

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Netherdrake/youtubeId"
	"github.com/jinzhu/now"
	"github.com/palantir/stacktrace"
)

func calcScheduleDate() string {
	fls, err := GetFakeliveSettings()
	if err != nil {
		fls.StartTime = "00:00"

	}

	if fls.StartTime == "" {
		fls.StartTime = "00:00"
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
	checkWriteErr(buf.WriteString(fmt.Sprint(fls.StartTime)))
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

type SmilPlaylist struct {
	Id               []byte `storm:"id"`
	Title            string
	Thumbnail        string
	Duration         string
	StartTimeSeconds int
	EndTimeSeconds   int
	VidType          videoType
	Scheduled        time.Time
	EndTime          time.Time
	Src              string
	StartSec         int
	Lenght           int
}

func appendLatestVideos(videos []Video, starttime time.Time) (smilPlaylist []SmilPlaylist, sstarttime time.Time) {

	for i := range videos {

		var currVideo Video
		err := stormdb.One("Id", videos[i].Id, &currVideo)
		if err != nil {
			log.Fatalln(err)
		}

		if currVideo.Disabled {
			log.Println(currVideo)
			continue
		}
		duration := videos[i].DurationSeconds
		if currVideo.StartTime != 0 || currVideo.EndTime != 0 {
			duration = time.Duration(currVideo.EndTime - currVideo.StartTime) * time.Second
		}

		location, err := GetVideoLocation(currVideo.Id)
		if err != nil {
			log.Println(stacktrace.Propagate(err, "video id : %d", currVideo))
			continue
		}

		var lenght= -1

		if currVideo.EndTime !=0{
			lenght=currVideo.EndTime-currVideo.StartTime
		}



		smilPlaylist = append(smilPlaylist, SmilPlaylist{
			Title:currVideo.Title,
			Thumbnail:currVideo.Thumbnail,
			VidType:   vod,
			Scheduled: starttime,
			StartTimeSeconds:currVideo.StartTime,
			EndTimeSeconds:currVideo.EndTime,
			EndTime:   starttime.Add(duration),
			Src:       location,
			Duration:duration.String(),
			StartSec:  currVideo.StartTime,
			Lenght:    lenght,
		})

		log.Println("-->",currVideo)
		starttime = starttime.Add(duration)
	}

	sstarttime = starttime

	return smilPlaylist, sstarttime
}

func genSmilPlaylistSlice(ids []Video, startTime string) (smilPlaylist []SmilPlaylist, videos []Video) {

	settings, err := GetFakeliveSettings()
	if err != nil {
		log.Fatalln(err)
	}

	startLatestVideosTimes := settings.RepeatTimes

	var StartTime time.Time

	StartTime = now.MustParse(startTime)
	StartPlaylistTime := StartTime

	var shouldContinue = true
	var startLatestIndex int

	for i := range ids {

		var currVideo Video
		err := stormdb.One("Id", ids[i].Id, &currVideo)
		if err != nil {
			log.Fatalln(err)
		}

		if currVideo.Disabled {
			log.Println(currVideo)
			continue
		}
		duration := ids[i].DurationSeconds
		if currVideo.StartTime != 0 || currVideo.EndTime != 0 {
			duration = time.Duration(currVideo.EndTime - currVideo.StartTime) * time.Second

		}

		if StartPlaylistTime.Add(24 * time.Hour).Before(StartTime) {
			break
		}

		location, err := GetVideoLocation(currVideo.Id)
		if err != nil {
			log.Println(stacktrace.Propagate(err, "video id : %d", currVideo.Id))
			continue
		}
		endVideoTime := StartTime.Add(duration)
		//
		if startLatestIndex < len(startLatestVideosTimes) {

			//check if at the end of the video we pass the startstream time


			log.Println(endVideoTime, startLatestVideosTimes[startLatestIndex].getHourSeconds())


			//09:30  09:15
			//the livestream starts before the end of the vod
			if endVideoTime.After(startLatestVideosTimes[startLatestIndex].getHourSeconds()) {
				startLatestVideosTimes[startLatestIndex].Once.Do(func() {
					//TODO resolve this mess <
					cuttime := endVideoTime.Sub(startLatestVideosTimes[startLatestIndex].getHourSeconds())

					log.Println("cuttime",cuttime)



					playtime := duration - cuttime

					smilPlaylist = append(smilPlaylist, SmilPlaylist{
						Title:currVideo.Title,
						Thumbnail:currVideo.Thumbnail,
						VidType:   vod,
						Duration:playtime.String(),
						StartTimeSeconds:currVideo.StartTime,
						EndTimeSeconds:currVideo.EndTime,
						Scheduled: StartTime,
						EndTime:   StartTime.Add(playtime),
						Src:       location,
						StartSec:  0,
						Lenght:    int(playtime / time.Second),
					})

					videos = append(videos, currVideo)

					StartTime = StartTime.Add(playtime)

					var vids []SmilPlaylist
					vids, StartTime = appendLatestVideos(ids[:3], StartTime)
					videos = append(videos, ids[:3]...)

					smilPlaylist = append(smilPlaylist, vids...)

					shouldContinue = true
					startLatestIndex++

				})

				if shouldContinue {
					shouldContinue = false
					continue
				}
				//calculate when we need to stop the vod

			}
		}

		smilPlaylist = append(smilPlaylist, SmilPlaylist{
			Title:currVideo.Title,
			Thumbnail:currVideo.Thumbnail,
			VidType:   vod,
			Duration:duration.String(),
			Scheduled: StartTime,
			StartTimeSeconds:currVideo.StartTime,
			EndTimeSeconds:currVideo.EndTime,
			EndTime:   StartTime.Add(duration),
			Src:       location,
			StartSec:  0,
			Lenght:    -1,
		})

		videos = append(videos, currVideo)

		StartTime = StartTime.Add(duration)
	}

	return
}

func genSmil(smilp []SmilPlaylist) string {
	const tpllive = `<playlist name="pl%d" playOnStream="fakelive" repeat="true" scheduled="%s">
			<video src="%s" start="%d" length="%d"/>
		</playlist>`
	smil := bytes.NewBuffer([]byte(""))
	smil.WriteString(
		`<smil>
		    <head>
            </head>
            <body>
                <stream name="fakelive"></stream>

`)
	for i := range smilp {
		var lenght = -1
		if smilp[i].EndTimeSeconds != 0 {
			lenght = smilp[i].EndTimeSeconds - smilp[i].StartTimeSeconds
		}
		smil.WriteString(fmt.Sprintf(tpllive, i, timeFormated(smilp[i].Scheduled), smilp[i].Src, smilp[i].StartTimeSeconds, lenght) + "\n")
	}

	smil.WriteString(`</body></smil>`)

	return smil.String()

}

/*

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

	startlivestreamtime, err := now.Parse(liveStream.StartTime)
	if err != nil {
		startlivestreamtime = time.Time{}
		log.Println(err)
	}
	endlivestreamtime, err := now.Parse(liveStream.EndTime)
	if err != nil {
		endlivestreamtime = time.Time{}
		log.Println(err)
	}

	if startlivestreamtime.Before(StartTime) {
		startlivestreamtime = startlivestreamtime.Add(time.Hour * 24)
		endlivestreamtime = endlivestreamtime.Add(time.Hour * 24)
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
*/

func GetVideoLocation(id int) (string, error) {
	var vid_clip Clip_files

	err := sqldb.Table("clip_files").Where("id_clip = ? and id_quality = 1", id).First(&vid_clip).Error
	if err != nil {
		return "", err
	}

	var vid Clips
	err = sqldb.Where("id = ? ", id).First(&vid).Error
	if err != nil {
		return "", err
	}

	var location string

	switch videoType(vid.Type) {
	case local:
		location = vid_clip.Vod_flash

	case vod:
		if strings.Contains(vid_clip.Vod_flash, "rtmp") {
			res := strings.Split(vid_clip.Vod_flash, "/")

			location = res[len(res) - 1]
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
