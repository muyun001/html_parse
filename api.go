package main

import (
	"html_parse_api/routers"
	"log"
)

func main() {
	router := routers.Load()

	err := router.Run(":9020")
	if err != nil {
		log.Fatalln(err)
	}
}
