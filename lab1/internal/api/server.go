package api

import (
    "github.com/gin-gonic/gin"
    "github.com/sirupsen/logrus"
    "log"
    "lab1/internal/app/handler"
    "lab1/internal/app/repository"
    "html/template" // Добавить этот импорт
    "strings"       // Добавить этот импорт
    "strconv"       // Добавить этот импорт
)

func StartServer() {
    log.Println("Starting server")

    repo, err := repository.NewRepository()
    if err != nil {
        logrus.Error("ошибка инициализации репозитория")
    }

    handler := handler.NewHandler(repo)

    r := gin.Default()
    
    // Добавляем функции для шаблонов ПЕРЕД LoadHTMLGlob
    r.SetFuncMap(template.FuncMap{
        "split": strings.Split,
        "replace": strings.ReplaceAll,
        "atoi": func(s string) int {
            i, _ := strconv.Atoi(s)
            return i
        },
        "add": func(a, b int) int {
            return a + b
        },
    })

    // добавляем наш html/шаблон
    r.LoadHTMLGlob("templates/*")

    r.Static("/static", "./resources")

    r.GET("/hello", handler.GetOrders)
    r.GET("/order/:id", handler.GetOrder)
    r.GET("/request", handler.GetRequest)

    r.Run()
    log.Println("Server down")
}