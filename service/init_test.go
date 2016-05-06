package service

import (
	"net/http"
	"net/http/httptest"

	"github.com/clawio/codes"
	"github.com/stretchr/testify/require"
)

func (suite *TestSuite) TestInit() {
	suite.MockAuthService.On("Verify", "").Once().Return(user, &codes.Response{}, nil)
	suite.MockMetaDataController.On("Init").Once().Return(nil)
	r, err := http.NewRequest("POST", initURL, nil)
	setToken(r)
	require.Nil(suite.T(), err)
	w := httptest.NewRecorder()
	suite.Server.ServeHTTP(w, r)
	require.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *TestSuite) TestInit_withError() {
	suite.MockAuthService.On("Verify", "").Once().Return(user, &codes.Response{}, nil)
	suite.MockMetaDataController.On("Init").Once().Return(codes.NewErr(99, ""))
	r, err := http.NewRequest("POST", initURL, nil)
	setToken(r)
	require.Nil(suite.T(), err)
	w := httptest.NewRecorder()
	suite.Server.ServeHTTP(w, r)
	require.Equal(suite.T(), http.StatusInternalServerError, w.Code)
}
