package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/marcogregorius/url-shortener/models"
)

func PingHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{"ping": "pong"}
	WriteJSON(w, http.StatusOK, data)
}

func GetShortlinkHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	s := models.Shortlink{Id: vars["id"]}
	if err := s.GetShortlink(); err != nil {
		msg := []string{(err.Error())}
		WriteError(w, http.StatusInternalServerError, msg)
		return
	}
	WriteJSON(w, http.StatusOK, s)
}

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	s := models.Shortlink{Id: vars["id"]}
	if err := s.GetShortlink(); err != nil {
		msg := []string{(err.Error())}
		WriteError(w, http.StatusInternalServerError, msg)
		return
	}
	http.Redirect(w, r, s.SourceUrl, http.StatusMovedPermanently)

	// increase visited counter
	s.Visited = s.Visited + 1
	s.LastVisitedAt.Time = time.Now()
	s.UpdateShortlink()
}

func CreateShortlinkHandler(w http.ResponseWriter, r *http.Request) {
	var s models.Shortlink
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		msg := []string{"invalid JSON body"}
		WriteError(w, http.StatusBadRequest, msg)
		return
	}

	// Do validation on JSON body
	if err := Validate(s); err != nil {
		WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := s.CreateShortlink(); err != nil {
		msg := []string{"database error"}
		WriteError(w, http.StatusInternalServerError, msg)
		return
	}
	WriteJSON(w, http.StatusOK, s)
}
