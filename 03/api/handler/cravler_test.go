package handler

import (
	"failed-interview/03/entity"
	"failed-interview/03/usecase/crawler/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/codegangsta/negroni"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func Test_listLinks(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	service := mock.NewMockUseCase(controller)
	r := mux.NewRouter()
	n := negroni.New()
	MakecrawlerHandlers(r, *n, service)
	path, err := r.GetRoute("listLinks").GetPathTemplate()
	assert.Nil(t, err)
	assert.Equal(t, "/v1/crawler", path)

	l := []*entity.Links{
		{URL: "url"},
		{URL: "url"},
	}

	service.EXPECT().
		GetList(gomock.Any(), gomock.Any()).
		Return(l)

	h := listLinks(service)

	ts := httptest.NewServer(h)
	defer ts.Close()

	payload := `{"links":"https://easymb.xyz","timeout":1}`

	resp, _ := http.Post(ts.URL+"/v1/crawler", "application/json", strings.
		NewReader(payload))

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func Test_listLinks_BadRequest(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	service := mock.NewMockUseCase(controller)
	r := mux.NewRouter()
	n := negroni.New()
	MakecrawlerHandlers(r, *n, service)
	path, err := r.GetRoute("listLinks").GetPathTemplate()
	assert.Nil(t, err)
	assert.Equal(t, "/v1/crawler", path)

	h := listLinks(service)

	ts := httptest.NewServer(h)
	defer ts.Close()

	payload := `{"links":1,"timeout":1}`

	resp, _ := http.Post(ts.URL+"/v1/crawler", "application/json", strings.
		NewReader(payload))

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
