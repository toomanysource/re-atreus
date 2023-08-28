// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"

	"github.com/toomanysource/atreus/app/relation/internal/biz"
	"github.com/toomanysource/atreus/app/relation/internal/conf"
	"github.com/toomanysource/atreus/app/relation/internal/data"
	"github.com/toomanysource/atreus/app/relation/internal/server"
	"github.com/toomanysource/atreus/app/relation/internal/service"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(confServer *conf.Server, client *conf.Client, jwt *conf.JWT, confData *conf.Data, logger log.Logger) (*kratos.App, func(), error) {
	db := data.NewMysqlConn(confData)
	cacheClient := data.NewRedisConn(confData)
	dataData, cleanup, err := data.NewData(db, cacheClient, logger)
	if err != nil {
		return nil, nil, err
	}
	clientConn := server.NewUserClient(client, logger)
	relationRepo := data.NewRelationRepo(dataData, clientConn, logger)
	relationUsecase := biz.NewRelationUsecase(relationRepo, jwt, logger)
	relationService := service.NewRelationService(relationUsecase, logger)
	grpcServer := server.NewGRPCServer(confServer, relationService, logger)
	httpServer := server.NewHTTPServer(confServer, jwt, relationService, logger)
	app := newApp(logger, grpcServer, httpServer)
	return app, func() {
		cleanup()
	}, nil
}
