package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"failed-interview/03/api/presenter"
	"failed-interview/03/usecase/crawler"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

func listLinks(service crawler.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Links   string `json:"links"`
			Timeout uint   `json:"timeout"`
		}

		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusBadRequest)

			w.Write([]byte(err.Error()))
			return
		}

		l := strings.Split(input.Links, "\n")

		data := service.GetList(l, input.Timeout)

		w.Header().Set("Content-Type", "application/json")

		var toJ []*presenter.Links
		for _, d := range data {
			toJ = append(toJ, &presenter.Links{
				URL:   d.URL,
				Title: d.Title,
				Error: d.Error,
			})
		}
		if err := json.NewEncoder(w).Encode(toJ); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal error"))
		}
	})
}

func MakecrawlerHandlers(r *mux.Router, n negroni.Negroni, service crawler.UseCase) {
	r.Handle("/v1/crawler", n.With(
		negroni.Wrap(listLinks(service)),
	)).Methods("POST", "OPTIONS").Name("listLinks")
}
