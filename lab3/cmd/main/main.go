package main

import (
	"fmt"
	"html/template"

	"lab3/internal/app/config"
	"lab3/internal/app/dsn"
	"lab3/internal/app/handler"
	"lab3/internal/app/repository"
	"lab3/internal/pkg"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

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