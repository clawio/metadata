package service

import (
	"net/http"

	"github.com/NYTimes/gizmo/server"
	"github.com/Sirupsen/logrus"
	"github.com/clawio/entities"
	"github.com/clawio/keys"
	"github.com/gorilla/context"
)

// Init retrieves the information about an object.
func (s *Service) Init(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, keys.UserKey).(entities.User)
	err := s.MetaDataController.Init(user)
	if err != nil {
		s.handleInitError(err, w)
		return
	}
}

func (s *Service) handleInitError(err error, w http.ResponseWriter) {
	server.Log.WithFields(logrus.Fields{
		"error": err,
	}).Error("error creating user home tree")
	w.WriteHeader(http.StatusInternalServerError)
	return
}
