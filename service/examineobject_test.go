package service

import (
	"net/http"
	"net/http/httptest"

	"github.com/clawio/codes"
	"github.com/clawio/entities/mocks"
	"github.com/stretchr/testify/require"
)

func (suite *TestSuite) TestExamine() {
	suite.MockAuthService.On("Verify", "").Once().Return(user, &codes.Response{}, nil)
	suite.MockMetaDataController.On("ExamineObject").Once().Return(&mocks.MockObjectInfo{}, nil)
	r, err := http.NewRequest("GET", "/clawio/v1/metadata/examine/myblob", nil)
	require.Nil(suite.T(), err)
	w := httptest.NewRecorder()
	suite.Server.ServeHTTP(w, r)
	require.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *TestSuite) TestExamine_withObjectNotFound() {
	suite.MockAuthService.On("Verify", "").Once().Return(user, &codes.Response{}, nil)
	suite.MockMetaDataController.On("ExamineObject").Once().Return(&mocks.MockObjectInfo{}, codes.NewErr(codes.NotFound, ""))
	r, err := http.NewRequest("GET", "/clawio/v1/metadata/examine/myblob", nil)
	require.Nil(suite.T(), err)
	w := httptest.NewRecorder()
	suite.Server.ServeHTTP(w, r)
	require.Equal(suite.T(), http.StatusNotFound, w.Code)
}

func (suite *TestSuite) TestExamine_withError() {
	suite.MockAuthService.On("Verify", "").Once().Return(user, &codes.Response{}, nil)
	suite.MockMetaDataController.On("ExamineObject").Once().Return(&mocks.MockObjectInfo{}, codes.NewErr(99, ""))
	r, err := http.NewRequest("GET", "/clawio/v1/metadata/examine/myblob", nil)
	require.Nil(suite.T(), err)
	w := httptest.NewRecorder()
	suite.Server.ServeHTTP(w, r)
	require.Equal(suite.T(), http.StatusInternalServerError, w.Code)
}
