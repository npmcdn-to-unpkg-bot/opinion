package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"

	"log"
)

func main() {

	db, err := gorm.Open(
		"mysql",
		"thesyncim:Kirk1zodiak@tcp(azorestv.com:3306)/azorestv?charset=utf8&parseTime=True",
	)

	if err != nil {
		log.Println(err, "error connect mysql server")
	}

	youtube := NewYoutubeDl("V:\\yt", db)
	_ = youtube

	ids,err :=youtube.GetAllIds()
	if err!=nil{
		log.Fatalln(err)
	}

	for i:=range ids{
		err=youtube.DownloadVideo(ids[i])
		if err!=nil{
			log.Fatalln(err)
		}
	}

	p,err:=GetPlaylist()
	if err!=nil{
		log.Fatalln(err)
	}

	GenPlaylist(p)
}
