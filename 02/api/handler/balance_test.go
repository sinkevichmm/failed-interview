package handler

import (
	"failed-interview/02/entity"
	"failed-interview/02/usecase/balance/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/codegangsta/negroni"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func Test_listBalances(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	service := mock.NewMockUseCase(controller)
	r := mux.NewRouter()
	n := negroni.New()
	MakebalanceHandlers(r, *n, service)
	path, err := r.GetRoute("listBalances").GetPathTemplate()
	assert.Nil(t, err)
	assert.Equal(t, "/v1/balance", path)

	b := &entity.Balance{
		ID: 1,
	}
	service.EXPECT().
		ListBalances().
		Return([]*entity.Balance{b}, nil)

	ts := httptest.NewServer(listBalances(service))

	defer ts.Close()
	res, err := http.Get(ts.URL)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func Test_listBalances_NotFound(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	service := mock.NewMockUseCase(controller)
	ts := httptest.NewServer(listBalances(service))

	defer ts.Close()
	service.EXPECT().
		ListBalances().
		Return(nil, entity.ErrNotFound)

	res, err := http.Get(ts.URL + "?title=balance+of+balances")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

func Test_updateBalance(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	service := mock.NewMockUseCase(controller)
	r := mux.NewRouter()
	n := negroni.New()
	MakebalanceHandlers(r, *n, service)
	path, err := r.GetRoute("updateBalance").GetPathTemplate()
	assert.Nil(t, err)
	assert.Equal(t, "/v1/balance", path)

	service.EXPECT().
		UpdateBalance(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil)

	h := updateBalance(service)

	ts := httptest.NewServer(h)
	defer ts.Close()

	payload := `{"idFrom":2,"idTo":1,"value":50}`

	resp, _ := http.Post(ts.URL+"/v1/book", "application/json", strings.
		NewReader(payload))

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func Test_updateBalance_NotFound(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	service := mock.NewMockUseCase(controller)
	r := mux.NewRouter()
	n := negroni.New()
	MakebalanceHandlers(r, *n, service)
	path, err := r.GetRoute("updateBalance").GetPathTemplate()
	assert.Nil(t, err)
	assert.Equal(t, "/v1/balance", path)

	service.EXPECT().
		UpdateBalance(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(entity.ErrNotFound)

	h := updateBalance(service)

	ts := httptest.NewServer(h)
	defer ts.Close()

	payload := `{"idFrom":20,"idTo":10,"value":50}`
	resp, _ := http.Post(ts.URL+"/v1/book", "application/json", strings.
		NewReader(payload))

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func Test_updateBalance_NotFound2(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	service := mock.NewMockUseCase(controller)
	r := mux.NewRouter()
	n := negroni.New()
	MakebalanceHandlers(r, *n, service)
	path, err := r.GetRoute("updateBalance").GetPathTemplate()
	assert.Nil(t, err)
	assert.Equal(t, "/v1/balance", path)

	service.EXPECT().
		UpdateBalance(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(entity.ErrIDFromNotFound)

	h := updateBalance(service)

	ts := httptest.NewServer(h)
	defer ts.Close()

	payload := `{"idFrom":20,"idTo":1,"value":50}`
	resp, _ := http.Post(ts.URL+"/v1/book", "application/json", strings.
		NewReader(payload))

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func Test_updateBalance_NotFound3(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	service := mock.NewMockUseCase(controller)
	r := mux.NewRouter()
	n := negroni.New()
	MakebalanceHandlers(r, *n, service)
	path, err := r.GetRoute("updateBalance").GetPathTemplate()
	assert.Nil(t, err)
	assert.Equal(t, "/v1/balance", path)

	service.EXPECT().
		UpdateBalance(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(entity.ErrIDToNotFound)

	h := updateBalance(service)

	ts := httptest.NewServer(h)
	defer ts.Close()

	payload := `{"idFrom":2,"idTo":10,"value":50}`
	resp, _ := http.Post(ts.URL+"/v1/book", "application/json", strings.
		NewReader(payload))

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func Test_updateBalance_IDEqual(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	service := mock.NewMockUseCase(controller)
	r := mux.NewRouter()
	n := negroni.New()
	MakebalanceHandlers(r, *n, service)
	path, err := r.GetRoute("updateBalance").GetPathTemplate()
	assert.Nil(t, err)
	assert.Equal(t, "/v1/balance", path)

	h := updateBalance(service)

	ts := httptest.NewServer(h)
	defer ts.Close()

	payload := `{"idFrom":2,"idTo":2,"value":50}`
	resp, _ := http.Post(ts.URL+"/v1/book", "application/json", strings.
		NewReader(payload))

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
