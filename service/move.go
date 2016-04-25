package service

import (
	"encoding/json"
	"net/http"

	"github.com/NYTimes/gizmo/server"
	"github.com/Sirupsen/logrus"
	"github.com/clawio/codes"
	"github.com/clawio/entities"
	"github.com/clawio/keys"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

// MoveObject retrieves the information about an object.
func (s *Service) MoveObject(w http.ResponseWriter, r *http.Request) {
	sourcePath := mux.Vars(r)["path"]
	targetPath := r.URL.Query().Get("target")
	user := context.Get(r, keys.UserKey).(entities.User)
	err := s.MetaDataController.MoveObject(user, sourcePath, targetPath)
	if err != nil {
		s.handleMoveObjectError(err, w)
		return
	}
}

func (s *Service) handleMoveObjectError(err error, w http.ResponseWriter) {
	if codeErr, ok := err.(*codes.Err); ok {
		if codeErr.Code == codes.NotFound {
			server.Log.WithFields(logrus.Fields{
				"error": err,
			}).Error("object not found")
			w.WriteHeader(http.StatusNotFound)
			return
		} else if codeErr.Code == codes.BadInputData {
			server.Log.WithFields(logrus.Fields{
				"error": err,
			}).Warn("object cannot be moved")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err)
			return
		}
	}
	server.Log.WithFields(logrus.Fields{
		"error": err,
	}).Error("error moving object")
	w.WriteHeader(http.StatusInternalServerError)
	return
}
