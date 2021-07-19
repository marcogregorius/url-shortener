package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/marcogregorius/url-shortener/models"
	log "github.com/sirupsen/logrus"
)

func (a *App) PingHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{"ping": "pong"}
	WriteJSON(w, http.StatusOK, data)
}

func (a *App) handleGetShortlink(w http.ResponseWriter, r *http.Request) (models.Shortlink, map[string]interface{}, error) {
	vars := mux.Vars(r)
	s := models.Shortlink{Id: vars["id"]}
	var res map[string]interface{}
	var err error
	if res, err = s.GetShortlink(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			msg := []string{("Shortlink not found")}
			WriteError(w, http.StatusNotFound, msg)
		default:
			msg := []string{(err.Error())}
			WriteError(w, http.StatusInternalServerError, msg)
		}
		return s, nil, err
	}
	return s, res, nil
}

func (a *App) GetShortlinkHandler(w http.ResponseWriter, r *http.Request) {
	_, res, err := a.handleGetShortlink(w, r)
	if err != nil {
		return
	}
	WriteJSON(w, http.StatusOK, res)
}

func (a *App) RedirectHandler(w http.ResponseWriter, r *http.Request) {
	s, res, err := a.handleGetShortlink(w, r)
	if err != nil {
		return
	}

	http.Redirect(w, r, res["source_url"].(string), http.StatusMovedPermanently)

	// update LastVisitedAt and Visited
	s.LastVisitedAt.Time, s.LastVisitedAt.Valid = time.Now(), true
	err = s.VisitShortlink(a.DB)
	if err != nil {
		log.Error(err)
	}
}

func (a *App) CreateShortlinkHandler(w http.ResponseWriter, r *http.Request) {
	var s models.Shortlink
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		msg := []string{"invalid JSON body"}
		fmt.Println(err)
		WriteError(w, http.StatusBadRequest, msg)
		return
	}

	// Do validation on JSON body
	if err := Validate(s); err != nil {
		WriteError(w, http.StatusBadRequest, err)
		return
	}

	var res map[string]interface{}
	var err error
	if res, err = s.CreateShortlink(a.DB); err != nil {
		msg := []string{"database error"}
		WriteError(w, http.StatusInternalServerError, msg)
		return
	}
	WriteJSON(w, http.StatusOK, res)
}
