package main

import (

	"github.com/carlescere/scheduler"
	"log"
	"net/http"
	"bytes"
)

var basedir="V:\\"

func main() {

	var err error
	_, err = scheduler.Every().Day().At("07:10").Run(func() {
		work()

	})

	if err != nil {
		log.Fatalln(err)

	}
}

func work()error{
	err:=downloadMissingYoutubeVideos()
	if err!=nil{
		return err
	}

	resp,err:=http.Get("http://localhost:1935/reload")
	if err!=nil{
		return err
	}

	var buf =&bytes.Buffer{}
	defer resp.Body.Close()
	_,err=buf.ReadFrom(resp.Body)
	if err!=nil{
		return err
	}

	_=buf


	return nil
}
