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

// ExamineObject retrieves the information about an object.
func (s *Service) ExamineObject(w http.ResponseWriter, r *http.Request) {
	path := mux.Vars(r)["path"]
	user := context.Get(r, keys.UserKey).(*entities.User)
	oinfo, err := s.MetaDataController.ExamineObject(user, path)
	if err != nil {
		s.handleExamineObjectError(err, w)
		return
	}
	if err := json.NewEncoder(w).Encode(oinfo); err != nil {
		s.handleExamineObjectError(err, w)
		return
	}
}

func (s *Service) handleExamineObjectError(err error, w http.ResponseWriter) {
	if codeErr, ok := err.(*codes.Err); ok {
		if codeErr.Code == codes.NotFound {
			server.Log.WithFields(logrus.Fields{
				"error": err,
			}).Error("object not found")
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}
	server.Log.WithFields(logrus.Fields{
		"error": err,
	}).Error("error examining object")
	w.WriteHeader(http.StatusInternalServerError)
	return
}
