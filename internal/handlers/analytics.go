package handlers

import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
)

type AnalyticsHandler struct{}

func NewAnalyticsHandler() *AnalyticsHandler {
    return &AnalyticsHandler{}
}

// DemandData представляет данные о спросе на продукт
type DemandData struct {
    ProductID    string    `json:"product_id"`
    Name         string    `json:"name"`
    CurrentScore float64   `json:"current_score"`
    History      []float64 `json:"history"`
    Forecast     float64   `json:"forecast"`
    Timestamp    time.Time `json:"timestamp"`
}

// GetDemand возвращает данные о спросе (для MVP используются тестовые данные)
func (h *AnalyticsHandler) GetDemand(c *gin.Context) {
    // Получаем product_id из запроса, если он есть
    productID := c.Query("product_id")
    if productID == "" {
        productID = "default-product-1"
    }
    
    // В реальной системе здесь будет запрос к сервису аналитики
    // для получения реальных данных о спросе
    data := DemandData{
        ProductID:    productID,
        Name:         "Тестовый продукт",
        CurrentScore: 76.5,
        History:      []float64{65.2, 68.7, 72.1, 74.8, 76.5},
        Forecast:     79.2,
        Timestamp:    time.Now(),
    }
    
    c.JSON(http.StatusOK, gin.H{
        "status": "success",
        "data":   data,
    })
}
