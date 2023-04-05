package main

import (
	"fmt"
	"titmouse/lib/db"
	"titmouse/lib/log"
)

func main() {
	log.Api().Init(nil)

	if err := db.Api().Init(&db.Cfg{
		Type: db.SQLite,
		Dsn:  "titmouse.db",
		Pool: &db.CfgPool{
			MaxIdle:     10,
			MaxOpen:     50,
			MaxLifeTime: 3600,
		},
	}); err != nil {
		panic(err)
	}

	fmt.Println("hello world")
}
