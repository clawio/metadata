package service

import (
	"net/http"
	"net/http/httptest"

	"github.com/clawio/codes"
	"github.com/stretchr/testify/require"
)

func (suite *TestSuite) TestDelete() {
	suite.MockMetaDataController.On("DeleteObject").Once().Return(nil)
	r, err := http.NewRequest("DELETE", deleteURL+"myblob", nil)
	setToken(r)
	require.Nil(suite.T(), err)
	w := httptest.NewRecorder()
	suite.Server.ServeHTTP(w, r)
	require.Equal(suite.T(), http.StatusOK, w.Code)
}
func (suite *TestSuite) TestDelete_withError() {
	suite.MockMetaDataController.On("DeleteObject").Once().Return(codes.NewErr(99, ""))
	r, err := http.NewRequest("DELETE", deleteURL+"myblob", nil)
	setToken(r)
	require.Nil(suite.T(), err)
	w := httptest.NewRecorder()
	suite.Server.ServeHTTP(w, r)
	require.Equal(suite.T(), http.StatusInternalServerError, w.Code)
}
