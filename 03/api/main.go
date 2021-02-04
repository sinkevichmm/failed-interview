package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"failed-interview/03/infrastructure/repository"
	"failed-interview/03/usecase/crawler"

	"failed-interview/03/api/handler"
	"failed-interview/03/api/middleware"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	apiPort = 8080
)

func main() {
	crawlerRepo := repository.NewCrawler()
	crawlerService := crawler.NewService(crawlerRepo)

	r := mux.NewRouter()

	n := negroni.New(
		negroni.HandlerFunc(middleware.Cors),

		negroni.NewLogger(),
	)

	handler.MakecrawlerHandlers(r, *n, crawlerService)

	http.Handle("/", r)

	logger := log.New(os.Stderr, "logger: ", log.Lshortfile)
	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         ":" + strconv.Itoa(apiPort),
		Handler:      context.ClearHandler(http.DefaultServeMux),
		ErrorLog:     logger,
	}
	err := srv.ListenAndServe()

	if err != nil {
		log.Fatal(err.Error())
	}
}
