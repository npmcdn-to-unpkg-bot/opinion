package fakelive2

import (
	"log"
	"os"
	"path/filepath"

	"github.com/Netherdrake/youtubeId"
	"github.com/otium/ytdl"
	"github.com/palantir/stacktrace"
	"github.com/pillash/mp4util"
)

type youtubeIdMapping struct {
	YtID     string
	VID      int64
	Duration int64
}

var Youtube2Platform = map[string]*youtubeIdMapping{}
var Platform2Youtube = map[int64]*youtubeIdMapping{}

func getYoutubeVideoDuration(ytid string) (int, error) {

	return mp4util.Duration(basedir + "/" + ytid + ".mp4")

}

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

func getYoutubeVideoIds() (list []string, err error) {

	var clip_files []Clip_files

	sqldb.DB()
	err = sqldb.Table("clip_files").Where("embed_flash <> ? and Id_quality = 1", "").Find(&clip_files).Error
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
