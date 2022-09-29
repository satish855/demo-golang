package main

import (
	"github.com/byteintellect/go_commons"
	"github.com/byteintellect/go_commons/config"
	"github.com/byteintellect/protos_go/users/v1"
	config2 "github.com/byteintellect/user_svc/config"
	"github.com/infobloxopen/atlas-app-toolkit/gorm/resource"
	"go.uber.org/zap"
	"log"
)

func getConfig() *config2.UserSvcConfig {
	var cfg config2.UserSvcConfig
	err := config.ReadFile("CONFIG_PATH", &cfg)
	if err != nil {
		log.Fatalf("error reading config")
	}
	return &cfg
}

func main() {
	cfg := getConfig()
	baseApp, err := go_commons.NewBaseApp(&cfg.BaseConfig)
	if err != nil {
		log.Fatalf("Error initializing application %v", err)
	}
	grpcServer, err := NewGRPCServer(cfg, baseApp)
	if err != nil {
		log.Fatalf("Error initializing grpc server %v", err)
	}
	doneC := make(chan error)

	// Init External
	go func() {
		doneC <- go_commons.ServeExternal(&cfg.BaseConfig, baseApp, grpcServer, usersv1.RegisterUserServiceHandlerFromEndpoint)
	}()
	if err := <-doneC; err != nil {
		baseApp.Logger().Fatal("Error Starting gRPC service", zap.Error(err))
	}
	resource.RegisterApplication(cfg.BaseConfig.AppName)
	resource.SetPlural()
}
