package service

import (
	"net/http"
	"net/http/httptest"

	"github.com/clawio/codes"
	"github.com/clawio/entities"
	"github.com/stretchr/testify/require"
)

func (suite *TestSuite) TestExamine() {
	suite.MockMetaDataController.On("ExamineObject").Once().Return(&entities.ObjectInfo{}, nil)
	r, err := http.NewRequest("GET", examineURL+"myblob", nil)
	setToken(r)
	require.Nil(suite.T(), err)
	w := httptest.NewRecorder()
	suite.Server.ServeHTTP(w, r)
	require.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *TestSuite) TestExamine_withObjectNotFound() {
	suite.MockMetaDataController.On("ExamineObject").Once().Return(&entities.ObjectInfo{}, codes.NewErr(codes.NotFound, ""))
	r, err := http.NewRequest("GET", examineURL+"myblob", nil)
	setToken(r)
	require.Nil(suite.T(), err)
	w := httptest.NewRecorder()
	suite.Server.ServeHTTP(w, r)
	require.Equal(suite.T(), http.StatusNotFound, w.Code)
}

func (suite *TestSuite) TestExamine_withError() {
	suite.MockMetaDataController.On("ExamineObject").Once().Return(&entities.ObjectInfo{}, codes.NewErr(99, ""))
	r, err := http.NewRequest("GET", examineURL+"myblob", nil)
	setToken(r)
	require.Nil(suite.T(), err)
	w := httptest.NewRecorder()
	suite.Server.ServeHTTP(w, r)
	require.Equal(suite.T(), http.StatusInternalServerError, w.Code)
}
