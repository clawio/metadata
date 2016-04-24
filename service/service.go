package service

import (
	"errors"
	"net/http"

	"github.com/NYTimes/gizmo/config"
	"github.com/clawio/keys"
	"github.com/clawio/metadata/metadatacontroller"
	"github.com/clawio/metadata/metadatacontroller/simple"
	"github.com/clawio/sdk"
	"github.com/gorilla/context"
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
	return "/clawio/v1/metadata"
}

// Middleware provides an http.Handler hook wrapped around all requests.
// In this implementation, we authenticate the request.
func (s *Service) Middleware(h http.Handler) http.Handler {
	return h
}

// Endpoints is a listing of all endpoints available in the Service.
func (s *Service) Endpoints() map[string]map[string]http.HandlerFunc {
	return map[string]map[string]http.HandlerFunc{
		"/metrics": {
			"GET": func(w http.ResponseWriter, r *http.Request) {
				prometheus.Handler().ServeHTTP(w, r)
			},
		},
		"/init": {
			"POST": prometheus.InstrumentHandlerFunc("/init", s.authenticateHandlerFunc(s.Init)),
		},
		"/examine/{path:.*}": {
			"GET": prometheus.InstrumentHandlerFunc("/examine", s.authenticateHandlerFunc(s.ExamineObject)),
		},
		"/listtree/{path:.*}": {
			"GET": prometheus.InstrumentHandlerFunc("/listtree", s.authenticateHandlerFunc(s.ListTree)),
		},
	}
}

func (s *Service) getTokenFromRequest(r *http.Request) string {
	if t := r.Header.Get("token"); t != "" {
		return t
	}
	return r.URL.Query().Get("token")
}

func (s *Service) authenticateHandlerFunc(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := s.getTokenFromRequest(r)
		user, _, err := s.SDK.Auth.Verify(token)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		context.Set(r, keys.UserKey, user)
		handler(w, r)
	}
}
