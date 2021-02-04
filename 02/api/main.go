package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"failed-interview/02/infrastructure/repository"
	"failed-interview/02/usecase/balance"

	"failed-interview/02/api/handler"
	"failed-interview/02/api/middleware"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = ""
	dbname   = "public"
	apiPort  = 8080
)

func main() {
	dataSourceName := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", dataSourceName)

	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		db.Close()
		log.Fatal(err)
	}

	balanceRepo := repository.NewbalancePG(db)
	balanceService := balance.NewService(balanceRepo)

	if err != nil {
		log.Fatal(err.Error())
	}

	r := mux.NewRouter()

	n := negroni.New(
		negroni.HandlerFunc(middleware.Cors),

		negroni.NewLogger(),
	)

	handler.MakebalanceHandlers(r, *n, balanceService)

	http.Handle("/", r)

	logger := log.New(os.Stderr, "logger: ", log.Lshortfile)
	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         ":" + strconv.Itoa(apiPort),
		Handler:      context.ClearHandler(http.DefaultServeMux),
		ErrorLog:     logger,
	}
	err = srv.ListenAndServe()

	if err != nil {
		log.Fatal(err.Error())
	}
}
