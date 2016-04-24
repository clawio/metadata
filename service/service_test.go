package service

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"errors"
	"github.com/NYTimes/gizmo/config"
	"github.com/NYTimes/gizmo/server"
	"github.com/clawio/codes"
	emocks "github.com/clawio/entities/mocks"
	mock_metadatacontroller "github.com/clawio/metadata/metadatacontroller/mock"
	"github.com/clawio/sdk"
	"github.com/clawio/sdk/mocks"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var user = &emocks.MockUser{Username: "test"}

type TestSuite struct {
	suite.Suite
	MockAuthService        *mocks.MockAuthService
	MockMetaDataController *mock_metadatacontroller.MetaDataController
	SDK                    *sdk.SDK
	Service                *Service
	Server                 *server.SimpleServer
}

func Test(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (suite *TestSuite) SetupTest() {
	cfg := &Config{
		Server: &config.Server{},
		General: &GeneralConfig{
			AuthenticationServiceBaseURL: "http://localhost:58001/clawio/v1/auth/",
		},
		MetaDataController: &MetaDataControllerConfig{
			Type:              "simple",
			SimpleMetaDataDir: "/tmp",
			SimpleTempDir:     "/tmp",
		},
	}
	mockAuthService := &mocks.MockAuthService{}
	s := &sdk.SDK{}
	s.Auth = mockAuthService

	svc := &Service{}
	svc.SDK = s
	svc.Config = cfg

	mockMetaDataController := &mock_metadatacontroller.MetaDataController{}
	svc.MetaDataController = mockMetaDataController
	suite.MockMetaDataController = mockMetaDataController

	suite.Service = svc
	suite.MockAuthService = mockAuthService
	serv := server.NewSimpleServer(cfg.Server)
	serv.Register(suite.Service)
	suite.Server = serv
	// create homedir for user test
	err := os.MkdirAll("/tmp/t/test", 0755)
	require.Nil(suite.T(), err)

	// configure user mock
	user.On("GetUsername").Return("test")
}

func (suite *TestSuite) TeardownTest() {
	os.Remove("/tmp/t/test")
}

func (suite *TestSuite) TestNew() {
	cfg := &Config{
		Server: &config.Server{},
		General: &GeneralConfig{
			AuthenticationServiceBaseURL: "http://localhost:58001/clawio/v1/auth/",
		},
		MetaDataController: &MetaDataControllerConfig{
			Type:              "simple",
			SimpleMetaDataDir: "/tmp",
			SimpleTempDir:     "/tmp",
		},
	}
	svc, err := New(cfg)
	require.Nil(suite.T(), err)
	require.NotNil(suite.T(), svc)
}
func (suite *TestSuite) TestNew_withNilConfig() {
	_, err := New(nil)
	require.NotNil(suite.T(), err)
}

func (suite *TestSuite) TestNew_withNilMetaDataControllerConfig() {
	cfg := &Config{
		General: &GeneralConfig{
			AuthenticationServiceBaseURL: "http://localhost:58001/clawio/v1/auth/",
		},
		Server: nil,
	}
	_, err := New(cfg)
	require.NotNil(suite.T(), err)
}
func (suite *TestSuite) TestNew_withNilGeneralConfig() {
	cfg := &Config{
		Server:             nil,
		MetaDataController: nil,
	}
	_, err := New(cfg)
	require.NotNil(suite.T(), err)
}
func (suite *TestSuite) TestPrefix() {
	require.Equal(suite.T(), "/clawio/v1/metadata", suite.Service.Prefix())
}

func (suite *TestSuite) TestMetrics() {
	r, err := http.NewRequest("GET", "/clawio/v1/metadata/metrics", nil)
	require.Nil(suite.T(), err)
	w := httptest.NewRecorder()
	suite.Server.ServeHTTP(w, r)
	require.Equal(suite.T(), 200, w.Code)
}

func (suite *TestSuite) TestgetTokenFromRequest_header() {
	r, err := http.NewRequest("GET", "/", nil)
	require.Nil(suite.T(), err)
	r.Header.Set("token", "mytoken")
	token := suite.Service.getTokenFromRequest(r)
	require.Equal(suite.T(), "mytoken", token)

}
func (suite *TestSuite) TestgetTokenFromRequest_query() {
	r, err := http.NewRequest("GET", "/", nil)
	require.Nil(suite.T(), err)
	values := r.URL.Query()
	values.Set("token", "mytoken")
	r.URL.RawQuery = values.Encode()
	token := suite.Service.getTokenFromRequest(r)
	require.Equal(suite.T(), "mytoken", token)
}
func (suite *TestSuite) TestAuthenticateHandlerFunc() {
	suite.MockMetaDataController.On("ExamineObject").Once().Return(&emocks.MockObjectInfo{}, nil)
	suite.MockAuthService.On("Verify", "mytoken").Once().Return(user, &codes.Response{}, nil)
	r, err := http.NewRequest("GET", "/clawio/v1/metadata/examine/myblob", nil)
	r.Header.Set("token", "mytoken")
	require.Nil(suite.T(), err)
	w := httptest.NewRecorder()
	suite.Server.ServeHTTP(w, r)
	require.True(suite.T(), 200 <= w.Code && w.Code <= 299)
}

func (suite *TestSuite) TestAuthenticateHandlerFunc_withBadToken() {
	suite.MockAuthService.On("Verify", "").Once().Return(user, &codes.Response{}, errors.New("test error"))
	r, err := http.NewRequest("GET", "/clawio/v1/metadata/examine/myblob", nil)
	require.Nil(suite.T(), err)
	w := httptest.NewRecorder()
	suite.Server.ServeHTTP(w, r)
	require.Equal(suite.T(), 401, w.Code)
}
