package main

import (

	"log"


	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/palantir/stacktrace"

)

var (
	sqldb               *gorm.DB
)

func init() {
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	var err error

	sqldb, err = gorm.Open(
		"mysql",
		"thesyncim:Kirk1zodiak@tcp(azorestv.com:3306)/azorestv?charset=utf8&parseTime=True",
	)

	if err != nil {
		log.Fatal(stacktrace.Propagate(err, "error connect mysql server"))
	}

}
