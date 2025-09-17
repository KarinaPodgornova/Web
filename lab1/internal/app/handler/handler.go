package handler

import (
  "github.com/gin-gonic/gin"
  "github.com/sirupsen/logrus"
  "lab1/internal/app/repository"
  "net/http"
  "strconv"
  "time"
  "strings"
)

type Handler struct {
  Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
  return &Handler{
    Repository: r,
  }
}

func (h *Handler) GetOrders(ctx *gin.Context) {
    var orders []repository.Order
    var err error

    // Получаем поисковый запрос
    searchQuery := ctx.Query("query")
    
    if searchQuery == "" {
        // Если запроса нет - показываем все товары
        orders, err = h.Repository.GetOrders()
        if err != nil {
            logrus.Error(err)
        }
    } else {
        // Если есть поисковый запрос - ищем по названию
        orders, err = h.Repository.GetOrdersByTitle(searchQuery)
        if err != nil {
            logrus.Error(err)
        }
    }

    ctx.HTML(http.StatusOK, "index.html", gin.H{
        "time":      time.Now().Format("15:04:05"),
        "orders":    orders, // Просто массив, а не map
        "cartCount": 2,
        "query":     searchQuery,
    })
}

func (h *Handler) GetOrder(ctx *gin.Context) {
  idStr := ctx.Param("id")
  id, err := strconv.Atoi(idStr)
  if err != nil {
    logrus.Error(err)
  }

  order, err := h.Repository.GetOrder(id)
  if err != nil {
    logrus.Error(err)
  }
  
  specsArray := strings.Split(order.Specs, "\n")

  ctx.HTML(http.StatusOK, "order.html", gin.H{
    "order": order,
    "specsArray": specsArray,
    "cartCount": 2,
  })
}

func (h *Handler) GetRequest(ctx *gin.Context) {
  // Получаем все товары
  orders, err := h.Repository.GetOrders()
  if err != nil {
    logrus.Error(err)
  }
  
  // Преобразуем массив в словарь для корзины
  cartItemsMap := make(map[int]repository.Order)
  // Берем первые 2 товара
  if len(orders) >= 2 {
    for i := 0; i < 2; i++ {
      cartItemsMap[orders[i].ID] = orders[i]
    }
  } else if len(orders) > 0 {
    for _, order := range orders {
      cartItemsMap[order.ID] = order
    }
  }

  ctx.HTML(http.StatusOK, "request.html", gin.H{
    "cartItems": cartItemsMap, // Теперь это map
    "cartCount": len(cartItemsMap),
  })
}

func (h *Handler) GetCart(ctx *gin.Context) {
  // Получаем все товары
  orders, err := h.Repository.GetOrders()
  if err != nil {
    logrus.Error(err)
  }
  
  // Преобразуем массив в словарь для корзины
  cartItemsMap := make(map[int]repository.Order)
  // Берем первые 2 товара
  if len(orders) >= 2 {
    for i := 0; i < 2; i++ {
      cartItemsMap[orders[i].ID] = orders[i]
    }
  } else {
    for _, order := range orders {
      cartItemsMap[order.ID] = order
    }
  }

  ctx.HTML(http.StatusOK, "request.html", gin.H{
    "cartItems": cartItemsMap, // Теперь это map
    "cartCount": len(cartItemsMap),
  })
}