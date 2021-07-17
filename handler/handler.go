package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/marcogregorius/url-shortener/models"
	log "github.com/sirupsen/logrus"
)

func PingHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{"ping": "pong"}
	WriteJSON(w, http.StatusOK, data)
}

func GetShortlinkHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	s := models.Shortlink{Id: vars["id"]}
	var res map[string]interface{}
	var err error
	if res, err = s.GetShortlink(); err != nil {
		msg := []string{(err.Error())}
		WriteError(w, http.StatusInternalServerError, msg)
		return
	}
	WriteJSON(w, http.StatusOK, res)
}

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	s := models.Shortlink{Id: vars["id"]}
	var res map[string]interface{}
	var err error
	if res, err = s.GetShortlink(); err != nil {
		msg := []string{(err.Error())}
		WriteError(w, http.StatusInternalServerError, msg)
		return
	}
	//fmt.Println(res["source_url"])
	//source := fmt.Sprintf(res["source_url"])
	http.Redirect(w, r, res["source_url"].(string), http.StatusMovedPermanently)

	// increase visited counter
	s.Visited = s.Visited + 1
	s.LastVisitedAt.Time, s.LastVisitedAt.Valid = time.Now(), true
	err = s.VisitShortlink(true)
	if err != nil {
		log.Error(err)
	}
}

func CreateShortlinkHandler(w http.ResponseWriter, r *http.Request) {
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
	if res, err = s.CreateShortlink(); err != nil {
		msg := []string{"database error"}
		WriteError(w, http.StatusInternalServerError, msg)
		return
	}
	WriteJSON(w, http.StatusOK, res)
}
