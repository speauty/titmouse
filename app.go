package main

import (
	"fmt"
	"titmouse/lib/log"
	"titmouse/model"
)

func main() {
	log.Api().Init(nil)

	tmpSettings := new(model.Settings)
	fmt.Println(tmpSettings.Load())
}
