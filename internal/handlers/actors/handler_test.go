package actors_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	handlerActors "vk_test/internal/handlers/actors"
	"vk_test/internal/infrastucture/teststore"
	"vk_test/internal/model"
	serviceActors "vk_test/internal/services/actors"

	"github.com/stretchr/testify/require"
)

func TestCreateActor(t *testing.T) {
	mainStorage := teststore.NewStorage()
	actorService := serviceActors.NewService(mainStorage)
	actorHandler := handlerActors.NewHandler(actorService)

	server := httptest.NewServer(http.HandlerFunc(actorHandler.CreateActor))

	payload, _ := json.Marshal(model.Actor{
		Name:      "Alexandr Petrov",
		Sex:       "мужчина",
		DateBirth: model.CivilTime(time.Date(2006, 1, 23, 0, 0, 0, 0, time.UTC)),
	})
	resp, _ := http.Post(server.URL, "application/json", bytes.NewReader(payload))

	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestUpdateActor(t *testing.T) {
	mainStorage := teststore.NewStorage()
	actorService := serviceActors.NewService(mainStorage)
	actorHandler := handlerActors.NewHandler(actorService)

	server := httptest.NewServer(http.HandlerFunc(actorHandler.UpdateActor))

	resp, _ := http.Get(server.URL)

	require.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
}

func TestDeleteActor(t *testing.T) {
	mainStorage := teststore.NewStorage()
	actorService := serviceActors.NewService(mainStorage)
	actorHandler := handlerActors.NewHandler(actorService)

	server := httptest.NewServer(http.HandlerFunc(actorHandler.DeleteActor))

	resp, _ := http.Get(server.URL)

	require.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
}

func TestGetActors(t *testing.T) {
	mainStorage := teststore.NewStorage()
	actorService := serviceActors.NewService(mainStorage)
	actorHandler := handlerActors.NewHandler(actorService)

	server := httptest.NewServer(http.HandlerFunc(actorHandler.GetActors))
	resp, _ := http.Get(server.URL)

	require.Equal(t, http.StatusOK, resp.StatusCode)
}
