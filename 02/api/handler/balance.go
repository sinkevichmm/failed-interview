package handler

import (
	"encoding/json"
	"errors"
	"failed-interview/02/api/presenter"
	"log"
	"net/http"

	"failed-interview/02/usecase/balance"

	"failed-interview/02/entity"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

// TODO: прибраться

func listBalances(service balance.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error reading balances"
		var data []*entity.Balance
		var err error

		data, err = service.ListBalances()

		w.Header().Set("Content-Type", "application/json")
		if err != nil && err != entity.ErrNotFound {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}

		if data == nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(errorMessage))
			return
		}
		var toJ []*presenter.Balance
		for _, d := range data {
			toJ = append(toJ, &presenter.Balance{
				ID:    d.ID,
				Value: d.Value,
			})
		}
		if err := json.NewEncoder(w).Encode(toJ); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
		}
	})
}

func updateBalance(service balance.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var input struct {
			IDFrom int `json:"idFrom"`
			IDTo   int `json:"idTo"`
			Value  int `json:"value"`
		}
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("bad input"))
			return
		}

		if input.IDFrom == input.IDTo {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("idFrom and IdTo is equal"))
			return
		}

		err = service.UpdateBalance(input.IDFrom, input.IDTo, input.Value)

		if err != nil {
			if errors.Is(err, entity.ErrNotFound) {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(err.Error()))
			}

			if errors.Is(err, entity.ErrIDFromNotFound) {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(err.Error()))
			}

			if errors.Is(err, entity.ErrIDToNotFound) {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(err.Error()))
			}

			if errors.Is(err, entity.ErrNotEnoughbalances) {
				w.WriteHeader(http.StatusNotModified)
				w.Write([]byte(err.Error()))
			}

			return
		}
	})
}

//MakebalanceHandlers make url handlers
func MakebalanceHandlers(r *mux.Router, n negroni.Negroni, service balance.UseCase) {
	r.Handle("/v1/balance", n.With(
		negroni.Wrap(listBalances(service)),
	)).Methods("GET", "OPTIONS").Name("listBalances")

	r.Handle("/v1/balance", n.With(
		negroni.Wrap(updateBalance(service)),
	)).Methods("PUT", "OPTIONS").Name("updateBalance")
}
