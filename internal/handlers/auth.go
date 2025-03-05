package handlers

import (
    "io"
    "log"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
)

type AuthHandler struct {
    authServiceAddr string // Base URL of Auth Service HTTP API
}

func NewAuthHandler(authServiceAddr string) *AuthHandler {
    return &AuthHandler{authServiceAddr: authServiceAddr}
}

// Login proxies the login request to Auth Service
func (h *AuthHandler) Login(c *gin.Context) {
    resp, err := h.proxyRequest(c, "/auth/login")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to Auth Service"})
        log.Printf("Login proxy error: %v", err)
        return
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
        return
    }
    c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
}

// Register proxies the register request to Auth Service
func (h *AuthHandler) Register(c *gin.Context) {
    resp, err := h.proxyRequest(c, "/auth/register")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to Auth Service"})
        log.Printf("Register proxy error: %v", err)
        return
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
        return
    }
    c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
}

// proxyRequest forwards the request to Auth Service
func (h *AuthHandler) proxyRequest(c *gin.Context, path string) (*http.Response, error) {
    url := h.authServiceAddr + path
    req, err := http.NewRequest(c.Request.Method, url, c.Request.Body)
    if err != nil {
        return nil, err
    }
    
    // Копируем все заголовки запроса
    for k, v := range c.Request.Header {
        req.Header[k] = v
    }
    
    // Устанавливаем контекст для предотвращения утечек
    req = req.WithContext(c.Request.Context())
    
    client := &http.Client{
        Timeout: 10 * time.Second, // Таймаут для предотвращения зависания
    }
    return client.Do(req)
}
