package main

import (
	"failed-interview/01/internal/server"
	"flag"
	"log"
)

func main() {
	httpport := flag.String("httpport", "", "httpserver port")
	grpcaddress := flag.String("grpcaddress", "", "grpc address:port")
	auth := flag.String("auth", "", "login:password")

	flag.Parse()

	m := server.NewApp(*httpport, *grpcaddress, *auth)

	log.Fatalln(m.Start().Error())
}
