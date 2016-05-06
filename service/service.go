package service

import (
	"errors"
	"net/http"

	"github.com/NYTimes/gizmo/config"
	"github.com/clawio/authentication/lib"
	"github.com/clawio/metadata/metadatacontroller"
	"github.com/clawio/metadata/metadatacontroller/simple"
	"github.com/clawio/sdk"
	"github.com/prometheus/client_golang/prometheus"
)

type (
	// Service implements server.Service and
	// handle all requests to the server.
	Service struct {
		Config             *Config
		SDK                *sdk.SDK
		MetaDataController metadatacontroller.MetaDataController
	}

	// Config is a struct that holds the
	// configuration for Service
	Config struct {
		Server             *config.Server
		General            *GeneralConfig
		MetaDataController *MetaDataControllerConfig
	}

	// GeneralConfig contains configuration parameters
	// for general parts of the service.
	GeneralConfig struct {
		BaseURL                      string
		JWTKey, JWTSigningMethod     string
		AuthenticationServiceBaseURL string
	}

	// MetaDataControllerConfig is a struct that holds
	// configuration parameters for a metadata controller.
	MetaDataControllerConfig struct {
		Type              string
		SimpleMetaDataDir string
		SimpleTempDir     string
	}
)

// New will instantiate and return
// a new Service that implements server.Service.
func New(cfg *Config) (*Service, error) {
	if cfg == nil {
		return nil, errors.New("config is nil")
	}
	if cfg.General == nil {
		return nil, errors.New("config.General is nil")
	}
	if cfg.MetaDataController == nil {
		return nil, errors.New("config.MetaDataController is  nil")
	}

	urls := &sdk.ServiceEndpoints{}
	urls.AuthServiceBaseURL = cfg.General.AuthenticationServiceBaseURL
	s := sdk.New(urls, nil)

	metadataController := getMetaDataController(cfg.MetaDataController)
	return &Service{Config: cfg, SDK: s, MetaDataController: metadataController}, nil
}

func getMetaDataController(cfg *MetaDataControllerConfig) metadatacontroller.MetaDataController {
	opts := &simple.Options{
		MetaDataDir: cfg.SimpleMetaDataDir,
		TempDir:     cfg.SimpleTempDir,
	}
	return simple.New(opts)
}

// Prefix returns the string prefix used for all endpoints within
// this service.
func (s *Service) Prefix() string {
	if s.Config.General.BaseURL == "" {
		return "/"
	}
	return s.Config.General.BaseURL
}

// Middleware provides an http.Handler hook wrapped around all requests.
// In this implementation, we authenticate the request.
func (s *Service) Middleware(h http.Handler) http.Handler {
	return h
}

// Endpoints is a listing of all endpoints available in the Service.
func (s *Service) Endpoints() map[string]map[string]http.HandlerFunc {
	authenticator := lib.NewAuthenticator(s.Config.General.JWTKey, s.Config.General.JWTSigningMethod)

	return map[string]map[string]http.HandlerFunc{
		"/metrics": {
			"GET": func(w http.ResponseWriter, r *http.Request) {
				prometheus.Handler().ServeHTTP(w, r)
			},
		},
		"/init": {
			"POST": prometheus.InstrumentHandlerFunc("/init", authenticator.JWTHandlerFunc(s.Init)),
		},
		"/examine/{path:.*}": {
			"GET": prometheus.InstrumentHandlerFunc("/examine", authenticator.JWTHandlerFunc(s.ExamineObject)),
		},
		"/list/{path:.*}": {
			"GET": prometheus.InstrumentHandlerFunc("/list", authenticator.JWTHandlerFunc(s.ListTree)),
		},
		"/move/{path:.*}": {
			"POST": prometheus.InstrumentHandlerFunc("/move", authenticator.JWTHandlerFunc(s.MoveObject)),
		},
		"/delete/{path:.*}": {
			"DELETE": prometheus.InstrumentHandlerFunc("/delete", authenticator.JWTHandlerFunc(s.DeleteObject)),
		},
	}
}
