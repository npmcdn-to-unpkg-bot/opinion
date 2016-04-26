package main

import (
	"github.com/otium/ytdl"

	"fmt"
	"github.com/Netherdrake/youtubeId"
	"github.com/jinzhu/gorm"
	"github.com/pillash/mp4util"
	"log"
	"os"
	"path/filepath"
)

type youtubeIdMapping struct {
	YtID     string
	VID      int
	Duration int
}

type Youtube struct {
	BaseDir          string
	Youtube2Platform map[string]*youtubeIdMapping
	Platform2Youtube map[int]*youtubeIdMapping
	db               *gorm.DB
}

func NewYoutubeDl(basedir string, db *gorm.DB) *Youtube {
	return &Youtube{
		BaseDir:          basedir,
		Youtube2Platform: make(map[string]*youtubeIdMapping),
		Platform2Youtube: make(map[int]*youtubeIdMapping),
		db:               db,
	}
}

func (yt *Youtube) DownloadVideo(id string) error {

	_, err := os.Stat(filepath.Join(yt.BaseDir, id+".mp4"))
	//if not exists local download it
	if os.IsNotExist(err) {
		vid, err := ytdl.GetVideoInfo("https://www.youtube.com/watch?v=" + id)
		if err != nil {
			return err
		}
		if len(vid.Formats) < 1 {
			return fmt.Errorf("No available Videos to Download %+v", vid)
		}

		var index int
		for i := range vid.Formats {
			if vid.Formats[i].Resolution == "1080p" && vid.Formats[i].Extension == "mp4" {
				index = i
				break
			}
		}



		out, err := os.Create(filepath.Join(yt.BaseDir, id+".mp4"))
		if err != nil {
			return err
		}
		defer out.Close()

		err = vid.Download(vid.Formats[index], out)
		if err != nil {
			return err
		}




	}

	return nil
}

func (yt *Youtube) GetDuration(id string) (int, error) {
	return mp4util.Duration(filepath.Join(yt.BaseDir, id+".mp4"))
}

func (yt *Youtube) GetAllIds() ([]string, error) {

	var (
		clip_files []Clip_files
		err        error
	)

	yt.db.DB()
	err = yt.db.Table("clip_files").Where("embed_flash <> ? and Id_quality = 1", "").Find(&clip_files).Error
	if err != nil {
		return nil, err
	}

	var ids []string

	for i := range clip_files {
		ytId, errr := youtubeId.New().Parse(clip_files[i].Embed_flash)
		if errr != nil {
			log.Printf("error getting video id for: %s id %d", clip_files[i].Embed_flash, clip_files[i].Id_clip)
			continue
		}

		ids = append(ids, ytId)

		yt.Youtube2Platform[ytId] = &youtubeIdMapping{
			YtID: ytId,
			VID:  int(clip_files[i].Id_clip),
		}

		yt.Platform2Youtube[int(clip_files[i].Id_clip)] = &youtubeIdMapping{
			YtID: ytId,
			VID:  int(clip_files[i].Id_clip),
		}
	}
	return ids, nil
}
