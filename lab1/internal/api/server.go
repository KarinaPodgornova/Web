package api

import (
    "github.com/gin-gonic/gin"
    "github.com/sirupsen/logrus"
    "log"
    "lab1/internal/app/handler"
    "lab1/internal/app/repository"
         // Добавить этот импорт
)

func StartServer() {
    log.Println("Starting server")

    repo, err := repository.NewRepository()
    if err != nil {
        logrus.Error("ошибка инициализации репозитория")
    }

    handler := handler.NewHandler(repo)

    r := gin.Default()
    
   
    // добавляем наш html/шаблон
    r.LoadHTMLGlob("templates/*")
    r.Static("/static", "./resources")

    r.GET("/", handler.GetOrders)
    r.GET("/order/:id", handler.GetOrder)
    r.GET("/request", handler.GetRequest)

    r.Run()
    log.Println("Server down")
}
