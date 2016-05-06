package service

import (
	"net/http"
	"net/http/httptest"

	"github.com/clawio/codes"
	"github.com/clawio/entities"
	"github.com/stretchr/testify/require"
)

var oinfos = []*entities.ObjectInfo{}

func (suite *TestSuite) TestListTree() {
	suite.MockAuthService.On("Verify", "").Once().Return(user, &codes.Response{}, nil)
	suite.MockMetaDataController.On("ListTree").Once().Return(oinfos, nil)
	r, err := http.NewRequest("GET", listURL+"mytree", nil)
	setToken(r)
	require.Nil(suite.T(), err)
	w := httptest.NewRecorder()
	suite.Server.ServeHTTP(w, r)
	require.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *TestSuite) TestListTree_withNotFoundError() {
	suite.MockAuthService.On("Verify", "").Once().Return(user, &codes.Response{}, nil)
	suite.MockMetaDataController.On("ListTree").Once().Return(oinfos, codes.NewErr(codes.NotFound, ""))
	r, err := http.NewRequest("GET", listURL+"mytree", nil)
	setToken(r)
	require.Nil(suite.T(), err)
	w := httptest.NewRecorder()
	suite.Server.ServeHTTP(w, r)
	require.Equal(suite.T(), http.StatusNotFound, w.Code)
}
func (suite *TestSuite) TestListTree_withBadInputError() {
	suite.MockAuthService.On("Verify", "").Once().Return(user, &codes.Response{}, nil)
	suite.MockMetaDataController.On("ListTree").Once().Return(oinfos, codes.NewErr(codes.BadInputData, ""))
	r, err := http.NewRequest("GET", listURL+"mytree", nil)
	setToken(r)
	require.Nil(suite.T(), err)
	w := httptest.NewRecorder()
	suite.Server.ServeHTTP(w, r)
	require.Equal(suite.T(), http.StatusBadRequest, w.Code)
}
func (suite *TestSuite) TestListTree_withError() {
	suite.MockAuthService.On("Verify", "").Once().Return(user, &codes.Response{}, nil)
	suite.MockMetaDataController.On("ListTree").Once().Return(oinfos, codes.NewErr(99, ""))
	r, err := http.NewRequest("GET", listURL+"mytree", nil)
	setToken(r)
	require.Nil(suite.T(), err)
	w := httptest.NewRecorder()
	suite.Server.ServeHTTP(w, r)
	require.Equal(suite.T(), http.StatusInternalServerError, w.Code)
}
