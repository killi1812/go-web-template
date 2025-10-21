package main

import (
	"template/app"
	"template/controller"
	"template/service"
	"template/util/minio"

	"go.uber.org/zap"

	_ "template/docs"
)

// @securitydefinitions.bearerauth BearerAuth

func init() {
	app.Setup()
}

func main() {
	err := app.LoadConfig()
	if err != nil {
		panic(err)
	}
	// Provide logger
	app.Provide(zap.S)
	app.Provide(minio.New)

	app.Provide(service.NewDiscordService)
	app.Provide(service.NewUserCrudService)
	app.Provide(service.NewAuthService)

	app.RegisterController(controller.NewGameCnt)
	app.RegisterController(controller.NewInfoCnt)
	app.RegisterController(controller.NewUserCtn)
	app.RegisterController(controller.NewAuthCtn)

	app.Start()
}
