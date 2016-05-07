package main

import (
	"flag"
	"github.com/NYTimes/gizmo/config"
	"github.com/NYTimes/gizmo/server"
	"github.com/clawio/metadata/service"
)

func main() {
	flag.Parse()
	var cfg *service.Config
	config.LoadJSONFile(*config.ConfigLocationCLI, &cfg)

	server.Init("metadata-service", cfg.Server)

	svc, err := service.New(cfg)
	if err != nil {
		server.Log.Fatal("unable to create service: ", err)
	}
	err = server.Register(svc)
	if err != nil {
		server.Log.Fatal("unable to register service: ", err)
	}

	err = server.Run()
	if err != nil {
		server.Log.Fatal("server encountered a fatal error: ", err)
	}
}
