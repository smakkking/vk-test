package films_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"vk_test/internal/infrastucture/teststore"
	"vk_test/internal/model"

	handlerFilms "vk_test/internal/handlers/films"
	serviceFilms "vk_test/internal/services/films"

	"github.com/stretchr/testify/require"
)

func TestCreateFilm(t *testing.T) {
	mainStorage := teststore.NewStorage()
	filmService := serviceFilms.NewService(mainStorage)
	filmHandler := handlerFilms.NewHandler(filmService)

	server := httptest.NewServer(http.HandlerFunc(filmHandler.CreateFilm))

	payload, _ := json.Marshal(&model.Film{
		Title:        "Казаки - разбойники",
		Description:  "",
		DateCreation: model.CivilTime(time.Date(1956, 7, 14, 0, 0, 0, 0, time.UTC)),
		Rating:       9,
		ActorIDList:  []int{1, 2},
	})
	resp, _ := http.Post(server.URL, "application/json", bytes.NewReader(payload))

	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestUpdateFilm(t *testing.T) {
	mainStorage := teststore.NewStorage()
	filmService := serviceFilms.NewService(mainStorage)
	filmHandler := handlerFilms.NewHandler(filmService)

	server := httptest.NewServer(http.HandlerFunc(filmHandler.UpdateFilm))

	resp, _ := http.Get(server.URL)

	require.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
}

func TestDeleteFilm(t *testing.T) {
	mainStorage := teststore.NewStorage()
	filmService := serviceFilms.NewService(mainStorage)
	filmHandler := handlerFilms.NewHandler(filmService)

	server := httptest.NewServer(http.HandlerFunc(filmHandler.DeleteFilm))

	resp, _ := http.Get(server.URL)

	require.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
}

func TestGetFilmsSorted(t *testing.T) {
	mainStorage := teststore.NewStorage()
	filmService := serviceFilms.NewService(mainStorage)
	filmHandler := handlerFilms.NewHandler(filmService)

	server := httptest.NewServer(http.HandlerFunc(filmHandler.GetFilms))

	resp, _ := http.Get(server.URL)

	require.Equal(t, http.StatusOK, resp.StatusCode)
}
