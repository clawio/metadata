package service

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	"github.com/NYTimes/gizmo/config"
	"github.com/NYTimes/gizmo/server"
	"github.com/clawio/authentication/lib"
	"github.com/clawio/entities"
	mock_metadatacontroller "github.com/clawio/metadata/metadatacontroller/mock"
	"github.com/clawio/sdk"
	"github.com/clawio/sdk/mocks"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var (
	examineURL string
	listURL    string
	deleteURL  string
	moveURL    string
	initURL    string
	metricsURL string
	user       = &entities.User{Username: "test"}
	jwtToken   string
)

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
			JWTKey:                       "secret",
			JWTSigningMethod:             "HS256",
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

	// Create the token
	authenticator := lib.NewAuthenticator(cfg.General.JWTKey, cfg.General.JWTSigningMethod)
	token, err := authenticator.CreateToken(user)
	require.Nil(suite.T(), err)
	jwtToken = token

	// set testing urls
	examineURL = path.Join(svc.Config.General.BaseURL, "/examine") + "/"
	listURL = path.Join(svc.Config.General.BaseURL, "/list") + "/"
	deleteURL = path.Join(svc.Config.General.BaseURL, "/delete") + "/"
	moveURL = path.Join(svc.Config.General.BaseURL, "/move") + "/"
	initURL = path.Join(svc.Config.General.BaseURL, "/init")
	metricsURL = path.Join(svc.Config.General.BaseURL, "/metrics")
}

func (suite *TestSuite) TeardownTest() {
	os.Remove("/tmp/t/test")
}

func (suite *TestSuite) TestNew() {
	cfg := &Config{
		Server: &config.Server{},
		General: &GeneralConfig{
			AuthenticationServiceBaseURL: "http://localhost:58001/api/auth/",
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
	suite.Service.Config.General.BaseURL = "/api/metadata"
	require.Equal(suite.T(), suite.Service.Config.General.BaseURL, suite.Service.Prefix())
}
func (suite *TestSuite) TestPrefix_withEmpty() {
	suite.Service.Config.General.BaseURL = ""
	require.Equal(suite.T(), "/", suite.Service.Prefix())
}

func (suite *TestSuite) TestMetrics() {
	r, err := http.NewRequest("GET", metricsURL, nil)
	require.Nil(suite.T(), err)
	w := httptest.NewRecorder()
	suite.Server.ServeHTTP(w, r)
	require.Equal(suite.T(), 200, w.Code)
}
func setToken(r *http.Request) {
	r.Header.Set("Authorization", "bearer "+jwtToken)

}
