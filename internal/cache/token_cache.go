package cache

import (
    "sync"
    "time"
)

type tokenInfo struct {
    valid     bool
    expiresAt time.Time
}

type TokenCache struct {
    cache map[string]tokenInfo
    mu    sync.RWMutex
    ttl   time.Duration
}

func NewTokenCache() *TokenCache {
    return &TokenCache{
        cache: make(map[string]tokenInfo),
        ttl:   15 * time.Minute, // Время жизни токена в кэше
    }
}

func (c *TokenCache) Get(token string) (bool, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    info, exists := c.cache[token]
    if !exists {
        return false, false
    }
    
    // Проверяем время жизни кэша
    if time.Now().After(info.expiresAt) {
        // Токен устарел, удаляем его из кэша
        // Используем горутину для избежания блокировки чтения
        go func() {
            c.mu.Lock()
            delete(c.cache, token)
            c.mu.Unlock()
        }()
        return false, false
    }
    
    return info.valid, true
}

func (c *TokenCache) Set(token string, valid bool) {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    c.cache[token] = tokenInfo{
        valid:     valid,
        expiresAt: time.Now().Add(c.ttl),
    }
}

// Очистка устаревших токенов
func (c *TokenCache) Cleanup() {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    now := time.Now()
    for token, info := range c.cache {
        if now.After(info.expiresAt) {
            delete(c.cache, token)
        }
    }
}
