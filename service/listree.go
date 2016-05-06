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

// ListTree retrieves the information about an object.
func (s *Service) ListTree(w http.ResponseWriter, r *http.Request) {
	path := mux.Vars(r)["path"]
	user := context.Get(r, keys.UserKey).(*entities.User)
	oinfos, err := s.MetaDataController.ListTree(user, path)
	if err != nil {
		s.handleListTreeError(err, w)
		return
	}
	if err := json.NewEncoder(w).Encode(oinfos); err != nil {
		s.handleListTreeError(err, w)
		return
	}
}

func (s *Service) handleListTreeError(err error, w http.ResponseWriter) {
	if codeErr, ok := err.(*codes.Err); ok {
		if codeErr.Code == codes.NotFound {
			server.Log.WithFields(logrus.Fields{
				"error": err,
			}).Warn("object not found")
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if codeErr.Code == codes.BadInputData {
			server.Log.WithFields(logrus.Fields{
				"error": err,
			}).Warn("object is not a tree")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err)
			return
		}
	}
	server.Log.WithFields(logrus.Fields{
		"error": err,
	}).Error("error listing tree")
	w.WriteHeader(http.StatusInternalServerError)
	return
}
