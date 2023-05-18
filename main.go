package main

import (
	"book_management_system_backend/bootstrap"
	_ "book_management_system_backend/docs"
	"book_management_system_backend/utils"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// @title          Book Management System Backend
// @version        0.0.1
// @description    This is a Book Management System backend for Fudan 2023 midterm Project of Database course.
// @termsOfService https://swagger.io/terms/

// @contact.name   JingYiJun
// @contact.url    https://www.jingyijun.xyz
// @contact.email  jingyijun3104@outlook.com

// @license.name  Apache 2.0
// @license.url   https://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8000
// @BasePath  /api

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

//go:generate go install github.com/swaggo/swag/cmd/swag@latest
//go:generate swag init
func main() {
	app := bootstrap.InitFiberApp()

	go func() {
		if innerErr := app.Listen("0.0.0.0:8000"); innerErr != nil {
			log.Println(innerErr)
		}
	}()

	interrupt := make(chan os.Signal, 1)

	// wait for CTRL-C interrupt
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-interrupt

	// close app
	err := app.Shutdown()
	if err != nil {
		utils.Logger.Error("app shutdown error", zap.Error(err))
	}

	_ = utils.Logger.Sync()
}
