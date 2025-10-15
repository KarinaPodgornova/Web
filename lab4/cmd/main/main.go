package main

import (
	"fmt"
	"html/template"

	"lab4/internal/app/config"
	"lab4/internal/app/dsn"
	"lab4/internal/app/handler"
	"lab4/internal/app/repository"
	"lab4/internal/pkg"
	_ "lab4/docs"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @title Current API
// @version 1.0
// @description API для определения необходимой силы тока
// @contact.name API Support
// @contact.url http://localhost
// @contact.email support@current.com
// @license.name MIT
// @host localhost
// @BasePath /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	router := gin.Default()
	conf, err := config.NewConfig()
	if err != nil {
		logrus.Fatalf("error loading config: %v", err)
	}

	router.SetFuncMap(template.FuncMap{
    "find_amperage": func(a, b float64) float64 {
        return a*1000 / b
    },
	})

	postgresString := dsn.FromEnv()
	fmt.Println(postgresString)

	rep, errRep := repository.NewRepository(postgresString)
	if errRep != nil {
		logrus.Fatalf("error initializing repository: %v", errRep)
	}

	hand := handler.NewHandler(rep)

	application := pkg.NewApp(conf, router, hand)
	application.RunApp()
}