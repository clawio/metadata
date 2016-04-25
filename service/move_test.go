package service

import (
	"net/http"
	"net/http/httptest"

	"github.com/clawio/codes"
	"github.com/stretchr/testify/require"
)

func (suite *TestSuite) TestMove() {
	suite.MockAuthService.On("Verify", "").Once().Return(user, &codes.Response{}, nil)
	suite.MockMetaDataController.On("MoveObject").Once().Return(nil)
	r, err := http.NewRequest("POST", "/clawio/v1/metadata/move/myblob", nil)
	require.Nil(suite.T(), err)
	values := r.URL.Query()
	values.Set("target", "otherblob")
	r.URL.RawQuery = values.Encode()
	w := httptest.NewRecorder()
	suite.Server.ServeHTTP(w, r)
	require.Equal(suite.T(), http.StatusOK, w.Code)
}
func (suite *TestSuite) TestMove_withNotFoundError() {
	suite.MockAuthService.On("Verify", "").Once().Return(user, &codes.Response{}, nil)
	suite.MockMetaDataController.On("MoveObject").Once().Return(codes.NewErr(codes.NotFound, ""))
	r, err := http.NewRequest("POST", "/clawio/v1/metadata/move/myblob", nil)
	require.Nil(suite.T(), err)
	w := httptest.NewRecorder()
	suite.Server.ServeHTTP(w, r)
	require.Equal(suite.T(), http.StatusNotFound, w.Code)
}
func (suite *TestSuite) TestMove_withBadInputError() {
	suite.MockAuthService.On("Verify", "").Once().Return(user, &codes.Response{}, nil)
	suite.MockMetaDataController.On("MoveObject").Once().Return(codes.NewErr(codes.BadInputData, ""))
	r, err := http.NewRequest("POST", "/clawio/v1/metadata/move/myblob", nil)
	require.Nil(suite.T(), err)
	w := httptest.NewRecorder()
	suite.Server.ServeHTTP(w, r)
	require.Equal(suite.T(), http.StatusBadRequest, w.Code)
}
func (suite *TestSuite) TestMove_withError() {
	suite.MockAuthService.On("Verify", "").Once().Return(user, &codes.Response{}, nil)
	suite.MockMetaDataController.On("MoveObject").Once().Return(codes.NewErr(99, ""))
	r, err := http.NewRequest("POST", "/clawio/v1/metadata/move/myblob", nil)
	require.Nil(suite.T(), err)
	w := httptest.NewRecorder()
	suite.Server.ServeHTTP(w, r)
	require.Equal(suite.T(), http.StatusInternalServerError, w.Code)
}
