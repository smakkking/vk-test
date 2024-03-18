package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func SendJSON(w http.ResponseWriter, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, fmt.Sprintf("handlers.SendJSON: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(data)
	if err != nil {
		http.Error(w, fmt.Sprintf("handlers.SendJSON: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
}
