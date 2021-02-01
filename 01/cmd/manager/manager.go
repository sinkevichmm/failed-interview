package main

import (
	"failed-interview/01/internal/manager"
	"flag"
	"log"
)

func main() {
	meta := flag.String("meta", "", "meta file for file storage")
	limit := flag.Uint("limit", 0, "limit file storage capacity")
	port := flag.String("port", "", "port to listen")
	auth := flag.String("auth", "", "auth key")

	flag.Parse()

	m := manager.NewApp(meta, limit, port, auth)

	log.Fatalln(m.GRPCFS.Start().Error())
}
