package main

import (
	"fmt"
	"github.com/byteintellect/go_commons"
	"github.com/byteintellect/go_commons/cache"
	db2 "github.com/byteintellect/go_commons/db"
	"github.com/byteintellect/go_commons/entity"
	"github.com/byteintellect/protos_go/users/v1"
	"github.com/byteintellect/user_svc/config"
	"github.com/byteintellect/user_svc/pkg/domain"
	svc2 "github.com/byteintellect/user_svc/pkg/svc"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	grpcZap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/go-grpc-middleware/validator"
	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/infobloxopen/atlas-app-toolkit/gateway"
	"github.com/infobloxopen/atlas-app-toolkit/requestid"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"time"
)

func NewGRPCServer(cfg *config.UserSvcConfig, app *go_commons.BaseApp) (*grpc.Server, error) {

	grpcMux := grpc.NewServer(
		grpc.KeepaliveParams(
			keepalive.ServerParameters{
				Time:    time.Duration(cfg.ServerConfig.KeepAliveTime) * time.Second,
				Timeout: time.Duration(cfg.ServerConfig.KeepAliveTimeOut) * time.Second,
			},
		),
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				// logging middleware
				grpcZap.UnaryServerInterceptor(app.Logger()),

				// Request-Id interceptor
				requestid.UnaryServerInterceptor(),

				// Metrics middleware
				app.GrpcMetrics().UnaryServerInterceptor(),

				// validation middleware
				grpc_validator.UnaryServerInterceptor(),

				// collection operators middleware
				gateway.UnaryServerInterceptor(),

				// trace middleware
				otelgrpc.UnaryServerInterceptor(otelgrpc.WithTracerProvider(app.Tracer())),
			),
		),
		grpc.StreamInterceptor(app.GrpcMetrics().StreamServerInterceptor()),
	)

	factory := getDomainMappings()
	addressSvc := svc2.NewAddressSvc(db2.NewGORMRepository(db2.WithDb(app.Db()), db2.WithCreator(factory.GetMapping("addresses"))))
	identitySvc := svc2.NewIdentitySvc(db2.NewGORMRepository(db2.WithDb(app.Db()), db2.WithCreator(factory.GetMapping("addresses"))))
	userSvcCache := cache.NewRedisCache(
		fmt.Sprintf("%s:%s", cfg.CacheConfig.Host, cfg.CacheConfig.Port),
		cfg.CacheConfig.Password,
		1,
		app.Logger(),
		factory.GetMapping("users"),
		app.Tracer())
	userSvcImpl := svc2.NewUserServiceServer(
		db2.NewGORMRepository(
			db2.WithDb(app.Db()),
			db2.WithCreator(factory.GetMapping("users"))),
		identitySvc,
		addressSvc,
		userSvcCache,
		app.Logger())
	usersv1.RegisterUserServiceServer(grpcMux, userSvcImpl)

	// Register reflection service on gRPC server.
	app.GrpcMetrics().InitializeMetrics(grpcMux)
	reflection.Register(grpcMux)
	grpcPrometheus.Register(grpcMux)
	return grpcMux, nil
}

func getDomainMappings() entity.DomainFactory {
	factory := entity.NewDomainFactory()
	factory.RegisterMapping("users", func() entity.Base {
		return &domain.User{}
	})
	factory.RegisterMapping("identities", func() entity.Base {
		return &domain.Identity{}
	})
	factory.RegisterMapping("addresses", func() entity.Base {
		return &domain.Address{}
	})
	return *factory
}
