package fakelive

import (
	"github.com/asdine/storm"
	"time"
	"log"
)

func syncPlaylist(videos []Video) error {
	now := time.Now()
	for i := range videos {
		var video Video
		err := stormdb.One("Id", videos[i].Id, &video)
		if err == storm.ErrNotFound || err == storm.ErrIndexNotFound {
			er := stormdb.Save(videos[i])
			if er != nil {
				log.Println(er)
				log.Println(videos[i], videos[i].Id)
				return err
			}
			continue

		} else if err == nil {
			continue
		}
		log.Println(videos[i], videos[i].Id)
		return err
	}

	log.Println("second")
	log.Println(time.Since(now))
	return nil
}

