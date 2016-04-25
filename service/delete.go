package service

import (
	"net/http"

	"github.com/NYTimes/gizmo/server"
	"github.com/Sirupsen/logrus"
	"github.com/clawio/entities"
	"github.com/clawio/keys"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

// DeleteObject retrieves the information about an object.
func (s *Service) DeleteObject(w http.ResponseWriter, r *http.Request) {
	path := mux.Vars(r)["path"]
	user := context.Get(r, keys.UserKey).(entities.User)
	err := s.MetaDataController.DeleteObject(user, path)
	if err != nil {
		s.handleDeleteObjectError(err, w)
		return
	}
}

func (s *Service) handleDeleteObjectError(err error, w http.ResponseWriter) {
	server.Log.WithFields(logrus.Fields{
		"error": err,
	}).Error("error deleting object")
	w.WriteHeader(http.StatusInternalServerError)
	return
}
